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

func exportSubaccountEntitlements(subaccountID string, configDir string) {

	dataBlock, err := readSubaccountEntilementsDataSource(subaccountID)
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

	jsonBytes, err := getEntitlementTfStateData(TmpFolder)
	if err != nil {
		log.Fatalf("error json.Marshal: %v", err)
		return
	}

	jsonString := string(jsonBytes)
	var data map[string]interface{}
	err = json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	importBlock, err := getEntitlementsImportBlock(data, subaccountID)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	if len(importBlock) == 0 {
		log.Println("no entitlement found for the given subaccount")
		return
	}

	importFileName := "subaccount_entitlements_import.tf"
	importFileName = filepath.Join(currentDir, configDir, importFileName)

	err = tfutils.CreateFileWithContent(importFileName, importBlock)
	if err != nil {
		log.Fatalf("create file %v failed!", dataBlockFile)
		return
	}

	log.Println("BTP subaccount entitlements has been exported. Please check " + configDir + " folder")

}

// this function read the data source document and return the data block to use to get the resoure state
func readSubaccountEntilementsDataSource(subaccountId string) (string, error) {
	choice := "btp_subaccount_entitlements"
	dsDoc, err := tfutils.GetDocsForResource("SAP", "btp", "btp", "data-sources", choice, BtpProviderVersion, "github.com")

	if err != nil {
		log.Fatalf("read doc failed")
		return "", err
	}
	dataBlock := strings.Replace(dsDoc.Import, dsDoc.Attributes["subaccount_id"], subaccountId, -1)
	return dataBlock, nil

}

// this function runs terraform actions and return the state values to create import block
func getEntitlementTfStateData(configDir string) ([]byte, error) {
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

	jsonBytes, err := json.Marshal(state.Values.RootModule.Resources[0].AttributeValues["values"])
	if err != nil {
		log.Fatalf("error json.Marshal: %v", err)
		return nil, err
	}

	return jsonBytes, nil

}

func getEntitlementsImportBlock(data map[string]interface{}, subaccountId string) (string, error) {
	choice := "btp_subaccount_entitlement"
	resource_doc, err := tfutils.GetDocsForResource("SAP", "btp", "btp", "resources", choice, BtpProviderVersion, "github.com")
	if err != nil {
		log.Fatalf("read doc failed")
		return "", err
	}

	var importBlock string
	for key, value := range data {

		template := strings.Replace(resource_doc.Import, "<resource_name>", strings.Replace(key, ":", "_", -1), -1)
		template = strings.Replace(template, "<subaccount_id>", subaccountId, -1)
		if subMap, ok := value.(map[string]interface{}); ok {
			for subKey, subValue := range subMap {
				template = strings.Replace(template, "<"+subKey+">", fmt.Sprintf("%v", subValue), -1)
			}
		}
		importBlock += template + "\n"
	}

	return importBlock, nil
}
