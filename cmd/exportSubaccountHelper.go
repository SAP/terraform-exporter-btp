package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func exportSubaccount(subaccountID string, configDir string, filterValues []string) {

	data, err := fetchImportConfiguration(subaccountID, SubaccountType, TmpFolder)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	importBlock, err := getSubaccountImportBlock(data, subaccountID, filterValues)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	err = writeImportConfiguration(configDir, SubaccountType, importBlock)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}
}

func getSubaccountImportBlock(data map[string]interface{}, subaccountId string, filterValues []string) (string, error) {

	if len(filterValues) != 0 {
		if filterValues[0] != fmt.Sprintf("%v", data["name"]) {
			log.Println("Error:", fmt.Errorf("subaccount %s not found. Please adjust it in the provided file", filterValues[0]))
			os.Exit(0)
			return "", nil
		}
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
