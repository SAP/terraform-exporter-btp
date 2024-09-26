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
		path, _ := cmd.Flags().GetString("path")

		output.PrintExportStartMessage()
		exportByJson(subaccount, path, resourceFileName, configDir)
		output.PrintExportSuccessMessage()
	},
}

func init() {
	var path string
	exportByJsonCmd.Flags().StringVarP(&path, "path", "p", "btpResources.json", "path to JSON file with list of resources")

	exportCmd.AddCommand(exportByJsonCmd)
}
