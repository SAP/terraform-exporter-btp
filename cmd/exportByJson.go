/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/SAP/terraform-exporter-btp/pkg/output"
	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"

	"github.com/spf13/cobra"
)

// exportByListCmd  represents the generate command
var exportByJsonCmd = &cobra.Command{
	Use:               "export-by-json",
	Short:             "Export resources based on a JSON file.",
	DisableAutoGenTag: true,
	Run: func(cmd *cobra.Command, args []string) {
		subaccount, _ := cmd.Flags().GetString("subaccount")
		directory, _ := cmd.Flags().GetString("directory")
		configDir, _ := cmd.Flags().GetString("config-dir")
		path, _ := cmd.Flags().GetString("path")

		_, iD := tfutils.GetExecutionLevelAndId(subaccount, directory)

		if configDir == configDirDefault {
			configDir = configDir + "_" + iD
		}

		output.PrintExportStartMessage()
		exportByJson(subaccount, directory, path, tfConfigFileName, configDir)
		output.PrintExportSuccessMessage()
	},
}

func init() {
	templateOptions := generateCmdHelpOptions{
		Description: getExportByJsonCmdDescription,
		Examples:    getExportByJsonCmdExamples,
	}

	var path string
	var configDir string
	var subaccount string
	var directory string

	exportByJsonCmd.Flags().StringVarP(&subaccount, "subaccount", "s", "", "ID of the subaccount")
	exportByJsonCmd.Flags().StringVarP(&directory, "directory", "d", "", "ID of the directory")
	exportByJsonCmd.MarkFlagsOneRequired("subaccount", "directory")
	exportByJsonCmd.MarkFlagsMutuallyExclusive("subaccount", "directory")

	exportByJsonCmd.Flags().StringVarP(&configDir, "config-dir", "c", configDirDefault, "folder for config generation")
	exportByJsonCmd.Flags().StringVarP(&path, "path", "p", "btpResources.json", "path to JSON file with list of resources")

	rootCmd.AddCommand(exportByJsonCmd)

	exportByJsonCmd.SetHelpTemplate(generateCmdHelp(exportByJsonCmd, templateOptions))
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
			output.ColorStringCyan("btptf export-by-json --subaccount"),
			output.ColorStringYellow("[Subaccount ID]"),
		),
		"Export the resources of a subaccount that are listed in a JSON file with a custom file name and in a custom directory.": fmt.Sprintf("%s %s %s %s",
			output.ColorStringCyan("btptf export-by-json --subaccount"),
			output.ColorStringYellow("[Subaccount ID]"),
			output.ColorStringCyan("--path"),
			output.ColorStringYellow("'\\BTP\\resources\\my-btp-resources.json'"),
		),
		"Export the resources of a directory that are listed in the JSON file from the default directory.": fmt.Sprintf("%s %s",
			output.ColorStringCyan("btptf export-by-json --sdirectory"),
			output.ColorStringYellow("[Directory ID]"),
		),
		"Export the resources of a directory that are listed in a JSON file with a custom file name and in a custom directory.": fmt.Sprintf("%s %s %s %s",
			output.ColorStringCyan("btptf export-by-json --directory"),
			output.ColorStringYellow("[Directory ID]"),
			output.ColorStringCyan("--path"),
			output.ColorStringYellow("'\\BTP\\resources\\my-btp-resources.json'"),
		),
	})
}
