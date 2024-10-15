/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	output "github.com/SAP/terraform-exporter-btp/pkg/output"
	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"

	"github.com/spf13/cobra"
)

// createJsonCmd represents the get command
var createJsonCmd = &cobra.Command{
	Use:               "create-json",
	Short:             "Store the list of resources in a subaccount into a JSON file",
	DisableAutoGenTag: true,
	Run: func(cmd *cobra.Command, args []string) {
		subaccount, _ := cmd.Flags().GetString("subaccount")
		directory, _ := cmd.Flags().GetString("directory")
		path, _ := cmd.Flags().GetString("path")
		resources, _ := cmd.Flags().GetString("resources")

		level, iD := tfutils.GetExecutionLevelAndId(subaccount, directory)

		if !isValidUuid(iD) {
			log.Fatalln(getUuidError(level, iD))
		}

		output.PrintInventoryCreationStartMessage()
		resourcesList := tfutils.GetResourcesList(resources, level)
		createJson(subaccount, directory, path, resourcesList)
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
	var directory string

	createJsonCmd.Flags().StringVarP(&subaccount, "subaccount", "s", "", "ID of the subaccount")
	createJsonCmd.Flags().StringVarP(&directory, "directory", "d", "", "ID of the directory")
	createJsonCmd.MarkFlagsOneRequired("subaccount", "directory")
	createJsonCmd.MarkFlagsMutuallyExclusive("subaccount", "directory")
	createJsonCmd.Flags().StringVarP(&path, "path", "p", "btpResources.json", "path to JSON file with list of resources")
	createJsonCmd.Flags().StringVarP(&resources, "resources", "r", "all", "comma seperated string for resources")

	rootCmd.AddCommand(createJsonCmd)

	createJsonCmd.SetHelpTemplate(generateCmdHelp(createJsonCmd, templateOptions))
}

func getCreateJsonCmdDescription(c *cobra.Command) string {

	var resources string
	for i, resource := range tfutils.AllowedResourcesSubaccount {
		if i == 0 {
			resources = output.ColorStringYellow(resource)
		} else {
			resources = resources + ", " + output.ColorStringYellow(resource)
		}
	}

	var resourcesDir string
	for i, resource := range tfutils.AllowedResourcesDirectory {
		if i == 0 {
			resourcesDir = output.ColorStringYellow(resource)
		} else {
			resourcesDir = resourcesDir + ", " + output.ColorStringYellow(resource)
		}
	}

	return generateCmdHelpDescription(c.Short,
		[]string{
			formatHelpNote(
				"Use this command to compile a list of all resources in a subaccount and store it into a file",
			),
			formatHelpNote(
				fmt.Sprintf("Valid resources on subaccount level are: "+resources+" or %s (default)",
					output.ColorStringYellow("all"),
				)),
			formatHelpNote(
				fmt.Sprintf("Valid resources on directory level are: "+resourcesDir+" or %s (default)",
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
		"Create a JSON file with all resources of a directory.": fmt.Sprintf("%s %s",
			output.ColorStringCyan("btptf create-json --directory"),
			output.ColorStringYellow("[Directory ID]"),
		),
		"Create a JSON file with resources 'directory' and 'entitlements' on directory level only.": fmt.Sprintf("%s%s %s %s",
			output.ColorStringCyan("btptf create-json --resources="),
			output.ColorStringYellow("'directory,entitlements'"),
			output.ColorStringCyan("--directory"),
			output.ColorStringYellow("[Directory ID]"),
		),
	})
}
