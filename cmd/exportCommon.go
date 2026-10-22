package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/SAP/terraform-exporter-btp/pkg/files"
	"github.com/SAP/terraform-exporter-btp/pkg/output"
	"github.com/SAP/terraform-exporter-btp/pkg/tfutils"
	"github.com/google/uuid"
)

const tfConfigFileName = "btp_resources.tf"
const configDirDefault = "generated_configurations_<account-id>"
const jsonFileDefault = "btpResources_<account-id>.json"

func generateConfigForResource(resource string, values []string, subaccountID, directoryID, organizationID, spaceID, configDir, resourceFileName string) (string, int) {
	tempConfigDir := resource + "-config"
	level, _ := tfutils.GetExecutionLevelAndId(subaccountID, directoryID, organizationID, spaceID)
	spec, err := tfutils.GetListResourceSpec(resource, func() string {
		if resource == tfutils.CmdCfSpaceRoleParameter {
			return tfutils.OrganizationLevel
		}
		return level
	}())
	if err != nil {
		log.Fatalf("error resolving list resource: %v", err)
	}

	tfutils.ExecPreExportSteps(tempConfigDir, level)
	if err := tfutils.ValidateTerraformQuerySupport(); err != nil {
		tfutils.CleanupProviderConfig(tempConfigDir)
		log.Fatalf("cannot query resources: %v", err)
	}
	if err := tfutils.PrepareQueryDirectory(tempConfigDir); err != nil {
		tfutils.CleanupProviderConfig(tempConfigDir)
		log.Fatalf("cannot prepare Terraform query: %v", err)
	}

	output.AddNewLine()
	spinner := output.StartSpinner("querying and generating " + strings.ToUpper(spec.ResourceType))
	resourceType, results, err := tfutils.QueryResource(tempConfigDir, resource, subaccountID, directoryID, organizationID, spaceID, values, true)
	if err != nil {
		tfutils.CleanupProviderConfig(tempConfigDir)
		fmt.Print("\r\n")
		log.Fatalf("error querying %s: %v", spec.ResourceType, err)
	}
	resourceType = tfutils.ExportResourceType(resourceType)
	if len(results) == 0 {
		output.StopSpinner(spinner)
		fmt.Println(output.ColorStringCyan("   no " + strings.ToUpper(resourceType) + " found"))
		tfutils.CleanupTempFiles(tempConfigDir)
		return resourceType, 0
	}
	if err := tfutils.WriteQueryResults(tempConfigDir, resourceFileName, resourceType, results); err != nil {
		tfutils.CleanupProviderConfig(tempConfigDir)
		log.Fatalf("error writing generated configuration for %s: %v", resourceType, err)
	}
	output.StopSpinner(spinner)
	tfutils.ExecPostQuerySteps(tempConfigDir, configDir, resourceFileName, strings.ToUpper(resourceType))
	return resourceType, len(results)
}

func isValidUuid(value string) bool { _, err := uuid.Parse(value); return err == nil }

func getUuidError(level, id string) string {
	switch level {
	case tfutils.SubaccountLevel:
		return fmt.Sprintf("Invalid subaccount ID: %s. Please provide a valid UUID.", id)
	case tfutils.DirectoryLevel:
		return fmt.Sprintf("Invalid directory ID: %s Please provide a valid UUID.", id)
	}
	return ""
}

func generateNextStepsDocument(configDir, subaccount, directory, organization, space string) {
	level, id := tfutils.GetExecutionLevelAndId(subaccount, directory, organization, space)
	level = strings.ReplaceAll(level, "Level", "")
	content := output.GetNextStepsTemplate(output.NextStepTemplateData{ConfigDir: configDir, UUID: id, Level: level})
	currentDir, err := os.Getwd()
	if err != nil {
		log.Println("error creating NextSteps.md file: ", err)
		return
	}
	if err := files.CreateFileWithContent(filepath.Join(currentDir, configDir, "NextSteps.md"), content); err != nil {
		log.Println("error creating NextSteps.md file: ", err)
	}
}
