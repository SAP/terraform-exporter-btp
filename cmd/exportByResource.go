package cmd

import (
	"fmt"
	"log"

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
		directory, _ := cmd.Flags().GetString("directory")
		configDir, _ := cmd.Flags().GetString("config-dir")
		resources, _ := cmd.Flags().GetString("resources")

		level, iD := tfutils.GetExecutionLevelAndId(subaccount, directory)

		if !isValidUuid(iD) {
			log.Fatalln(getUuidError(level, iD))
		}

		if configDir == configDirDefault {
			configDir = configDir + "_" + iD
		}

		output.PrintExportStartMessage()
		tfutils.SetupConfigDir(configDir, true)

		resourcesList := tfutils.GetResourcesList(resources, level)
		for _, resourceToImport := range resourcesList {
			generateConfigForResource(resourceToImport, nil, subaccount, directory, configDir, tfConfigFileName)
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
	var directory string

	exportByResourceCmd.Flags().StringVarP(&subaccount, "subaccount", "s", "", "ID of the subaccount")
	exportByResourceCmd.Flags().StringVarP(&directory, "directory", "d", "", "ID of the directory")
	exportByResourceCmd.MarkFlagsOneRequired("subaccount", "directory")
	exportByResourceCmd.MarkFlagsMutuallyExclusive("subaccount", "directory")

	exportByResourceCmd.Flags().StringVarP(&configDir, "config-dir", "c", configDirDefault, "folder for config generation")
	exportByResourceCmd.Flags().StringVarP(&resources, "resources", "r", "all", "comma seperated string for resources")

	rootCmd.AddCommand(exportByResourceCmd)

	exportByResourceCmd.SetHelpTemplate(generateCmdHelp(exportByResourceCmd, templateOptions))
}

func getExportByResourceCmdDescription(c *cobra.Command) string {

	var resources string
	for i, resource := range tfutils.AllowedResourcesSubaccount {
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
			output.ColorStringCyan("--resources="),
			output.ColorStringYellow("'subaccount,entitlements'"),
		),
		"Export a diretory together with all its contained resources.": fmt.Sprintf("%s %s",
			output.ColorStringCyan("btptf export --directory"),
			output.ColorStringYellow("[Directory ID]"),
		),
		"Export a directory with entitlements only.": fmt.Sprintf("%s %s %s%s",
			output.ColorStringCyan("btptf export --directory"),
			output.ColorStringYellow("[Directory ID]"),
			output.ColorStringCyan("--resources="),
			output.ColorStringYellow("'directory,entitlements'"),
		),
	})
}
