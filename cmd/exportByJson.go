/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/SAP/terraform-exporter-btp/output"

	"github.com/spf13/cobra"
)

// exportByListCmd  represents the generate command
var exportByJsonCmd = &cobra.Command{
	Use:               "by-json",
	Short:             "export resources based on a JSON file.",
	Long:              `Use this command to export resources from the JSON file that is generated using the create-json command.`,
	DisableAutoGenTag: true,
	Run: func(cmd *cobra.Command, args []string) {
		subaccount, _ := cmd.InheritedFlags().GetString("subaccount")
		configDir, _ := cmd.InheritedFlags().GetString("config-dir")
		resourceFileName, _ := cmd.InheritedFlags().GetString("resource-file-name")
		jsonFile, _ := cmd.Flags().GetString("from")

		output.PrintExportStartMessage()
		exportByJson(subaccount, jsonFile, resourceFileName, configDir)
		output.PrintExportSuccessMessage()
	},
}

func init() {
	var jsonFile string
	exportByJsonCmd.Flags().StringVarP(&jsonFile, "from", "p", "btpResources.json", "path to JSON file with resources")

	exportCmd.AddCommand(exportByJsonCmd)
}
