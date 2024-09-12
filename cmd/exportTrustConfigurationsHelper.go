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

func exportTrustConfigurations(subaccountID string, configDir string) {
	dataBlock, err := readSubaccountTrustConfigurationsDataSource(subaccountID)
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

	importBlock, err := getTrustConfigurationsImportBlock(data, subaccountID)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	if len(importBlock) == 0 {
		log.Println("No trust configuration found for the given subaccount")
		return
	}

	importFileName := "subaccount_trust_configurations_import.tf"
	importFileName = filepath.Join(currentDir, configDir, importFileName)

	err = tfutils.CreateFileWithContent(importFileName, importBlock)
	if err != nil {
		log.Fatalf("create file %s failed!", dataBlockFile)
		return
	}

	log.Println("subaccount trust configuration has been exported. Please check " + configDir + " folder")

}

func readSubaccountTrustConfigurationsDataSource(subaccountId string) (string, error) {
	choice := "btp_subaccount_trust_configurations"
	dsDoc, err := tfutils.GetDocsForResource("SAP", "btp", "btp", "data-sources", choice, BtpProviderVersion, "github.com")

	if err != nil {
		log.Fatalf("read doc failed")
		return "", err
	}
	dataBlock := strings.Replace(dsDoc.Import, dsDoc.Attributes["subaccount_id"], subaccountId, -1)
	return dataBlock, nil

}

func getTrustConfigurationsImportBlock(data map[string]interface{}, subaccountId string) (string, error) {
	choice := "btp_subaccount_trust_configuration"
	resource_doc, err := tfutils.GetDocsForResource("SAP", "btp", "btp", "resources", choice, BtpProviderVersion, "github.com")
	if err != nil {
		log.Fatalf("read doc failed!")
		return "", err
	}

	var importBlock string
	trusts := data["values"].([]interface{})

	for x, value := range trusts {
		trust := value.(map[string]interface{})
		template := strings.Replace(resource_doc.Import, "<resource_name>", "trust"+fmt.Sprint(x), -1)
		template = strings.Replace(template, "<subaccount_id>", subaccountId, -1)
		template = strings.Replace(template, "<origin>", fmt.Sprintf("%v", trust["origin"]), -1)
		importBlock += template + "\n"
	}

	return importBlock, nil
}
