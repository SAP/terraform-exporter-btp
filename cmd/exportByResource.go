package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/SAP/terraform-exporter-btp/internal/cfcli"
	"github.com/SAP/terraform-exporter-btp/pkg/output"
	"github.com/SAP/terraform-exporter-btp/pkg/resume"
	tfcleantypes "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	tfcleanorchestrator "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/orchestrator"
	"github.com/SAP/terraform-exporter-btp/pkg/tfutils"

	"github.com/spf13/cobra"
)

// exportByResourceCmd represents the exportAll command
var exportByResourceCmd = &cobra.Command{
	Use:               "export",
	Short:             "Export resources from SAP BTP",
	DisableAutoGenTag: true,
	Run: func(cmd *cobra.Command, args []string) {
		subaccount, _ := cmd.Flags().GetString("subaccount")
		directory, _ := cmd.Flags().GetString("directory")
		organization, _ := cmd.Flags().GetString("organization")
		configDir, _ := cmd.Flags().GetString("config-dir")
		resources, _ := cmd.Flags().GetString("resources")

		backendPath, _ := cmd.Flags().GetString("backend-path")
		backendType, _ := cmd.Flags().GetString("backend-type")
		backendConfigOptions, _ := cmd.Flags().GetStringSlice("backend-config")
		backendConfig := tfutils.BackendConfig{
			PathToBackendConfig: backendPath,
			BackendType:         backendType,
			BackendConfig:       backendConfigOptions,
		}

		space := ""

		resultStore := make(map[string]int)

		level, iD := tfutils.GetExecutionLevelAndId(subaccount, directory, organization, space)

		if !isValidUuid(iD) {
			log.Fatalln(getUuidError(level, iD))
		}

		if configDir == configDirDefault {
			configDirParts := strings.Split(configDir, "_")
			configDir = configDirParts[0] + "_" + configDirParts[1] + "_" + iD
		}

		output.PrintExportStartMessage()
		tfutils.SetupConfigDir(configDir, true, level)

		resourcesList := tfutils.GetResourcesList(resources, level)

		exportLog, _ := resume.GetExistingExportLog(configDir)

		var fullExportLog resume.Log
		if len(exportLog) > 0 {
			resourcesList = tfutils.FilterResourcesListByLog(resourcesList, exportLog)
			fullExportLog, _ = resume.GetExistingExportLogComplete(configDir)
		}

		for _, resourceToImport := range resourcesList {
			if resourceToImport == tfutils.CmdCfSpaceRoleParameter {
				var finalCount int
				var resourceType string
				spaces, err := cfcli.GetSpaceList(organization)
				if err != nil {
					tfutils.CleanupProviderConfig()
					log.Fatalln(fmt.Errorf("unable to get space list for space role. err = %s", err))
				}
				for _, spaceID := range spaces {
					space := spaceID
					var count int
					resourceType, count = generateConfigForResource(resourceToImport, nil, subaccount, directory, organization, space, configDir, tfConfigFileName)
					finalCount = finalCount + count
				}
				resultStore[resourceType] = finalCount
				_ = resume.WriteExportLog(configDir, resourceToImport, resourceType, finalCount)

			} else {
				resourceType, count := generateConfigForResource(resourceToImport, nil, subaccount, directory, organization, space, configDir, tfConfigFileName)
				resultStore[resourceType] = count
				_ = resume.WriteExportLog(configDir, resourceToImport, resourceType, count)
			}
		}

		levelIds := tfcleantypes.LevelIds{
			SubaccountId: subaccount,
			DirectoryId:  directory,
			CfOrgId:      organization,
		}

		tfcleanorchestrator.CleanUpGeneratedCode(configDir, level, levelIds, &resultStore, backendConfig)
		tfutils.FinalizeTfConfig(configDir)
		generateNextStepsDocument(configDir, subaccount, directory, organization, space)
		tfutils.CleanupProviderConfig()
		_ = resume.RemoveExportLog(configDir)
		resultStoreNew := resume.MergeSummaryTable(resultStore, fullExportLog)
		output.RenderSummaryTable(resultStoreNew)
		output.PrintExportSuccessMessage()
	},
}

func init() {
	templateOptionsHelp := generateCmdHelpOptions{
		Description:     getExportByResourceCmdDescription,
		DescriptionNote: getExportCmdDescriptionNote,
		Examples:        getExportByResourceCmdExamples,
	}

	templateOptionsUsage := generateCmdHelpOptions{
		Description:     getEmtptySection,
		DescriptionNote: getEmtptySection,
		Examples:        getExportByResourceCmdExamples,
		Debugging:       getEmtptySection,
		Footer:          getEmtptySection,
	}

	var resources string
	var configDir string
	var subaccount string
	var directory string
	var organization string
	var backendPath string
	var backendType string

	exportByResourceCmd.Flags().StringVarP(&subaccount, "subaccount", "s", "", "ID of the subaccount")
	exportByResourceCmd.Flags().StringVarP(&directory, "directory", "d", "", "ID of the directory")
	exportByResourceCmd.Flags().StringVarP(&organization, "organization", "o", "", "ID of the Cloud Foundry organization")
	exportByResourceCmd.MarkFlagsOneRequired("subaccount", "directory", "organization")
	exportByResourceCmd.MarkFlagsMutuallyExclusive("subaccount", "directory", "organization")

	exportByResourceCmd.Flags().StringVarP(&configDir, "config-dir", "c", configDirDefault, "Directory for the Terraform code")
	exportByResourceCmd.Flags().StringVarP(&resources, "resources", "r", "all", "Comma-separated list of resources to be included")

	exportByResourceCmd.Flags().StringVarP(&backendPath, "backend-path", "b", "", "Path to the Terraform backend sample file")
	exportByResourceCmd.Flags().StringVar(&backendType, "backend-type", "", "Type of the Terraform backend")
	exportByResourceCmd.Flags().StringSlice("backend-config", []string{}, "Backend configuration")
	exportByResourceCmd.MarkFlagsMutuallyExclusive("backend-path", "backend-type")
	exportByResourceCmd.MarkFlagsRequiredTogether("backend-type", "backend-config")

	rootCmd.AddCommand(exportByResourceCmd)

	exportByResourceCmd.SetHelpTemplate(generateCmdHelp(exportByResourceCmd, templateOptionsHelp))
	exportByResourceCmd.SetUsageTemplate(generateCmdHelp(exportByResourceCmd, templateOptionsUsage))
}

func getExportByResourceCmdDescription(c *cobra.Command) string {

	var resources string
	for i, resource := range tfutils.AllowedResourcesSubaccount {
		if i == 0 {
			resources = resource
		} else {
			resources = resources + ", " + resource
		}
	}

	var resourcesDir string
	for i, resource := range tfutils.AllowedResourcesDirectory {
		if i == 0 {
			resourcesDir = resource
		} else {
			resourcesDir = resourcesDir + ", " + resource
		}
	}

	var resourcesEnv string
	for i, resource := range tfutils.AllowedResourcesOrganization {
		if i == 0 {
			resourcesEnv = resource
		} else {
			resourcesEnv = resourcesEnv + ", " + resource
		}
	}

	mainText := `Use this command to export resources from SAP BTP per account level (subaccount, directory, or environment instance). The command will create a directory with the Terraform configuration files and import blocks for the following resources in your specified account level:`
	return generateCmdHelpDescription(mainText,
		[]string{
			formatHelpNote(
				fmt.Sprint("For directories: " + resourcesDir),
			),
			formatHelpNote(
				fmt.Sprint("For subaccounts: " + resources),
			),
			formatHelpNote(
				"For environment instances: " + resourcesEnv,
			),
		})
}

func getExportCmdDescriptionNote(c *cobra.Command) string {
	point1 := formatHelpNote("We recommend to run this command only if you’re familiar with the Terraform resources in your SAP BTP accounts. For a safer approach, use 'btptf export-by-json'.")
	point2 := formatHelpNote("You must specify one of --subaccount, --directory, or --environment-instance.")

	content := fmt.Sprintf("%s\n%s", point1, point2)

	return getSectionWithHeader("Note", content)
}

func getExportByResourceCmdExamples(c *cobra.Command) string {

	return generateCmdHelpCustomExamplesBlock(map[string]string{
		"Export a directory that manages entitlements, but no users": fmt.Sprintf("%s %s %s %s",
			output.ColorStringCyan("btptf export --directory"),
			output.ColorStringYellow("<directory ID>"),
			output.ColorStringCyan("--resources"),
			output.ColorStringYellow("'directory,entitlements'"),
		),
		"Export a directory that doesn't manage entitlements or users": fmt.Sprintf("%s %s %s %s",
			output.ColorStringCyan("btptf export --directory"),
			output.ColorStringYellow("<directory ID>"),
			output.ColorStringCyan("--resources"),
			output.ColorStringYellow("'directory'"),
		),
		"Export a directory that manages entitlements and users": fmt.Sprintf("%s %s",
			output.ColorStringCyan("btptf export --directory"),
			output.ColorStringYellow("<directory ID>"),
		),
		"Export the entitlements of a subaccount": fmt.Sprintf("%s %s %s %s",
			output.ColorStringCyan("btptf export --subaccount"),
			output.ColorStringYellow("<subaccount ID>"),
			output.ColorStringCyan("--resources"),
			output.ColorStringYellow("'entitlements'"),
		),
		"Export the subscriptions of a subaccount": fmt.Sprintf("%s %s %s %s",
			output.ColorStringCyan("btptf export --subaccount"),
			output.ColorStringYellow("<subaccount ID>"),
			output.ColorStringCyan("--resources"),
			output.ColorStringYellow("'subscriptions'"),
		),
		"Export the roles and role collections of a subaccount": fmt.Sprintf("%s %s %s %s",
			output.ColorStringCyan("btptf export --subaccount"),
			output.ColorStringYellow("<subaccount ID>"),
			output.ColorStringCyan("--resources"),
			output.ColorStringYellow("'roles,role-collections'"),
		),
		"Export a subaccount with a sample backend configuration file": fmt.Sprintf("%s %s %s %s",
			output.ColorStringCyan("btptf export --subaccount"),
			output.ColorStringYellow("<subaccount ID>"),
			output.ColorStringCyan("--backend-path"),
			output.ColorStringYellow("backend.tf"),
		),
		"Export a subaccount with parameters for the backend configuration": fmt.Sprintf("%s %s %s %s %s %s %s %s",
			output.ColorStringCyan("btptf export --subaccount"),
			output.ColorStringYellow("<subaccount ID>"),
			output.ColorStringCyan("--backend-type"),
			output.ColorStringYellow("azurerm"),
			output.ColorStringCyan("--backend-config"),
			output.ColorStringYellow("'resource_group_name=rg-terraform-state'"),
			output.ColorStringCyan("--backend-config"),
			output.ColorStringYellow("'storage_account_name=terraformstatestorage'"),
		),
		"Export the spaces of a Cloud Foundry organization": fmt.Sprintf("%s %s %s %s",
			output.ColorStringCyan("btptf export --organization"),
			output.ColorStringYellow("<organization ID>"),
			output.ColorStringCyan("--resources"),
			output.ColorStringYellow("'spaces'"),
		),
		"Export the users of a Cloud Foundry organization": fmt.Sprintf("%s %s %s %s",
			output.ColorStringCyan("btptf export --organization"),
			output.ColorStringYellow("<organization ID>"),
			output.ColorStringCyan("--resources"),
			output.ColorStringYellow("'users'"),
		),
	})
}
