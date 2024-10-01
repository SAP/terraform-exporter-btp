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
	spinner := output.StartSpinner("crafting import block for " + techResourceNameLong)

	data, err := tfutils.FetchImportConfiguration(subaccountId, resourceType, tfutils.TmpFolder)
	if err != nil {
		tfutils.CleanupProviderConfig(tempConfigDir)
		log.Fatalf("error fetching impport configuration for %s: %v", resourceType, err)
	}

	importBlock, err := importProvider.GetImportBlock(data, subaccountId, values)
	if err != nil {
		tfutils.CleanupProviderConfig(tempConfigDir)
		log.Fatalf("error crafting import block: %v", err)
	}

	if len(importBlock) == 0 {
		output.StopSpinner(spinner)

		fmt.Println(output.ColorStringCyan("   no " + techResourceNameLong + " found for the given subaccount"))

		// Just clean up the temporary files, remaining setup remains untouched
		tfutils.CleanupTempFiles(tempConfigDir)
		fmt.Println(output.ColorStringGrey("   temporary files deleted"))

	} else {

		err = files.WriteImportConfiguration(tempConfigDir, resourceType, importBlock)
		if err != nil {
			tfutils.CleanupProviderConfig(tempConfigDir)
			log.Fatalf("error writing import configuration for %s: %v", resourceType, err)
		}

		output.StopSpinner(spinner)

		tfutils.ExecPostExportSteps(tempConfigDir, configDir, resourceFileName, techResourceNameLong)
	}
}
