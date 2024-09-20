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

func exportSubaccountEnvironmentInstances(subaccountID string, configFolder string, filterValues []string) {

	output.AddNewLine()
	spinner, err := output.StartSpinner("crafting import block for " + strings.ToUpper(tfutils.SubaccountEnvironmentInstanceType))
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	data, err := tfutils.FetchImportConfiguration(subaccountID, tfutils.SubaccountEnvironmentInstanceType, tfutils.TmpFolder)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	importBlock, err := getSubaccountEnvironmentInstanceBlock(data, subaccountID, filterValues)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	if len(importBlock) == 0 {
		log.Println("No environment instance found for the given subaccount")
		return
	}

	err = files.WriteImportConfiguration(configFolder, tfutils.SubaccountEnvironmentInstanceType, importBlock)
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

func getSubaccountEnvironmentInstanceBlock(data map[string]interface{}, subaccountId string, filterValues []string) (string, error) {

	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.SubaccountEnvironmentInstanceType)
	if err != nil {
		return "", err
	}

	var importBlock string
	environmentInstances := data["values"].([]interface{})

	if len(filterValues) != 0 {
		var subaccountAllEnvInstances []string

		for _, value := range environmentInstances {
			environmentInstance := value.(map[string]interface{})
			subaccountAllEnvInstances = append(subaccountAllEnvInstances, fmt.Sprintf("%v", environmentInstance["environment_type"]))
			if slices.Contains(filterValues, fmt.Sprintf("%v", environmentInstance["environment_type"])) {
				importBlock += templateEnvironmentInstanceImport(environmentInstance, subaccountId, resourceDoc)
			}
		}

		missingEnvInstance, subset := isSubset(subaccountAllEnvInstances, filterValues)

		if !subset {
			return "", fmt.Errorf("environment instance %s not found in the subaccount. Please adjust it in the provided file", missingEnvInstance)

		}
	} else {

		for _, value := range environmentInstances {
			environmentInstance := value.(map[string]interface{})
			importBlock += templateEnvironmentInstanceImport(environmentInstance, subaccountId, resourceDoc)
		}
	}

	return importBlock, nil
}

func templateEnvironmentInstanceImport(environmentInstance map[string]interface{}, subaccountId string, resourceDoc tfutils.EntityDocs) string {
	template := strings.Replace(resourceDoc.Import, "<resource_name>", fmt.Sprintf("%v", environmentInstance["environment_type"]), -1)
	template = strings.Replace(template, "<subaccount_id>", subaccountId, -1)
	template = strings.Replace(template, "<environment_instance_id>", fmt.Sprintf("%v", environmentInstance["id"]), -1)
	return template + "\n"
}
