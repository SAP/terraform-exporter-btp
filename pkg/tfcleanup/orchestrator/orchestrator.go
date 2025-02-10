package orchestrator

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/SAP/terraform-exporter-btp/internal/btpcli"
	"github.com/SAP/terraform-exporter-btp/pkg/output"
	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	providerprocessor "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/provider_processor"
	resourceprocessor "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/resource_processor"
	"github.com/SAP/terraform-exporter-btp/pkg/tfutils"
)

func CleanUpGeneratedCode(configFolder string, level string) {
	if os.Getenv("BTPTF_EXPERIMENTAL") == "" {
		return
	}

	output.AddNewLine()
	spinner := output.StartSpinner("🧪 making the Terraform configuration even better")

	currentDir, err := os.Getwd()
	if err != nil {
		tfutils.CleanupProviderConfig()
		fmt.Print("\r\n")
		log.Fatalf("error getting current directory: %v", err)
	}

	terraformConfigPath := filepath.Join(currentDir, configFolder)

	err = orchestrateCodeCleanup(terraformConfigPath, level)

	if err != nil {
		fmt.Print("\r\n")
		log.Printf("error improving Terraform configuration: %v", err)
		log.Println("skipping improvement steps")
	}

	output.StopSpinner(spinner)
}

func orchestrateCodeCleanup(dir string, level string) error {
	dir = filepath.Clean(dir)

	_, err := os.Lstat(dir)
	if err != nil {
		log.Printf("Failed to stat %q: %s\n", dir, err)
		return err
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		log.Printf("Failed to read directory %q: %s", dir, err)
		return err
	}

	contentToCreate := make(generictools.VariableContent)
	dependencyAddresses := generictools.NewDepedendcyAddresses()

	btpClient, err := btpcli.GetLoggedInClient()

	if err != nil {
		return err
	}

	for _, file := range files {
		// We only process the resources and the provider files
		path := filepath.Join(dir, file.Name())

		if file.Name() == "btp_resources.tf" {
			f := generictools.GetHclFile(path)
			resourceprocessor.ProcessResources(f, level, &contentToCreate, &dependencyAddresses, btpClient)
			generictools.ProcessChanges(f, path)
		} else if file.Name() == "provider.tf" {
			f := generictools.GetHclFile(path)
			providerprocessor.ProcessProvider(f, &contentToCreate)
			generictools.ProcessChanges(f, path)
		}
	}

	if len(contentToCreate) > 0 {
		generictools.CreateVariablesFile(contentToCreate, dir)
	}
	return nil
}
