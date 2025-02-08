package orchestrator

import (
	"log"
	"os"
	"path/filepath"

	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	providerprocessor "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/provider_processor"
	resourceprocessor "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/resource_processor"
)

func OrchestrateCodeCleanup(dir string) error {
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

	for _, file := range files {
		// We only process the resources and the provider files
		path := filepath.Join(dir, file.Name())

		if file.Name() == "btp_resources.tf" {
			f := generictools.GetHclFile(path)
			resourceprocessor.ProcessResources(f, &contentToCreate)
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
