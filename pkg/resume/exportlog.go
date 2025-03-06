package resume

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func WriteExportLog(configDir string, resource string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %v", err)
	}

	logFileName := filepath.Join(currentDir, configDir, "importlog.json")

	// Check if the log file exists
	_, err = os.Stat(logFileName)
	if os.IsNotExist(err) {
		// Create a new log file with the initial resource
		initialContent := map[string][]string{"resources": {resource}}
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
	var logData map[string][]string
	err = json.Unmarshal(fileContent, &logData)
	if err != nil {
		return fmt.Errorf("error unmarshaling JSON content: %v", err)
	}

	// Append the new resource to the resources array
	logData["resources"] = append(logData["resources"], resource)

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
	var logData map[string][]string
	err = json.Unmarshal(fileContent, &logData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON content: %v", err)
	}

	return logData["resources"], nil

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
