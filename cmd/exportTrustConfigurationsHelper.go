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

func exportTrustConfigurations(subaccountID string, configDir string, filterValues []string) {
	dataBlock, err := readDataSource(subaccountID, SubaccountTrustConfigurationType)
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

	jsonBytes, err := GetTfStateData(TmpFolder, SubaccountTrustConfigurationType)
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

	importBlock, err := getTrustConfigurationsImportBlock(data, subaccountID, filterValues)
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
