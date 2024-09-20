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

func exportSubaccountTrustConfigurations(subaccountID string, configDir string, filterValues []string) {

	output.AddNewLine()
	spinner, err := output.StartSpinner("crafting import block for " + strings.ToUpper(tfutils.SubaccountTrustConfigurationType))
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	data, err := tfutils.FetchImportConfiguration(subaccountID, tfutils.SubaccountTrustConfigurationType, tfutils.TmpFolder)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	importBlock, err := getSubaccountTrustConfigurationsImportBlock(data, subaccountID, filterValues)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	if len(importBlock) == 0 {
		log.Println("No trust configuration found for the given subaccount")
		return
	}

	err = files.WriteImportConfiguration(configDir, tfutils.SubaccountTrustConfigurationType, importBlock)
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

func getSubaccountTrustConfigurationsImportBlock(data map[string]interface{}, subaccountId string, filterValues []string) (string, error) {

	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.SubaccountTrustConfigurationType)
	if err != nil {
		log.Fatalf("read doc failed!")
		return "", err
	}

	var importBlock string
	trusts := data["values"].([]interface{})

	if len(filterValues) != 0 {
		var subaccountAllTrusts []string

		for x, value := range trusts {
			trust := value.(map[string]interface{})
			subaccountAllTrusts = append(subaccountAllTrusts, fmt.Sprintf("%v", trust["origin"]))
			if slices.Contains(filterValues, fmt.Sprintf("%v", trust["origin"])) {
				importBlock += templateTrustImport(x, trust, subaccountId, resourceDoc)
			}
		}

		missingTrust, subset := isSubset(subaccountAllTrusts, filterValues)

		if !subset {
			return "", fmt.Errorf("trust configuration %s not found in the subaccount. Please adjust it in the provided file", missingTrust)
		}
	} else {
		for x, value := range trusts {
			trust := value.(map[string]interface{})
			importBlock += templateTrustImport(x, trust, subaccountId, resourceDoc)
		}
	}

	return importBlock, nil
}

func templateTrustImport(x int, trust map[string]interface{}, subaccountId string, resourceDoc tfutils.EntityDocs) string {
	template := strings.Replace(resourceDoc.Import, "<resource_name>", "trust"+fmt.Sprint(x), -1)
	template = strings.Replace(template, "<subaccount_id>", subaccountId, -1)
	template = strings.Replace(template, "<origin>", fmt.Sprintf("%v", trust["origin"]), -1)
	return template + "\n"
}
