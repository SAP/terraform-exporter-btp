/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// exportCmd represents the export command
var resFile string
var configDir string

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export specific resources from an SAP BTP subaccount",
	Long: `
This command is used when you want to export resources of SAP BTP.

You have two options:

- by-json: export resources from a json file that is generated using the create-list command.
- by-resource: export resources you specify by type.

By default, the CLI it will generate the import files and a resource configuration file.
The directory for the configuration files has as default value 'generated_configurations'.
The resource configuration file has as default value 'btp_resources.tf'.

You can change the default values for the directory by using the flag --config-dir.
You can change the name of the resource configuration file by using the flag --resource-file-name.


The command will fail if a resource file already exists`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Invalid command\n\nUse 'btptf resource --help' for syntax instructions.")
	},
}

func init() {
	exportCmd.PersistentFlags().StringVarP(&resFile, "resource-file-name", "f", "btp_resources.tf", "filename for resource config generation")
	exportCmd.PersistentFlags().StringVarP(&configDir, "config-dir", "o", "generated_configurations", "folder for config generation")

	rootCmd.AddCommand(exportCmd)
}
