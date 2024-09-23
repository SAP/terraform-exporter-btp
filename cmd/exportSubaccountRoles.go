package cmd

import (
	"btptfexport/output"
	"btptfexport/tfutils"
	"strings"

	"github.com/spf13/cobra"
)

// exportSubaccountRolesCmd represents the exportSubaccountRoles command
var exportSubaccountRolesCmd = &cobra.Command{
	Use:               "roles",
	Short:             "export roles of a subaccount",
	Long:              `export roles will export all the roless of the given subaccount and generate resource configuration for it`,
	DisableAutoGenTag: true,
	Run: func(cmd *cobra.Command, args []string) {
		subaccount, _ := cmd.Flags().GetString("subaccount")
		resourceFileName, _ := cmd.Flags().GetString("resourceFileName")
		configDir, _ := cmd.Flags().GetString("config-output-dir")

		output.PrintExportStartMessage()
		tfutils.SetupConfigDir(configDir, true)
		exportSubaccountRoles(subaccount, configDir, nil)
		tfutils.GenerateConfig(resourceFileName, configDir, true, strings.ToUpper(tfutils.SubaccountRoleType))
		output.PrintExportSuccessMessage()
	},
}

func init() {
	exportCmd.AddCommand(exportSubaccountRolesCmd)
	var subaccount string
	var resourceFile string
	var configDir string
	exportSubaccountRolesCmd.Flags().StringVarP(&subaccount, "subaccount", "s", "", "Id of the subaccount")
	_ = exportSubaccountRolesCmd.MarkFlagRequired("subaccount")
	exportSubaccountRolesCmd.Flags().StringVarP(&resourceFile, "resourceFileName", "f", "resources.tf", "filename for resource config generation")
	exportSubaccountRolesCmd.Flags().StringVarP(&configDir, "config-output-dir", "o", "generated_configurations", "folder for config generation")
}
