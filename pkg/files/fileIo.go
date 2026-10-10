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

	// Avoid following a pre-existing symlink at the target path.
	if info, lerr := os.Lstat(fileName); lerr == nil && info.Mode()&os.ModeSymlink != 0 {
		_ = os.Remove(fileName)
	}

	// Create the file with owner-only permissions as it may hold credentials
	// or exported resource data.
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}

	defer func() {
		if tempErr := file.Close(); tempErr != nil {
			err = tempErr
		}
	}()

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

	defer func() {
		if tempErr := sourceFile.Close(); tempErr != nil {
			err = tempErr
		}
	}()

	destFile, err := os.OpenFile(dest, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer func() {
		if tempErr := destFile.Close(); tempErr != nil {
			err = tempErr
		}
	}()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func GetFullPath(inputPath string) string {
	normalizedInputPath := filepath.FromSlash(inputPath)
	fullPath, _ := filepath.Abs(normalizedInputPath)
	fullPath = filepath.Clean(fullPath)

	return fullPath
}
