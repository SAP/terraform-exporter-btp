package cmd

import (
	"fmt"

	"github.com/SAP/terraform-exporter-btp/pkg/output"
	"github.com/SAP/terraform-exporter-btp/pkg/tfutils"

	"github.com/spf13/cobra"
)

// exportByResourceCmd represents the exportAll command
var exportByResourceCmd = &cobra.Command{
	Use:               "export",
	Short:             "Export resources of a subaccount",
	DisableAutoGenTag: true,
	Run: func(cmd *cobra.Command, args []string) {
		subaccount, _ := cmd.Flags().GetString("subaccount")
		configDir, _ := cmd.Flags().GetString("config-dir")
		resources, _ := cmd.Flags().GetString("resources")

		if configDir == configDirDefault {
			configDir = configDir + "_" + subaccount
		}

		output.PrintExportStartMessage()
		tfutils.SetupConfigDir(configDir, true)

		resourcesList := tfutils.GetResourcesList(resources)
		for _, resourceToImport := range resourcesList {
			generateConfigForResource(resourceToImport, nil, subaccount, configDir, tfConfigFileName)
		}

		tfutils.FinalizeTfConfig(configDir)
		tfutils.CleanupProviderConfig()
		output.PrintExportSuccessMessage()
	},
}

func init() {
	templateOptions := generateCmdHelpOptions{
		Description: getExportByResourceCmdDescription,
		Examples:    getExportByResourceCmdExamples,
	}

	var resources string
	var configDir string
	var subaccount string
	exportByResourceCmd.Flags().StringVarP(&subaccount, "subaccount", "s", "", "Id of the subaccount")
	_ = exportByResourceCmd.MarkFlagRequired("subaccount")
	exportByResourceCmd.Flags().StringVarP(&configDir, "config-dir", "c", configDirDefault, "folder for config generation")
	exportByResourceCmd.Flags().StringVarP(&resources, "resources", "r", "all", "comma seperated string for resources")

	rootCmd.AddCommand(exportByResourceCmd)

	exportByResourceCmd.SetUsageTemplate(generateCmdHelp(exportByResourceCmd, templateOptions))
	exportByResourceCmd.SetHelpTemplate(generateCmdHelp(exportByResourceCmd, templateOptions))
}

func getExportByResourceCmdDescription(c *cobra.Command) string {

	var resources string
	for i, resource := range tfutils.AllowedResources {
		if i == 0 {
			resources = output.ColorStringYellow(resource)
		} else {
			resources = resources + ", " + output.ColorStringYellow(resource)
		}
	}

	return generateCmdHelpDescription(c.Short,
		[]string{
			formatHelpNote(
				"Use this command to export SAP BTP resources specified by subaccount ID and, optionally, resource types.",
			),
			formatHelpNote(
				fmt.Sprintf("By default, the command will export all resources of a subaccount. "+
					"You can specify a subset with the %s flag.",
					output.ColorStringCyan("--resources"),
				)),
			formatHelpNote(
				fmt.Sprintf("Valid resources are: "+resources+" or %s (default)",
					output.ColorStringYellow("all"),
				)),
			formatHelpNote(
				fmt.Sprintf("Mixing %s with other resources will throw an error.",
					output.ColorStringYellow("all"),
				)),
		})
}

func getExportByResourceCmdExamples(c *cobra.Command) string {

	return generateCmdHelpCustomExamplesBlock(map[string]string{
		"Export a subaccount together with all its contained resources.": fmt.Sprintf("%s %s",
			output.ColorStringCyan("btptf export --subaccount"),
			output.ColorStringYellow("[Subaccount ID]"),
		),
		"Export a subaccount with entitlements only.": fmt.Sprintf("%s %s %s%s",
			output.ColorStringCyan("btptf export --subaccount"),
			output.ColorStringYellow("[Subaccount ID]"),
			output.ColorStringCyan("--resource="),
			output.ColorStringYellow("'subaccount,entitlements'"),
		),
	})
}
