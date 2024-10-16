package cmd

import (
	"fmt"
	"log"
	"strings"

	files "github.com/SAP/terraform-exporter-btp/pkg/files"
	output "github.com/SAP/terraform-exporter-btp/pkg/output"
	tfimportprovider "github.com/SAP/terraform-exporter-btp/pkg/tfimportprovider"
	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
	"github.com/google/uuid"
)

const tfConfigFileName = "btp_resources.tf"
const configDirDefault = "generated_configurations"

func generateConfigForResource(resource string, values []string, subaccountId string, directoryId string, configDir string, resourceFileName string) {
	tempConfigDir := resource + "-config"

	level, iD := tfutils.GetExecutionLevelAndId(subaccountId, directoryId)

	importProvider, _ := tfimportprovider.GetImportBlockProvider(resource, level)
	resourceType := importProvider.GetResourceType()
	techResourceNameLong := strings.ToUpper(resourceType)

	tfutils.ExecPreExportSteps(tempConfigDir)

	output.AddNewLine()
	spinner := output.StartSpinner("crafting import block for " + techResourceNameLong)

	data, err := tfutils.FetchImportConfiguration(subaccountId, directoryId, resourceType, tfutils.TmpFolder)
	if err != nil {
		tfutils.CleanupProviderConfig(tempConfigDir)
		fmt.Print("\r\n")
		log.Fatalf("error fetching impport configuration for %s: %v", resourceType, err)
	}

	importBlock, err := importProvider.GetImportBlock(data, iD, values)
	if err != nil {
		tfutils.CleanupProviderConfig(tempConfigDir)
		fmt.Print("\r\n")
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
			fmt.Print("\r\n")
			log.Fatalf("error writing import configuration for %s: %v", resourceType, err)
		}

		output.StopSpinner(spinner)

		tfutils.ExecPostExportSteps(tempConfigDir, configDir, resourceFileName, techResourceNameLong)
	}
}

func isValidUuid(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func getUuidError(level string, iD string) string {
	if level == tfutils.SubaccountLevel {
		return fmt.Sprintf("Invalid subaccount ID: %s. Please provide a valid UUID.", iD)
	} else if level == tfutils.DirectoryLevel {
		return fmt.Sprintf("Invalid directory ID: %s Please provide a valid UUID.", iD)
	}
	return ""
}
