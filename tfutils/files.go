package tfutils

import (
	"fmt"
	"os"
)

func CreateFileWithContent(fileName string, content string) error {

	// Create the file
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("error:", err)
		return nil
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
