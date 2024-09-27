/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	output "github.com/SAP/terraform-exporter-btp/output"
	tfutils "github.com/SAP/terraform-exporter-btp/tfutils"

	"github.com/spf13/cobra"
)

// createJsonCmd represents the get command
var createJsonCmd = &cobra.Command{
	Use:   "create-json",
	Short: "Store the list of resources in a subaccount into a JSON file",
	Long: `create-json command compiles a list of all resources in a subaccount and store it into a file.

Examples:

btptf create-json --resources=subaccount,entitlements -s <subaccount-id>
btptf create-json --resources=all -s <subaccount-id> -p <file-name.json>

Valid resources are:
- subaccount
- entitlements
- subscriptions
- environment-instances
- trust-configurations
- service-instances
- service-bindings
- roles
- role-collections

OR

- all

Mixing "all" with other resources will throw an error.
`,
	DisableAutoGenTag: true,
	Run: func(cmd *cobra.Command, args []string) {
		subaccount, _ := cmd.InheritedFlags().GetString("subaccount")
		path, _ := cmd.Flags().GetString("path")
		resources, _ := cmd.Flags().GetString("resources")

		output.PrintInventoryCreationStartMessage()
		resourcesList := tfutils.GetResourcesList(resources)
		createJson(subaccount, path, resourcesList)
		output.PrintInventoryCreationSuccessMessage()
	},
}

func init() {
	var path string
	var resources string
	createJsonCmd.Flags().StringVarP(&path, "path", "p", "btpResources.json", "path to JSON file with list of resources")
	createJsonCmd.Flags().StringVarP(&resources, "resources", "r", "all", "comma seperated string for resources")

	rootCmd.AddCommand(createJsonCmd)
}
