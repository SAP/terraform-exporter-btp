package orchestrator

import (
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/SAP/terraform-exporter-btp/internal/btpcli"
	"github.com/SAP/terraform-exporter-btp/pkg/output"
	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	providerprocessor "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/provider_processor"
	resourceprocessor "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/resource_processor"
	"github.com/SAP/terraform-exporter-btp/pkg/tfutils"
	"github.com/SAP/terraform-exporter-btp/pkg/toggles"
)

func CleanUpGeneratedCode(configFolder *os.Root, level string, levelIds generictools.LevelIds, resultStore *map[string]int, backendConfig tfutils.BackendConfig, importsOnly bool) {
	if toggles.IsCodeCleanupDeactivated() {
		return
	}

	output.AddNewLine()
	spinner := output.StartSpinner("📈 making the Terraform configuration even better")

	err := orchestrateCodeCleanup(configFolder, level, levelIds, resultStore, backendConfig, importsOnly)

	if err != nil {
		if importsOnly {
			tfutils.CleanupProviderConfig()
			log.Fatalf("error preparing imports-only configuration: %v", err)
		}
		fmt.Print("\r\n")
		log.Printf("error improving Terraform configuration: %v", err)
		log.Println("skipping improvement steps")
	}

	output.StopSpinner(spinner)
}

func orchestrateCodeCleanup(dir *os.Root, level string, levelIds generictools.LevelIds, resultStore *map[string]int, backendConfig tfutils.BackendConfig, importsOnly bool) error {
	entries, err := fs.ReadDir(dir.FS(), ".")
	if err != nil {
		return err
	}

	contentToCreate := make(generictools.VariableContent)
	dependencyAddresses := generictools.NewDependencyAddresses()

	var btpClient *btpcli.ClientFacade
	if !importsOnly {
		btpClient, err = btpcli.GetLoggedInClient()
		if err != nil {
			return err
		}
	}

	for _, file := range entries {
		// We only process the resources and the provider files
		path := file.Name()

		if file.Name() == "btp_resources.tf" {
			f, err := generictools.GetRootedHclFile(dir, path)
			if err != nil {
				return err
			}
			if importsOnly {
				if level == tfutils.SubaccountLevel && !toggles.IsEntitlementModuleGenerationDeactivated() {
					resourceprocessor.CollectEntitlementImports(f, &dependencyAddresses)
				}
			} else {
				resourceprocessor.ProcessResources(f, level, &contentToCreate, &dependencyAddresses, btpClient, levelIds)
				if err := generictools.ProcessChanges(dir, f, path); err != nil {
					return err
				}
			}
		} else if file.Name() == "provider.tf" {
			f, err := generictools.GetRootedHclFile(dir, path)
			if err != nil {
				return err
			}
			providerprocessor.ProcessProvider(f, &contentToCreate, backendConfig)
			if err := generictools.ProcessChanges(dir, f, path); err != nil {
				return err
			}
		}
	}

	// Remove unused imports
	if err := generictools.RemoveUnusedImports(dir, &dependencyAddresses.BlocksToRemove, resultStore); err != nil {
		return err
	}

	//Manipulate the import block to comprise the module addresses
	if err := resourceprocessor.AppendImportBlocksForEntitlementModule(dir, dependencyAddresses.EntitlementAddress, levelIds); err != nil {
		return err
	}

	err = generictools.RemoveEmptyFiles(dir)
	if err != nil {
		return err
	}

	if len(contentToCreate) > 0 {
		return generictools.CreateVariablesFile(contentToCreate, dir)
	}
	return nil
}
