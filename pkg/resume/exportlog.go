package resume

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
)

type Logentry struct {
	Resource     string `json:"resource"`
	ResourceType string `json:"resource_type"`
	Count        int    `json:"count"`
}

type Log struct {
	Resources []Logentry `json:"resources"`
}

func WriteExportLog(configDir string, resource string, resourceType string, count int) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %v", err)
	}

	logFileName := filepath.Join(currentDir, configDir, "importlog.json")

	// Check if the log file exists
	_, err = os.Stat(logFileName)
	if os.IsNotExist(err) {
		// Create a new log file with the initial resource
		initialContent := Log{
			Resources: []Logentry{
				{
					Resource:     resource,
					ResourceType: resourceType,
					Count:        count,
				},
			},
		}
		contentBytes, err := json.MarshalIndent(initialContent, "", "  ")
		if err != nil {
			return fmt.Errorf("error marshaling initial content: %v", err)
		}
		err = os.WriteFile(logFileName, contentBytes, 0644)
		if err != nil {
			return fmt.Errorf("error creating file %s: %v", logFileName, err)
		}
		return nil
	} else if err != nil {
		return fmt.Errorf("error checking file %s: %v", logFileName, err)
	}

	// Read the existing log file
	fileContent, err := os.ReadFile(logFileName)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", logFileName, err)
	}

	// Parse the existing JSON content
	var logData Log
	err = json.Unmarshal(fileContent, &logData)
	if err != nil {
		return fmt.Errorf("error unmarshaling JSON content: %v", err)
	}

	// Append the new resource to the resources array
	logData.Resources = append(logData.Resources, Logentry{
		Resource:     resource,
		ResourceType: resourceType,
		Count:        count,
	})

	// Marshal the updated content back to JSON
	updatedContent, err := json.MarshalIndent(logData, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling updated content: %v", err)
	}

	// Write the updated JSON back to the file
	err = os.WriteFile(logFileName, updatedContent, 0644)
	if err != nil {
		return fmt.Errorf("error writing to file %s: %v", logFileName, err)
	}

	return nil
}

func GetExistingExportLog(configDir string) ([]string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error getting current directory: %v", err)
	}

	logFileName := filepath.Join(currentDir, configDir, "importlog.json")

	// Check if the log file exists
	_, err = os.Stat(logFileName)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("error checking file %s: %v", logFileName, err)
	}

	// Read the existing log file
	fileContent, err := os.ReadFile(logFileName)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %v", logFileName, err)
	}

	// Parse the existing JSON content
	var logData Log
	err = json.Unmarshal(fileContent, &logData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON content: %v", err)
	}

	// Extract the resource names from the log
	var resources []string
	for _, entry := range logData.Resources {
		resources = append(resources, entry.Resource)
	}

	return resources, nil
}

func GetExistingExportLogComplete(configDir string) (Log, error) {
	var logData Log

	currentDir, err := os.Getwd()
	if err != nil {
		return logData, fmt.Errorf("error getting current directory: %v", err)
	}

	logFileName := filepath.Join(currentDir, configDir, "importlog.json")

	// Check if the log file exists
	_, err = os.Stat(logFileName)
	if os.IsNotExist(err) {
		return logData, nil
	} else if err != nil {
		return logData, fmt.Errorf("error checking file %s: %v", logFileName, err)
	}

	// Read the existing log file
	fileContent, err := os.ReadFile(logFileName)
	if err != nil {
		return logData, fmt.Errorf("error reading file %s: %v", logFileName, err)
	}

	// Parse the existing JSON content
	err = json.Unmarshal(fileContent, &logData)
	if err != nil {
		return logData, fmt.Errorf("error unmarshaling JSON content: %v", err)
	}

	return logData, nil
}

func RemoveExportLog(configDir string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %v", err)
	}

	logFileName := filepath.Join(currentDir, configDir, "importlog.json")

	// Check if the log file exists
	_, err = os.Stat(logFileName)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return fmt.Errorf("error checking file %s: %v", logFileName, err)
	}

	err = os.Remove(logFileName)
	if err != nil {
		return fmt.Errorf("error removing file %s: %v", logFileName, err)
	}

	return nil
}

func MergeSummaryTable(resultStore map[string]int, logData Log) map[string]int {
	// if there are no resources in the logData, return the resultStore
	if len(logData.Resources) == 0 {
		return resultStore
	}

	// We have resumed processing, so we must merge the logData with the resultStore
	resultStoreNew := make(map[string]int)
	for _, entry := range logData.Resources {
		resultStoreNew[entry.ResourceType] = entry.Count
	}

	maps.Copy(resultStoreNew, resultStore)

	return resultStoreNew
}
