package cmd

import (
	"btptfexport/tfutils"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func fetchImportConfiguration(subaccountID string, resourceType ResourceName, tmpFolder string) (map[string]interface{}, error) {
	dataBlock, err := readDataSource(subaccountID, resourceType)
	if err != nil {
		return nil, fmt.Errorf("error reading data source: %v", err)
	}

	dataBlockFile := filepath.Join(tmpFolder, "main.tf")
	err = tfutils.CreateFileWithContent(dataBlockFile, dataBlock)
	if err != nil {
		return nil, fmt.Errorf("create file %s failed: %v", dataBlockFile, err)
	}

	jsonBytes, err := GetTfStateData(tmpFolder, resourceType)
	if err != nil {
		return nil, fmt.Errorf("error getting Terraform state data: %v", err)
	}

	var data map[string]interface{}
	err = json.Unmarshal(jsonBytes, &data)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return data, nil
}

func writeImportConfiguration(configDir string, resourceType ResourceName, importBlock string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %v", err)
	}

	importFileName := fmt.Sprintf("%s_import.tf", resourceType)
	importFileName = filepath.Join(currentDir, configDir, importFileName)

	err = tfutils.CreateFileWithContent(importFileName, importBlock)
	if err != nil {
		return fmt.Errorf("create file %s failed: %v", importFileName, err)
	}

	log.Println(string(resourceType) + " exported. Please check " + configDir + " folder")

	return nil
}
