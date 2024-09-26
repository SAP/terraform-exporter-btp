package cmd

import (
	"fmt"
	"log"
	"strings"

	files "github.com/SAP/terraform-exporter-btp/files"
	output "github.com/SAP/terraform-exporter-btp/output"
	tfimportprovider "github.com/SAP/terraform-exporter-btp/tfimportprovider"
	tfutils "github.com/SAP/terraform-exporter-btp/tfutils"
)

func generateConfigForResource(resource string, values []string, subaccountId string, configDir string, resourceFileName string) {
	tempConfigDir := resource + "-config"

	importProvider, _ := tfimportprovider.GetImportBlockProvider(resource)
	resourceType := importProvider.GetResourceType()
	techResourceNameLong := strings.ToUpper(resourceType)

	tfutils.ExecPreExportSteps(tempConfigDir)

	output.AddNewLine()
	spinner, err := output.StartSpinner("crafting import block for " + techResourceNameLong)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	data, err := tfutils.FetchImportConfiguration(subaccountId, resourceType, tfutils.TmpFolder)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	importBlock, err := importProvider.GetImportBlock(data, subaccountId, values)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	if len(importBlock) == 0 {
		err = output.StopSpinner(spinner)
		if err != nil {
			log.Fatalf("error: %v", err)
			return
		}

		fmt.Println(output.ColorStringCyan("No " + techResourceNameLong + "found for the given subaccount"))

		tfutils.CleanupTempFiles(tempConfigDir)

	} else {

		err = files.WriteImportConfiguration(configDir, resourceType, importBlock)
		if err != nil {
			log.Fatalf("error: %v", err)
			return
		}

		err = output.StopSpinner(spinner)
		if err != nil {
			log.Fatalf("error: %v", err)
			return
		}
		tfutils.ExecPostExportSteps(tempConfigDir, configDir, resourceFileName, techResourceNameLong)
	}

}
