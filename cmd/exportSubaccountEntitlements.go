package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

// exportSubaccountEntitlementsCmd represents the exportSubaccountEntitlements command
var exportSubaccountEntitlementsCmd = &cobra.Command{
	Use:               "entitlements",
	Short:             "export entitlements of a subaccount",
	Long:              `export entitlements will export all the entitlements of the given subaccount and generate resource configuration for it`,
	DisableAutoGenTag: true,
	Run: func(cmd *cobra.Command, args []string) {
		subaccount, _ := cmd.Flags().GetString("subaccount")
		resourceFileName, _ := cmd.Flags().GetString("resourceFileName")
		configDir, _ := cmd.Flags().GetString("config-output-dir")

		printExportStartMessage()
		setupConfigDir(configDir, true)
		exportSubaccountEntitlements(subaccount, configDir, nil)
		generateConfig(resourceFileName, configDir, true, strings.ToUpper(string(SubaccountEntitlementType)))
		printExportSuccessMessage()
	},
}

func init() {
	exportCmd.AddCommand(exportSubaccountEntitlementsCmd)
	var subaccount string
	var resourceFile string
	var configDir string
	exportSubaccountEntitlementsCmd.Flags().StringVarP(&subaccount, "subaccount", "s", "", "Id of the subaccount")
	_ = exportSubaccountEntitlementsCmd.MarkFlagRequired("subaccount")
	exportSubaccountEntitlementsCmd.Flags().StringVarP(&resourceFile, "resourceFileName", "f", "resources.tf", "filename for resource config generation")
	exportSubaccountEntitlementsCmd.Flags().StringVarP(&configDir, "config-output-dir", "o", "generated_configurations", "folder for config generation")
}
