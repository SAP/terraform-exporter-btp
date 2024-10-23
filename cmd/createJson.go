/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"strings"

	output "github.com/SAP/terraform-exporter-btp/pkg/output"
	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"

	"github.com/spf13/cobra"
)

// createJsonCmd represents the get command
var createJsonCmd = &cobra.Command{
	Use:               "create-json",
	Short:             "Create a JSON file with a list of resources",
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

		if path == jsonFileDefault {
			pathParts := strings.Split(path, "_")
			path = pathParts[0] + "_" + iD + ".json"
		}

		output.PrintInventoryCreationStartMessage()
		resourcesList := tfutils.GetResourcesList(resources, level)
		createJson(subaccount, directory, path, resourcesList)
		output.PrintInventoryCreationSuccessMessage()
	},
}

func init() {
	templateOptionsHelp := generateCmdHelpOptions{
		Description:     getCreateJsonCmdDescription,
		DescriptionNote: getCreateJsonUsageNote,
		Examples:        getCreateJsonCmdExamples,
	}

	templateOptionsUsage := generateCmdHelpOptions{
		Description:     getEmtptySection,
		DescriptionNote: getEmtptySection,
		Examples:        getCreateJsonCmdExamples,
		Debugging:       getEmtptySection,
		Footer:          getEmtptySection,
	}

	var path string
	var resources string
	var subaccount string
	var directory string

	createJsonCmd.Flags().StringVarP(&subaccount, "subaccount", "s", "", "ID of the subaccount")
	createJsonCmd.Flags().StringVarP(&directory, "directory", "d", "", "ID of the directory")
	createJsonCmd.MarkFlagsOneRequired("subaccount", "directory")
	createJsonCmd.MarkFlagsMutuallyExclusive("subaccount", "directory")
	createJsonCmd.Flags().StringVarP(&path, "path", "p", jsonFileDefault, "Full path to JSON file with list of resources")
	createJsonCmd.Flags().StringVarP(&resources, "resources", "r", "all", "Comma-separated list of resources to be included")

	rootCmd.AddCommand(createJsonCmd)

	createJsonCmd.SetHelpTemplate(generateCmdHelp(createJsonCmd, templateOptionsHelp))
	createJsonCmd.SetUsageTemplate(generateCmdHelp(createJsonCmd, templateOptionsUsage))
}

func getCreateJsonCmdDescription(c *cobra.Command) string {

	var resources string
	for i, resource := range tfutils.AllowedResourcesSubaccount {
		if i == 0 {
			resources = resource
		} else {
			resources = resources + ", " + resource
		}
	}

	var resourcesDir string
	for i, resource := range tfutils.AllowedResourcesDirectory {
		if i == 0 {
			resourcesDir = resource
		} else {
			resourcesDir = resourcesDir + ", " + resource
		}
	}

	mainText := `Use this command to create a JSON file that lists all the resources for a directory, subaccount, or environment instance. This lets you easily edit the resources in the file before you export them.

Depending on the account level you specify, the JSON file will include the following resources:`

	return generateCmdHelpDescription(mainText,
		[]string{
			formatHelpNote(
				fmt.Sprint("For directories: " + resourcesDir),
			),
			formatHelpNote(
				fmt.Sprint("For subaccounts: " + resources),
			),
			formatHelpNote(
				"For environment instances: TBD",
			),
		})
}

func getCreateJsonUsageNote(c *cobra.Command) string {
	return getSectionWithHeader("Note", "You must specify one of --subaccount, --directory, or --environment-instance.")
}

func getCreateJsonCmdExamples(c *cobra.Command) string {
	return generateCmdHelpCustomExamplesBlock(map[string]string{
		"Create a JSON file for a directory with all of its resources": fmt.Sprintf("%s %s",
			output.ColorStringCyan("btptf create-json --directory"),
			output.ColorStringYellow("<directory ID>"),
		),
		"Create a JSON file for a subaccount with all of its resources": fmt.Sprintf("%s %s",
			output.ColorStringCyan("btptf create-json --subaccount"),
			output.ColorStringYellow("<subaccount ID>"),
		),
		"Create a JSON file for the entitlements of a subaccount": fmt.Sprintf("%s%s %s %s",
			output.ColorStringCyan("btptf create-json --resources="),
			output.ColorStringYellow("'subaccount,entitlements'"),
			output.ColorStringCyan("--subaccount"),
			output.ColorStringYellow("<subaccount ID>"),
		),

		"Create a JSON file for the roles and role collections of a subaccount": fmt.Sprintf("%s%s %s %s",
			output.ColorStringCyan("btptf create-json --resources="),
			output.ColorStringYellow("'roles,role-collections'"),
			output.ColorStringCyan("--subaccount"),
			output.ColorStringYellow("<subaccount ID>"),
		),
	})
}
