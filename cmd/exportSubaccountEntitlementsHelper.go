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

	jsonBytes, err := GetTfStateData(TmpFolder, SubaccountEntitlementType)
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
