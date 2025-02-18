package testutils

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// Define source files (initial state, target state)
type testFileStruc struct {
	sourceFilePath string
	targetFilePath string
}

type testFileMapping map[string]testFileStruc

var testFiles = testFileMapping{
	"provider_btp": {
		sourceFilePath: "provider_btp_source.tf",
		targetFilePath: "provider_btp_target.tf",
	},
	"provider_cf": {
		sourceFilePath: "provider_cf_source.tf",
		targetFilePath: "provider_cf_target.tf",
	},
}

// Read a hcl files from disc
func GetHclFilesById(id string) (sourceHclFile *hclwrite.File, targetHclFile *hclwrite.File) {
	// Read hcl file from disc
	sourceFilePath, targetFilePath := getFilePathsbyId(id)

	currentDir, _ := os.Getwd()
	// We navigate up one level to land in the tfcleanup directory
	currentDir = filepath.Dir(currentDir)
	srcPath := filepath.Join(currentDir, "testutils", "testdata", sourceFilePath)
	trgtPath := filepath.Join(currentDir, "testutils", "testdata", targetFilePath)

	srcFile, err := os.ReadFile(srcPath)
	if err != nil {
		log.Printf("Failed to read file %q: %s", srcPath, err)
		return
	}

	trgtFile, err := os.ReadFile(trgtPath)
	if err != nil {
		log.Printf("Failed to read file %q: %s", trgtPath, err)
		return
	}

	sourceHclFile, diags := hclwrite.ParseConfig(srcFile, srcPath, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		for _, diag := range diags {
			if diag.Subject != nil {
				log.Printf("[%s:%d] %s: %s", diag.Subject.Filename, diag.Subject.Start.Line, diag.Summary, diag.Detail)
			} else {
				log.Printf("%s: %s", diag.Summary, diag.Detail)
			}
		}
		return
	}

	targetHclFile, diags = hclwrite.ParseConfig(trgtFile, trgtPath, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		for _, diag := range diags {
			if diag.Subject != nil {
				log.Printf("[%s:%d] %s: %s", diag.Subject.Filename, diag.Subject.Start.Line, diag.Summary, diag.Detail)
			} else {
				log.Printf("%s: %s", diag.Summary, diag.Detail)
			}
		}
		return
	}

	return sourceHclFile, targetHclFile
}

func AreHclFilesEqual(testResultHclFile *hclwrite.File, targetHclFile *hclwrite.File) error {

	if bytes.Equal(testResultHclFile.Bytes(), targetHclFile.Bytes()) {
		return nil
	} else {
		return fmt.Errorf("HCL files are not equal")
	}

}

func getFilePathsbyId(id string) (sourceFilePath string, targetFilePath string) {
	return testFiles[id].sourceFilePath, testFiles[id].targetFilePath
}
