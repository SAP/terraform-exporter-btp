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

func exportSubaccountServiceBindings(subaccountID string, configDir string, filterValues []string) {

	output.AddNewLine()
	spinner, err := output.StartSpinner("crafting import block for " + strings.ToUpper(tfutils.SubaccountServiceBindingType))
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	data, err := tfutils.FetchImportConfiguration(subaccountID, tfutils.SubaccountServiceBindingType, tfutils.TmpFolder)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	importBlock, err := getSubaccountServiceBindingImportBlock(data, subaccountID, filterValues)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	if len(importBlock) == 0 {
		log.Println("no service binding found for the given subaccount")
		return
	}

	err = files.WriteImportConfiguration(configDir, tfutils.SubaccountServiceBindingType, importBlock)
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

func getSubaccountServiceBindingImportBlock(data map[string]interface{}, subaccountId string, filterValues []string) (string, error) {

	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.SubaccountServiceBindingType)
	if err != nil {
		log.Fatalf("read doc failed!")
		return "", err
	}

	var importBlock string
	serviceBindings := data["values"].([]interface{})

	if len(filterValues) != 0 {
		var subaccountAllServiceBindings []string

		for _, value := range serviceBindings {
			binding := value.(map[string]interface{})
			resourceName := fmt.Sprintf("%v", binding["name"])
			subaccountAllServiceBindings = append(subaccountAllServiceBindings, resourceName)
			if slices.Contains(filterValues, resourceName) {
				importBlock += templateServiceBindingImport(binding, subaccountId, resourceDoc)
			}
		}

		missingBinding, subset := isSubset(subaccountAllServiceBindings, filterValues)

		if !subset {
			return "", fmt.Errorf("service binding %s not found in the subaccount. Please adjust it in the provided file", missingBinding)
		}

	} else {
		for _, value := range serviceBindings {
			binding := value.(map[string]interface{})
			importBlock += templateServiceBindingImport(binding, subaccountId, resourceDoc)
		}
	}

	return importBlock, nil
}

func templateServiceBindingImport(binding map[string]interface{}, subaccountId string, resourceDoc tfutils.EntityDocs) string {
	resourceName := output.FormatServiceBindingResourceName(fmt.Sprintf("%v", binding["name"]))
	template := strings.Replace(resourceDoc.Import, "<resource_name>", resourceName, -1)
	template = strings.Replace(template, "<subaccount_id>", subaccountId, -1)
	template = strings.Replace(template, "<service_binding_id>", fmt.Sprintf("%v", binding["id"]), -1)
	return template + "\n"
}
