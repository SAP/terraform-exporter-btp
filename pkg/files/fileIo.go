package files

import (
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

	_, err = os.Stat(logFileName)

	if err == nil {
		// Read the file and append the resource
		file, err := os.OpenFile(logFileName, os.O_RDWR|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("error opening file %s: %v", logFileName, err)
		}
		defer file.Close()

		_, err = file.WriteString(fmt.Sprintf(", \"%s\"]}", resource))
		if err != nil {
			return fmt.Errorf("error writing to file %s: %v", logFileName, err)
		}
	} else if os.IsNotExist(err) {
		err = CreateFileWithContent(logFileName, fmt.Sprintf("{\"resources\": [\"%s\"]}", resource))
		if err != nil {
			return fmt.Errorf("create file %s failed: %v", logFileName, err)
		}
		return nil
	} else {
		return fmt.Errorf("error reading file %s: %v", logFileName, err)
	}

	return nil

}
