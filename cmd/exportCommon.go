package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/SAP/terraform-exporter-btp/pkg/output"
	"github.com/SAP/terraform-exporter-btp/pkg/tfutils"
	"github.com/google/uuid"
)

const tfConfigFileName = "btp_resources.tf"
const configDirDefault = "generated_configurations_<account-id>"
const jsonFileDefault = "btpResources_<account-id>.json"

func generateConfigForResource(resource string, values []string, subaccountID, directoryID, organizationID, spaceID string, configDir *os.Root, resourceFileName string, importsOnly bool) (string, int) {
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

	tempConfigDir, err := tfutils.NewQueryDirectory(level)
	if err != nil {
		tfutils.CleanupProviderConfig()
		log.Fatal(err)
	}
	if err := tfutils.ValidateTerraformQuerySupport(); err != nil {
		tfutils.CleanupProviderConfig(tempConfigDir)
		log.Fatalf("cannot query resources: %v", err)
	}
	if err := tfutils.PrepareQueryDirectory(tempConfigDir.Root); err != nil {
		tfutils.CleanupProviderConfig(tempConfigDir)
		log.Fatalf("cannot prepare Terraform query: %v", err)
	}

	output.AddNewLine()
	spinner := output.StartSpinner("querying and generating " + strings.ToUpper(spec.ResourceType))
	resourceType, results, err := tfutils.QueryResource(tempConfigDir.Root, resource, subaccountID, directoryID, organizationID, spaceID, values, true)
	if err != nil {
		tfutils.CleanupProviderConfig(tempConfigDir)
		fmt.Print("\r\n")
		log.Fatalf("error querying %s: %v", spec.ResourceType, err)
	}
	if len(results) == 0 {
		output.StopSpinner(spinner)
		fmt.Println(output.ColorStringCyan("   no " + strings.ToUpper(resourceType) + " found"))
		tfutils.CleanupTempFiles(tempConfigDir)
		return resourceType, 0
	}
	if err := tfutils.WriteQueryResults(tempConfigDir.Root, resourceFileName, resourceType, results); err != nil {
		tfutils.CleanupProviderConfig(tempConfigDir)
		log.Fatalf("error writing generated configuration for %s: %v", resourceType, err)
	}
	output.StopSpinner(spinner)
	tfutils.ExecPostQuerySteps(tempConfigDir, configDir, resourceFileName, strings.ToUpper(resourceType), importsOnly)
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

func generateNextStepsDocument(exportDir *os.Root, configDir, subaccount, directory, organization, space string) {
	level, id := tfutils.GetExecutionLevelAndId(subaccount, directory, organization, space)
	level = strings.ReplaceAll(level, "Level", "")
	content := output.GetNextStepsTemplate(output.NextStepTemplateData{ConfigDir: configDir, UUID: id, Level: level})
	if err := exportDir.WriteFile("NextSteps.md", []byte(content), 0600); err != nil {
		log.Println("error creating NextSteps.md file: ", err)
	}
}
