package files

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const maxJsonFileSize = 5242880 // 5 MB

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

func WriteImportConfiguration(configDir *os.Root, resourceType string, importBlock string) error {
	if err := ValidateFilename(resourceType); err != nil {
		return err
	}
	return configDir.WriteFile(resourceType+"_import.tf", []byte(importBlock), 0600)
}

func CopyImportFiles(srcDir, destDir *os.Root) error {
	entries, err := fs.ReadDir(srcDir.FS(), ".")
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), "_import.tf") {
			continue
		}
		if err := CopyFile(srcDir, destDir, entry.Name(), os.O_APPEND); err != nil {
			return fmt.Errorf("copy import %s: %w", entry.Name(), err)
		}
	}
	return nil
}

// CopyFile reads and writes relative to retained directory handles. In
// particular, neither a filename nor a symlink can redirect either operation
// outside its directory.
func CopyFile(srcDir, destDir *os.Root, name string, flag int) error {
	if err := ValidateFilename(name); err != nil {
		return err
	}
	source, err := srcDir.Open(name)
	if err != nil {
		return err
	}
	defer func() { _ = source.Close() }()
	info, err := source.Stat()
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("%q is not a regular file", name)
	}
	destination, err := destDir.OpenFile(name, os.O_CREATE|os.O_WRONLY|flag, 0600)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(destination, source)
	closeErr := destination.Close()
	return errors.Join(copyErr, closeErr)
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

func GetFullPath(inputPath string) string {
	normalizedInputPath := filepath.FromSlash(inputPath)
	fullPath, _ := filepath.Abs(normalizedInputPath)
	fullPath = filepath.Clean(fullPath)

	return fullPath
}
