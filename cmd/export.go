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
		//fmt.Println("please provide the resource to be imported with this commnad. Supported resources are subaccount, entilements, environment-instances, subscriptions, trust-configurations")
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
