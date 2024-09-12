/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "generate-resources-list",
	Short: "Store the list of resources in a subaccount into a json file",
	Long: `generate-resources-list command will get all the resource list or specified resource list in a subaccount.
It will then store this list into a file.

For example:

btptfexport generate-resources-list --resources=subaccount,entitlements -s <subaccount-id>
btptfexport generate-resources-list --resources=all -s <subaccount-id> -j <file-name.json>

Valid resources are:
- subaccount
- entitlements
- subscriptions
- environment-instances
- trust-configurations

OR

- all

Mixing "all" with other resources will throw an error.
`,
	DisableAutoGenTag: true,
	Run: func(cmd *cobra.Command, args []string) {
		subaccount, _ := cmd.Flags().GetString("subaccount")
		fileName, _ := cmd.Flags().GetString("json-out")
		resources, _ := cmd.Flags().GetString("resources")
		getResourcesInfo(subaccount, fileName, resources)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	var subaccount string
	var fileName string
	var resources string
	getCmd.Flags().StringVarP(&subaccount, "subaccount", "s", "", "Id of the subaccount")
	_ = getCmd.MarkFlagRequired("subaccount")
	getCmd.Flags().StringVarP(&fileName, "json-out", "j", "btpResources.json", "json file for list of resources")
	getCmd.Flags().StringVarP(&resources, "resources", "r", "all", "comma seperated string for resources")

}
