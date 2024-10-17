/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"path/filepath"
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
			configDirParts := strings.Split(configDir, "_")
			configDir = configDirParts[0] + "_" + configDirParts[1] + "_" + iD
		}

		if path == jsonFileDefault {
			pathParts := strings.Split(path, "_")
			path = pathParts[0] + "_" + iD + ".json"
		}

		output.PrintExportStartMessage()
		exportByJson(subaccount, directory, path, tfConfigFileName, configDir)
		output.PrintExportSuccessMessage()
	},
}

func init() {
	templateOptionsHelp := generateCmdHelpOptions{
		Description:     getExportByJsonCmdDescription,
		DescriptionNote: getExportByJsonCmdDescriptionNote,
		Examples:        getExportByJsonCmdExamples,
	}

	templateOptionsUsage := generateCmdHelpOptions{
		Description:     getEmtptySection,
		DescriptionNote: getEmtptySection,
		Examples:        getCreateJsonCmdExamples,
		Debugging:       getEmtptySection,
		Footer:          getEmtptySection,
	}

	var path string
	var configDir string
	var subaccount string
	var directory string

	exportByJsonCmd.Flags().StringVarP(&subaccount, "subaccount", "s", "", "ID of the subaccount")
	exportByJsonCmd.Flags().StringVarP(&directory, "directory", "d", "", "ID of the directory")
	exportByJsonCmd.MarkFlagsOneRequired("subaccount", "directory")
	exportByJsonCmd.MarkFlagsMutuallyExclusive("subaccount", "directory")

	exportByJsonCmd.Flags().StringVarP(&configDir, "config-dir", "c", configDirDefault, "Directory for the Terraform code")
	exportByJsonCmd.Flags().StringVarP(&path, "path", "p", jsonFileDefault, "Full path to JSON file with list of resources")

	rootCmd.AddCommand(exportByJsonCmd)

	exportByJsonCmd.SetHelpTemplate(generateCmdHelp(exportByJsonCmd, templateOptionsHelp))
	exportByJsonCmd.SetUsageTemplate(generateCmdHelp(exportByJsonCmd, templateOptionsUsage))
}

func getExportByJsonCmdDescription(c *cobra.Command) string {

	mainText := `Use this command to export resources from SAP BTP using a JSON file. The export is always per subaccount, directory, or environment instance. Create the JSON file with 'btptf create-json' and edit it as needed before exporting.`
	return generateCmdHelpDescription(mainText, nil)
}

func getExportByJsonCmdDescriptionNote(c *cobra.Command) string {
	return getSectionWithHeader("Note", "You must specify one of --subaccount, --directory, or --environment-instance.")
}

func getExportByJsonCmdExamples(c *cobra.Command) string {

	filePathSubaccount := filepath.Join("BTP", "resources", "my-btp-subaccount.json")
	filePathDirectory := filepath.Join("BTP", "resources", "my-btp-directory.json")

	return generateCmdHelpCustomExamplesBlock(map[string]string{
		"Export the resources of a directory that are listed in a JSON with a custom file name and in a custom directory": fmt.Sprintf("%s %s %s %s",
			output.ColorStringCyan("btptf export-by-json --directory"),
			output.ColorStringYellow("[directory ID]"),
			output.ColorStringCyan("--path"),
			output.ColorStringYellow("'"+filePathDirectory+"'"),
		),
		"Export the resources of a directory from JSON file from the default directory": fmt.Sprintf("%s %s",
			output.ColorStringCyan("btptf export-by-json --directory"),
			output.ColorStringYellow("[directory ID]"),
		),
		"Export the resources of a subaccount that are listed in a JSON file with a custom file name and in a custom directory": fmt.Sprintf("%s %s %s %s",
			output.ColorStringCyan("btptf export-by-json --subaccount"),
			output.ColorStringYellow("[subaccount ID]"),
			output.ColorStringCyan("--path"),
			output.ColorStringYellow("'"+filePathSubaccount+"'"),
		),
		"Export the resources of a subaccount that are listed in the JSON file from the default directory": fmt.Sprintf("%s %s",
			output.ColorStringCyan("btptf export-by-json --subaccount"),
			output.ColorStringYellow("[subaccount ID]"),
		),
	})
}
