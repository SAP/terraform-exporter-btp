/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "from-file",
	Short: "export resources from a json file.",
	Long: `Use this command to export resources from the json file that is generated using the get command.
You can removes resource names from this config file, if you want to selectively import resources`,
	Run: func(cmd *cobra.Command, args []string) {
		resourceFileName, _ := cmd.Flags().GetString("resourceFileName")
		subaccount, _ := cmd.Flags().GetString("subaccount")
		configDir, _ := cmd.Flags().GetString("config-output-dir")
		jsonFile, _ := cmd.Flags().GetString("resource-file-path")

		exportFromFile(subaccount, jsonFile, resourceFileName, configDir)
	},
}

func init() {
	exportCmd.AddCommand(generateCmd)
	var subaccount string
	var resFile string
	var jsonFile string
	var configDir string

	generateCmd.Flags().StringVarP(&subaccount, "subaccount", "s", "", "Id of the subaccount")
	generateCmd.MarkFlagRequired("subaccount")
	generateCmd.Flags().StringVarP(&resFile, "resourceFileName", "f", "resources.tf", "filename for resource config generation")
	generateCmd.Flags().StringVarP(&jsonFile, "resource-file-path", "", "", "json file having subaccount resources info")
	generateCmd.MarkFlagRequired("json-file")
	generateCmd.Flags().StringVarP(&configDir, "config-output-dir", "o", "generated_configurations", "folder for config generation")

}
