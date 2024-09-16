package cmd

import (
	"btptfexport/tfutils"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func exportSubaccount(subaccountID string, configDir string, filterValues []string) {

	dataBlock, err := readDataSource(subaccountID, SubaccountType)
	if err != nil {
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

	jsonBytes, err := GetTfStateData(TmpFolder, SubaccountType)
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

	if len(filterValues) != 0 {
		if filterValues[0] != fmt.Sprintf("%v", data["name"]) {
			log.Println("Error:", fmt.Errorf("subaccount %s not found. Please adjust it in the provided file", filterValues[0]))
			os.Exit(0)
		}
	}

	importBlock, err := getSubaccountImportBlock(data, subaccountID, nil)
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

func getSubaccountImportBlock(data map[string]interface{}, subaccountId string, filterValues []string) (string, error) {

	if len(filterValues) != 0 {
		return "", fmt.Errorf("filter values for subaccounts are not supported")
	}

	resourceDoc, err := getDocByResourceName(ResourcesKind, SubaccountType)
	if err != nil {
		return "", err
	}

	var importBlock string
	saName := fmt.Sprintf("%v", data["name"])
	template := strings.Replace(resourceDoc.Import, "<resource_name>", saName, -1)
	template = strings.Replace(template, "<subaccount_id>", subaccountId, -1)
	importBlock += template + "\n"

	return importBlock, nil
}
