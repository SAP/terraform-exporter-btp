package cmd

import (
	"btptfexport/tfutils"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func exportSubaccountEntitlements(subaccountID string, configDir string, filterValues []string) {

	dataBlock, err := readDataSource(subaccountID, SubaccountEntitlementType)
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

	importBlock, err := getEntitlementsImportBlock(data, subaccountID, filterValues)
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

func getEntitlementsImportBlock(data map[string]interface{}, subaccountId string, filterValues []string) (string, error) {

	resourceDoc, err := getDocByResourceName(ResourcesKind, SubaccountEntitlementType)
	if err != nil {
		return "", err
	}

	var importBlock string

	if len(filterValues) != 0 {
		var subaccountAllEntitlements []string
		for key, value := range data {
			subaccountAllEntitlements = append(subaccountAllEntitlements, strings.Replace(key, ":", "_", -1))
			if slices.Contains(filterValues, strings.Replace(key, ":", "_", -1)) {
				importBlock += templateEntitlementImport(key, value, subaccountId, resourceDoc)
			}
		}

		missingEntitlement, subset := isSubset(subaccountAllEntitlements, filterValues)

		if !subset {
			return "", fmt.Errorf("entitlement %s not found in the subaccount. Please adjust it in the provided file", missingEntitlement)
		}

	} else {
		for key, value := range data {
			importBlock += templateEntitlementImport(key, value, subaccountId, resourceDoc)
		}
	}

	return importBlock, nil
}

func templateEntitlementImport(key string, value interface{}, subaccountId string, resourceDoc tfutils.EntityDocs) string {
	template := strings.Replace(resourceDoc.Import, "<resource_name>", strings.Replace(key, ":", "_", -1), -1)
	template = strings.Replace(template, "<subaccount_id>", subaccountId, -1)
	if subMap, ok := value.(map[string]interface{}); ok {
		for subKey, subValue := range subMap {
			template = strings.Replace(template, "<"+subKey+">", fmt.Sprintf("%v", subValue), -1)
		}
	}
	return template + "\n"
}
