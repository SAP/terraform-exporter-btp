package cmd

import (
	"btptfexporter/tfutils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-exec/tfexec"
)

func exportSubaccount(subaccountID string, configDir string) {

	dataBlock, err := readSubaccountDataSource(subaccountID)
	if err != nil {
		log.Fatalf("error getting data source: %v", err)
		return
	}

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("error getting current directory: %v", err)
		return
	}
	dataBlockFile := filepath.Join(TmpFolder, "main.tf")
	err = tfutils.CreateFileWithContent(dataBlockFile, dataBlock)
	if err != nil {
		log.Fatalf("create file %s failed!", dataBlockFile)
		return
	}

	jsonBytes, err := getSubaccountTfStateData(TmpFolder)
	if err != nil {
		log.Fatalf("error json.Marshal: %v", err)
		return
	}

	var data map[string]interface{}
	subaccountName := string(jsonBytes)
	err = json.Unmarshal([]byte(subaccountName), &data)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	importBlock, err := getSubaccountImportBlock(data, subaccountID)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	importFileName := "subaccount_import.tf"
	importFileName = filepath.Join(currentDir, configDir, importFileName)

	err = tfutils.CreateFileWithContent(importFileName, importBlock)
	if err != nil {
		log.Fatalf("create file %s failed!", dataBlockFile)
		return
	}

	log.Println("btp subaccount has been exported. Please check " + configDir + " folder")
}

func readSubaccountDataSource(subaccountId string) (string, error) {
	choice := "btp_subaccount"
	dsDoc, err := tfutils.GetDocsForResource("SAP", "btp", "btp", "data-sources", choice, BtpProviderVersion, "github.com")

	if err != nil {
		log.Fatalf("read doc failed!")
		return "", err
	}
	dataBlock := strings.Replace(dsDoc.Import, "The ID of the subaccount", subaccountId, -1)
	return dataBlock, nil

}

func getSubaccountTfStateData(configDir string) ([]byte, error) {
	execPath, err := exec.LookPath("terraform")
	if err != nil {
		log.Fatalf("error finding Terraform: %v", err)
		return nil, err
	}
	// create a new Terraform instance
	tf, err := tfexec.NewTerraform(configDir, execPath)
	if err != nil {
		log.Fatalf("error running NewTerraform: %v", err)
		return nil, err
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		log.Fatalf("error running Init: %v", err)
		return nil, err
	}
	err = tf.Apply(context.Background())
	if err != nil {
		log.Fatalf("error running Apply: %v", err)
		return nil, err
	}

	state, err := tf.Show(context.Background())
	if err != nil {
		log.Fatalf("error running Show: %v", err)
		return nil, err
	}

	jsonBytes, err := json.Marshal(state.Values.RootModule.Resources[0].AttributeValues)
	if err != nil {
		log.Fatalf("error json.Marshal: %v", err)
		return nil, err
	}

	return jsonBytes, nil

}

func getSubaccountImportBlock(data map[string]interface{}, subaccountId string) (string, error) {
	choice := "btp_subaccount"
	resource_doc, err := tfutils.GetDocsForResource("SAP", "btp", "btp", "resources", choice, BtpProviderVersion, "github.com")
	if err != nil {
		log.Fatalf("read doc failed")
		return "", err
	}

	var importBlock string
	saName := fmt.Sprintf("%v", data["name"])
	template := strings.Replace(resource_doc.Import, "<resource_name>", saName, -1)
	template = strings.Replace(template, "<subaccount_id>", subaccountId, -1)
	importBlock += template + "\n"

	return importBlock, nil
}
