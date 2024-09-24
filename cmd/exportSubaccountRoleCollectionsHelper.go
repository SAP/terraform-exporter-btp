package cmd

import (
	"btptfexport/files"
	"btptfexport/output"
	"btptfexport/tfutils"
	"fmt"
	"log"
	"slices"
	"strings"
)

func exportSubaccountRoleCollections(subaccountID string, configDir string, filterValues []string) {

	output.AddNewLine()
	spinner, err := output.StartSpinner("crafting import block for " + strings.ToUpper(tfutils.SubaccountRoleCollectionType))
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	data, err := tfutils.FetchImportConfiguration(subaccountID, tfutils.SubaccountRoleCollectionType, tfutils.TmpFolder)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	importBlock, err := getSubaccountRoleCollectionsImportBlock(data, subaccountID, filterValues)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	if len(importBlock) == 0 {
		log.Println("no role collection found for the given subaccount")
		return
	}

	err = files.WriteImportConfiguration(configDir, tfutils.SubaccountRoleCollectionType, importBlock)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	err = output.StopSpinner(spinner)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}
}

func getSubaccountRoleCollectionsImportBlock(data map[string]interface{}, subaccountId string, filterValues []string) (string, error) {

	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.SubaccountRoleCollectionType)
	if err != nil {
		log.Fatalf("read doc failed!")
		return "", err
	}

	var importBlock string
	roleCollections := data["values"].([]interface{})

	if len(filterValues) != 0 {
		var subaccountAllRoleCollections []string

		for _, value := range roleCollections {
			roleCollection := value.(map[string]interface{})
			resourceName := output.FormatRoleCollectionResourceName(fmt.Sprintf("%v", roleCollection["name"]))
			subaccountAllRoleCollections = append(subaccountAllRoleCollections, resourceName)
			if slices.Contains(filterValues, resourceName) {
				importBlock += templateRoleCollectionImport(roleCollection, subaccountId, resourceDoc)
			}
		}

		missingRoleCollection, subset := isSubset(subaccountAllRoleCollections, filterValues)

		if !subset {
			return "", fmt.Errorf("role collection %s not found in the subaccount. Please adjust it in the provided file", missingRoleCollection)
		}

	} else {
		for _, value := range roleCollections {
			roleCollection := value.(map[string]interface{})
			importBlock += templateRoleCollectionImport(roleCollection, subaccountId, resourceDoc)
		}
	}

	return importBlock, nil
}

func templateRoleCollectionImport(roleCollection map[string]interface{}, subaccountId string, resourceDoc tfutils.EntityDocs) string {

	resourceDoc.Import = strings.Replace(resourceDoc.Import, "'", "", -1)
	resourceName := output.FormatRoleCollectionResourceName(fmt.Sprintf("%v", roleCollection["name"]))
	template := strings.Replace(resourceDoc.Import, "<resource_name>", resourceName, -1)
	template = strings.Replace(template, "<subaccount_id>", subaccountId, -1)
	template = strings.Replace(template, "<name>", fmt.Sprintf("%v", roleCollection["name"]), -1)
	return template + "\n"
}
