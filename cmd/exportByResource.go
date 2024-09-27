package cmd

import (
	"github.com/SAP/terraform-exporter-btp/output"
	"github.com/SAP/terraform-exporter-btp/tfutils"

	"github.com/spf13/cobra"
)

// exportByResourceCmd represents the exportAll command
var exportByResourceCmd = &cobra.Command{
	Use:   "by-resource",
	Short: "export resources of a subaccount",
	Long: `by-resource command exports the resources of a subaccount as specified.

Examples:

btptf export by-resource --resources=subaccount,entitlements -s <subaccount-id>
btptf export by-resource --resources=all -s <subaccount-id> -p <file-name.json>

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

Mixing "all" with other resources will throw an error.`,
	DisableAutoGenTag: true,
	Run: func(cmd *cobra.Command, args []string) {
		subaccount, _ := cmd.InheritedFlags().GetString("subaccount")
		resourceFileName, _ := cmd.InheritedFlags().GetString("resource-file-name")
		configDir, _ := cmd.InheritedFlags().GetString("config-dir")

		resources, _ := cmd.Flags().GetString("resources")

		output.PrintExportStartMessage()
		tfutils.SetupConfigDir(configDir, true)

		resourcesList := tfutils.GetResourcesList(resources)
		for _, resourceToImport := range resourcesList {
			generateConfigForResource(resourceToImport, nil, subaccount, configDir, resourceFileName)
		}

		tfutils.FinalizeTfConfig(configDir)
		output.PrintExportSuccessMessage()
	},
}

func init() {
	var resources string
	exportByResourceCmd.Flags().StringVarP(&resources, "resources", "r", "all", "comma seperated string for resources")

	exportCmd.AddCommand(exportByResourceCmd)
}
