/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/SAP/terraform-exporter-btp/pkg/output"

	"github.com/spf13/cobra"
)

// exportByListCmd  represents the generate command
var exportByJsonCmd = &cobra.Command{
	Use:               "by-json",
	Short:             "Export resources based on a JSON file.",
	DisableAutoGenTag: true,
	Run: func(cmd *cobra.Command, args []string) {
		subaccount, _ := cmd.InheritedFlags().GetString("subaccount")
		configDir, _ := cmd.InheritedFlags().GetString("config-dir")
		path, _ := cmd.Flags().GetString("path")

		if configDir == configDirDefault {
			configDir = configDir + "_" + subaccount
		}
		output.PrintExportStartMessage()
		exportByJson(subaccount, path, tfConfigFileName, configDir)
		output.PrintExportSuccessMessage()
	},
}

func init() {
	templateOptions := generateCmdHelpOptions{
		Description: getExportByJsonCmdDescription,
		Examples:    getExportByJsonCmdExamples,
	}

	var path string
	exportByJsonCmd.Flags().StringVarP(&path, "path", "p", "btpResources.json", "path to JSON file with list of resources")
	exportByJsonCmd.SetUsageTemplate(generateCmdHelp(exportByJsonCmd, templateOptions))
	exportByJsonCmd.SetHelpTemplate(generateCmdHelp(exportByJsonCmd, templateOptions))
	exportCmd.AddCommand(exportByJsonCmd)
}

func getExportByJsonCmdDescription(c *cobra.Command) string {

	return generateCmdHelpDescription(c.Short,
		[]string{
			formatHelpNote(
				"Use this command to export resources from the JSON file.",
			),
			formatHelpNote(
				fmt.Sprintf("You create the JSON file via the %s command.",
					output.ColorStringCyan("create-json"),
				)),
		})
}

func getExportByJsonCmdExamples(c *cobra.Command) string {

	return generateCmdHelpCustomExamplesBlock(map[string]string{
		"Export the resources of a subaccount that are listed in the JSON file from the default directory.": fmt.Sprintf("%s %s",
			output.ColorStringCyan("btptf export by-json --subaccount"),
			output.ColorStringYellow("[Subaccount ID]"),
		),
		"Export the resources of a subaccount that are listed in a JSON file with a custom file name and in a custom directory.": fmt.Sprintf("%s %s %s %s",
			output.ColorStringCyan("btptf export by-json --subaccount"),
			output.ColorStringYellow("[Subaccount ID]"),
			output.ColorStringCyan("--path"),
			output.ColorStringYellow("'\\BTP\\resources\\my-btp-resources.json'"),
		),
	})
}
