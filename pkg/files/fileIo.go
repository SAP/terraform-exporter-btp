package files

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

const maxJsonFileSize = 5242880 // 5 MB

func DeleteSourceFolder(srcDir string) {
	err := os.RemoveAll(srcDir)
	if err != nil {
		fmt.Print("\r\n")
		log.Fatalf("error deleting source folder %s: %v", srcDir, err)
	}
}

func CreateFileWithContent(fileName string, content string) error {

	// Create the file
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}
	defer file.Close()

	// Write content to the file
	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}
	return nil
}

func WriteImportConfiguration(configDir string, resourceType string, importBlock string) error {

	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %v", err)
	}

	importFileName := fmt.Sprintf("%s_import.tf", resourceType)
	importFileName = filepath.Join(currentDir, configDir, importFileName)

	err = CreateFileWithContent(importFileName, importBlock)
	if err != nil {
		return fmt.Errorf("create file %s failed: %v", importFileName, err)
	}

	return nil
}

func CopyImportFiles(srcDir, destDir string) error {
	// Find all files ending with "_import.tf" in the source directory
	files, err := filepath.Glob(filepath.Join(srcDir, "*_import.tf"))
	if err != nil {
		return fmt.Errorf("error finding files: %v", err)
	}

	// Copy each file to the destination directory
	for _, srcFile := range files {
		destFile := filepath.Join(destDir, filepath.Base(srcFile))

		err := copyFile(srcFile, destFile)
		if err != nil {
			return fmt.Errorf("error copying file %s to %s: %v", srcFile, destFile, err)
		}
	}
	return nil
}

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func IsFileSizeValid(filePath string) (bool, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		fmt.Print("\r\n")
		return false, fmt.Errorf("error reading JSON file: %v", err)
	}

	if fileInfo.Size() > maxJsonFileSize {
		return false, nil
	}
	return true, nil
}

func copyFile(src, dest string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.OpenFile(dest, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

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
