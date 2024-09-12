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

func exportSubaccountSubscriptions(subaccountID string, configDir string) {

	dataBlock, err := readSubaccountSubscriptionDataSource(subaccountID)
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

	jsonBytes, err := getSubscriptionsTfStateData(TmpFolder)
	if err != nil {
		log.Fatalf("error json.Marshal: %v", err)
		return
	}

	var data map[string]interface{}
	jsonString := string(jsonBytes)

	err = json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	importBlock, err := getSubscriptionsImportBlock(data, subaccountID)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	if len(importBlock) == 0 {
		log.Println("No subscription found for the given subaccount")
		return
	}

	importFileName := "subaccount_subscriptions_import.tf"
	importFileName = filepath.Join(currentDir, configDir, importFileName)

	err = tfutils.CreateFileWithContent(importFileName, importBlock)
	if err != nil {
		log.Fatalf("create file %v failed!", dataBlockFile)
		return
	}

	log.Println("subaccount subscriptions has been exported. Please check " + configDir + " folder")
}

// this function read the data source document and return the data block to use to get the resoure state
func readSubaccountSubscriptionDataSource(subaccountId string) (string, error) {
	choice := "btp_subaccount_subscriptions"
	dsDoc, err := tfutils.GetDocsForResource("SAP", "btp", "btp", "data-sources", choice, BtpProviderVersion, "github.com")

	if err != nil {
		log.Fatalf("read doc failed")
		return "", err
	}
	dataBlock := strings.Replace(dsDoc.Import, dsDoc.Attributes["subaccount_id"], subaccountId, -1)
	return dataBlock, nil

}

func getSubscriptionsTfStateData(configDir string) ([]byte, error) {
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

func getSubscriptionsImportBlock(data map[string]interface{}, subaccountId string) (string, error) {
	choice := "btp_subaccount_subscription"
	resource_doc, err := tfutils.GetDocsForResource("SAP", "btp", "btp", "resources", choice, BtpProviderVersion, "github.com")
	if err != nil {
		log.Fatalf("read doc failed")
		return "", err
	}

	var importBlock string
	subscriptions := data["values"].([]interface{})

	for _, value := range subscriptions {
		subscription := value.(map[string]interface{})
		if fmt.Sprintf("%v", subscription["state"]) != "NOT_SUBSCRIBED" {
			template := strings.Replace(resource_doc.Import, "<resource_name>", strings.Replace(fmt.Sprintf("%v", subscription["app_name"]), "-", "_", -1), -1)
			template = strings.Replace(template, "<subaccount_id>", subaccountId, -1)
			template = strings.Replace(template, "<app_name>", fmt.Sprintf("%v", subscription["app_name"]), -1)
			template = strings.Replace(template, "<plan_name>", fmt.Sprintf("%v", subscription["plan_name"]), -1)
			importBlock += template + "\n"
		}
	}

	return importBlock, nil
}
