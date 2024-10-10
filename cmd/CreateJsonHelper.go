package cmd

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	files "github.com/SAP/terraform-exporter-btp/pkg/files"
	output "github.com/SAP/terraform-exporter-btp/pkg/output"
	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
)

func createJson(subaccount string, directory string, fileName string, resources []string) {
	if len(resources) == 0 {
		log.Fatal("please provide the btp resources you want to get using --resources flag or provide 'all' to get all resources")
	}

	tfutils.ConfigureProvider()

	spinner := output.StartSpinner("collecting resources")

	result, err := tfutils.ReadDataSources(subaccount, directory, resources)
	if err != nil {
		tfutils.CleanupProviderConfig()
		log.Fatalf("error reading data sources: %v", err)
	}

	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		tfutils.CleanupProviderConfig()
		log.Fatalf("error processing JSON of data source: %s", err)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		tfutils.CleanupProviderConfig()
		log.Fatalf("error getting current directory: %s", err)
	}

	dataBlockFile := filepath.Join(currentDir, fileName)
	err = files.CreateFileWithContent(dataBlockFile, string(jsonBytes))
	if err != nil {
		tfutils.CleanupProviderConfig()
		log.Fatalf("create file %s failed!", dataBlockFile)
	}

	tfutils.CleanupProviderConfig()

	output.StopSpinner(spinner)
}
