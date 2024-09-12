package cmd

import (
	"btptfexporter/tfutils"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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

	jsonBytes, err := GetTfStateData(TmpFolder)
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
