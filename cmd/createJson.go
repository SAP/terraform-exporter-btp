/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	output "github.com/SAP/terraform-exporter-btp/output"
	tfutils "github.com/SAP/terraform-exporter-btp/tfutils"

	"github.com/spf13/cobra"
)

// createJsonCmd represents the get command
var createJsonCmd = &cobra.Command{
	Use:               "create-json",
	Short:             "Store the list of resources in a subaccount into a JSON file",
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
	templateOptions := generateCmdHelpOptions{
		Description: getCreateJsonCmdDescription,
		Examples:    getCreateJsonCmdExamples,
	}

	var path string
	var resources string
	var subaccount string

	createJsonCmd.Flags().StringVarP(&subaccount, "subaccount", "s", "", "Id of the subaccount")
	_ = createJsonCmd.MarkFlagRequired("subaccount")
	createJsonCmd.Flags().StringVarP(&path, "path", "p", "btpResources.json", "path to JSON file with list of resources")
	createJsonCmd.Flags().StringVarP(&resources, "resources", "r", "all", "comma seperated string for resources")

	createJsonCmd.SetUsageTemplate(generateCmdHelp(createJsonCmd, templateOptions))
	createJsonCmd.SetHelpTemplate(generateCmdHelp(createJsonCmd, templateOptions))

	rootCmd.AddCommand(createJsonCmd)
}

func getCreateJsonCmdDescription(c *cobra.Command) string {

	var resources string
	for i, resource := range tfutils.AllowedResources {
		if i == 0 {
			resources = output.ColorStringYellow(resource)
		} else {
			resources = resources + ", " + output.ColorStringYellow(resource)
		}
	}

	return generateCmdHelpDescription(c.Short,
		[]string{
			formatHelpNote(
				"Use this command to compile a list of all resources in a subaccount and store it into a file",
			),
			formatHelpNote(
				fmt.Sprintf("Valid resources are: "+resources+" or %s (default)",
					output.ColorStringYellow("all"),
				)),
			formatHelpNote(
				fmt.Sprintf("Mixing %s with other resources will throw an error.",
					output.ColorStringYellow("all"),
				)),
		})
}

func getCreateJsonCmdExamples(c *cobra.Command) string {

	return generateCmdHelpCustomExamplesBlock(map[string]string{
		"Create a JSON file with all resources of a subaccount.": fmt.Sprintf("%s %s",
			output.ColorStringCyan("btptf create-json --subaccount"),
			output.ColorStringYellow("[Subaccount ID]"),
		),
		"Create a JSON file with resources 'subaccount' and 'entitlements' only.": fmt.Sprintf("%s%s %s %s",
			output.ColorStringCyan("btptf create-json --resources="),
			output.ColorStringYellow("'subaccount,entitlements'"),
			output.ColorStringCyan("--subaccount"),
			output.ColorStringYellow("[Subaccount ID]"),
		),
	})
}
