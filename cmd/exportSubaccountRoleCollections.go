package cmd

import (
	"btptfexport/output"
	"btptfexport/tfutils"
	"strings"

	"github.com/spf13/cobra"
)

// exportSubaccountRoleCollectionsCmd represents the exportSubaccountRoleCollections command
var exportSubaccountRoleCollectionsCmd = &cobra.Command{
	Use:               "role-collections",
	Short:             "export roles collections of a subaccount",
	Long:              `export role-collections will export all the role collections of the given subaccount and generate resource configuration for it`,
	DisableAutoGenTag: true,
	Run: func(cmd *cobra.Command, args []string) {
		subaccount, _ := cmd.Flags().GetString("subaccount")
		resourceFileName, _ := cmd.Flags().GetString("resourceFileName")
		configDir, _ := cmd.Flags().GetString("config-output-dir")

		output.PrintExportStartMessage()
		tfutils.SetupConfigDir(configDir, true)
		exportSubaccountRoleCollections(subaccount, configDir, nil)
		tfutils.GenerateConfig(resourceFileName, configDir, true, strings.ToUpper(tfutils.SubaccountRoleCollectionType))
		output.PrintExportSuccessMessage()
	},
}

func init() {
	exportCmd.AddCommand(exportSubaccountRoleCollectionsCmd)
	var subaccount string
	var resourceFile string
	var configDir string
	exportSubaccountRoleCollectionsCmd.Flags().StringVarP(&subaccount, "subaccount", "s", "", "Id of the subaccount")
	_ = exportSubaccountRoleCollectionsCmd.MarkFlagRequired("subaccount")
	exportSubaccountRoleCollectionsCmd.Flags().StringVarP(&resourceFile, "resourceFileName", "f", "resources.tf", "filename for resource config generation")
	exportSubaccountRoleCollectionsCmd.Flags().StringVarP(&configDir, "config-output-dir", "o", "generated_configurations", "folder for config generation")
}
