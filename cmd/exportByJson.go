/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/SAP/terraform-exporter-btp/pkg/output"
	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"

	"github.com/spf13/cobra"
)

// exportByListCmd  represents the generate command
var exportByJsonCmd = &cobra.Command{
	Use:               "export-by-json",
	Short:             "Export resources from SAP BTP via JSON file",
	DisableAutoGenTag: true,
	Run: func(cmd *cobra.Command, args []string) {
		subaccount, _ := cmd.Flags().GetString("subaccount")
		directory, _ := cmd.Flags().GetString("directory")
		configDir, _ := cmd.Flags().GetString("config-dir")
		path, _ := cmd.Flags().GetString("path")

		level, iD := tfutils.GetExecutionLevelAndId(subaccount, directory)

		if !isValidUuid(iD) {
			log.Fatalln(getUuidError(level, iD))
		}

		if configDir == configDirDefault {
			configDir = configDir + "_" + iD
		}

		if path == jsonFileDefault {
			pathParts := strings.Split(path, ".")
			path = pathParts[0] + "_" + iD + "." + pathParts[1]
		}

		output.PrintExportStartMessage()
		exportByJson(subaccount, directory, path, tfConfigFileName, configDir)
		output.PrintExportSuccessMessage()
	},
}

func init() {
	templateOptions := generateCmdHelpOptions{
		Description:     getExportByJsonCmdDescription,
		DescriptionNote: getExportByJsonCmdDescriptionNote,
		Examples:        getExportByJsonCmdExamples,
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

	mainText := `Use this command to export resources from SAP BTP using a JSON file. The export is always per subaccount, directory, or environment instance. Create the JSON file with 'btptf create-json' and edit it as needed before exporting.`
	return generateCmdHelpDescription(mainText, nil)
}

func getExportByJsonCmdDescriptionNote(c *cobra.Command) string {
	return getSectionWithHeader("Note", "You must specify one of --subaccount, --directory, or --environment-instance.")
}

func getExportByJsonCmdExamples(c *cobra.Command) string {
	return generateCmdHelpCustomExamplesBlock(map[string]string{
		"Export the resources of a directory that are listed in a JSON with a custom file name and in a custom directory": fmt.Sprintf("%s %s %s %s",
			output.ColorStringCyan("btptf export-by-json --directory"),
			output.ColorStringYellow("[directory ID]"),
			output.ColorStringCyan("--path"),
			output.ColorStringYellow("'\\BTP\\resources\\my-btp-directory.json'"),
		),
		"Export the resources of a directory from JSON file from the default directory": fmt.Sprintf("%s %s",
			output.ColorStringCyan("btptf export-by-json --directory"),
			output.ColorStringYellow("[directory ID]"),
		),
		"Export the resources of a subaccount that are listed in a JSON file with a custom file name and in a custom directory": fmt.Sprintf("%s %s %s %s",
			output.ColorStringCyan("btptf export-by-json --subaccount"),
			output.ColorStringYellow("[subaccount ID]"),
			output.ColorStringCyan("--path"),
			output.ColorStringYellow("'\\BTP\\resources\\my-btp-subaccount.json'"),
		),
		"Export the resources of a subaccount that are listed in the JSON file from the default directory": fmt.Sprintf("%s %s",
			output.ColorStringCyan("btptf export-by-json --subaccount"),
			output.ColorStringYellow("[subaccount ID]"),
		),
	})
}
