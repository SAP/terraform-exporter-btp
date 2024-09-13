/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "resource",
	Short: "Export specific btp resources from a subaccount",
	Long: `
This command is used when you need to export specific resources.
By default, it will generate the <resource_name>_import.tf (import file) and resources.tf (resource file) files.
The resources.tf file can be renamed by using the flag --resourceFileName.
The command will fail if a resource file already exists`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Invalid command\n\nUse 'btptfexporter export --help' for syntax instructions.\n\nERROR")
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
