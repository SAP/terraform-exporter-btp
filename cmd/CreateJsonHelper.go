package cmd

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	files "github.com/SAP/terraform-exporter-btp/files"
	output "github.com/SAP/terraform-exporter-btp/output"
	tfutils "github.com/SAP/terraform-exporter-btp/tfutils"
)

func createJson(subaccount string, fileName string, resources []string) {
	if len(resources) == 0 {
		log.Fatal("please provide the btp resources you want to get using --resources flag or provide 'all' to get all resources")
		return
	}

	tfutils.ConfigureProvider()

	spinner, err := output.StartSpinner("collecting resources")
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	result, err := tfutils.ReadDataSources(subaccount, resources)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("error: %s", err)
		return
	}

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("error getting current directory: %s", err)
		return
	}

	dataBlockFile := filepath.Join(currentDir, fileName)
	err = files.CreateFileWithContent(dataBlockFile, string(jsonBytes))
	if err != nil {
		log.Fatalf("create file %s failed!", dataBlockFile)
		return
	}

	err = output.StopSpinner(spinner)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

}
