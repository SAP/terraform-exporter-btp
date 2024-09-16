package cmd

import (
	"btptfexport/tfutils"
	"fmt"
	"log"
	"slices"
	"strings"
)

func exportSubaccountEntitlements(subaccountID string, configDir string, filterValues []string) {

	data, err := fetchImportConfiguration(subaccountID, SubaccountEntitlementType, TmpFolder)
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

	err = writeImportConfiguration(configDir, SubaccountEntitlementType, importBlock)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}
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
