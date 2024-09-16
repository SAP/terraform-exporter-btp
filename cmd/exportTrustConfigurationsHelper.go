package cmd

import (
	"btptfexport/tfutils"
	"fmt"
	"log"
	"slices"
	"strings"
)

func exportSubaccountTrustConfigurations(subaccountID string, configDir string, filterValues []string) {
	data, err := fetchImportConfiguration(subaccountID, SubaccountTrustConfigurationType, TmpFolder)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	importBlock, err := getTrustConfigurationsImportBlock(data, subaccountID, filterValues)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	if len(importBlock) == 0 {
		log.Println("No trust configuration found for the given subaccount")
		return
	}

	err = writeImportConfiguration(configDir, SubaccountTrustConfigurationType, importBlock)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}
}

func getTrustConfigurationsImportBlock(data map[string]interface{}, subaccountId string, filterValues []string) (string, error) {
	resourceDoc, err := getDocByResourceName(ResourcesKind, SubaccountTrustConfigurationType)
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
