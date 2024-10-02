package cmd

import (
	"fmt"

	"github.com/SAP/terraform-exporter-btp/output"
	"github.com/SAP/terraform-exporter-btp/tfutils"

	"github.com/spf13/cobra"
)

// exportByResourceCmd represents the exportAll command
var exportByResourceCmd = &cobra.Command{
	Use:               "by-resource",
	Short:             "Export resources of a subaccount",
	DisableAutoGenTag: true,
	Run: func(cmd *cobra.Command, args []string) {
		subaccount, _ := cmd.InheritedFlags().GetString("subaccount")
		resourceFileName, _ := cmd.InheritedFlags().GetString("resource-file-name")
		configDir, _ := cmd.InheritedFlags().GetString("config-dir")

		resources, _ := cmd.Flags().GetString("resources")

		output.PrintExportStartMessage()
		tfutils.SetupConfigDir(configDir, true)

		resourcesList := tfutils.GetResourcesList(resources)
		for _, resourceToImport := range resourcesList {
			generateConfigForResource(resourceToImport, nil, subaccount, configDir, resourceFileName)
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
	exportByResourceCmd.Flags().StringVarP(&resources, "resources", "r", "all", "comma seperated string for resources")

	exportByResourceCmd.SetUsageTemplate(generateCmdHelp(exportByResourceCmd, templateOptions))
	exportByResourceCmd.SetHelpTemplate(generateCmdHelp(exportByResourceCmd, templateOptions))

	exportCmd.AddCommand(exportByResourceCmd)
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
			output.ColorStringCyan("btptf export by-resource --subaccount"),
			output.ColorStringYellow("[Subaccount ID]"),
		),
		"Export a subaccount with entitlements only.": fmt.Sprintf("%s %s %s%s",
			output.ColorStringCyan("btptf export by-resource --subaccount"),
			output.ColorStringYellow("[Subaccount ID]"),
			output.ColorStringCyan("--resource="),
			output.ColorStringYellow("'subaccount,entitlements'"),
		),
	})
}
