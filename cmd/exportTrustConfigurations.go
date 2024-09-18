package cmd

import (
	"btptfexport/output"
	"btptfexport/tfutils"
	"strings"

	"github.com/spf13/cobra"
)

// exportTrustConfigurationsCmd represents the exportSubaccountTrustConfigurations command
var exportTrustConfigurationsCmd = &cobra.Command{
	Use:               "trust-configurations",
	Short:             "export trust configurations of a subaccount",
	Long:              `export trust-configurations will export trust configurations of the given subaccount and generate resource configuration for it`,
	DisableAutoGenTag: true,
	Run: func(cmd *cobra.Command, args []string) {
		subaccount, _ := cmd.Flags().GetString("subaccount")
		resourceFileName, _ := cmd.Flags().GetString("resourceFileName")
		configDir, _ := cmd.Flags().GetString("config-output-dir")

		output.PrintExportStartMessage()
		tfutils.SetupConfigDir(configDir, true)
		exportSubaccountTrustConfigurations(subaccount, configDir, nil)
		tfutils.GenerateConfig(resourceFileName, configDir, true, strings.ToUpper(tfutils.SubaccountTrustConfigurationType))
		output.PrintExportSuccessMessage()
	},
}

func init() {
	exportCmd.AddCommand(exportTrustConfigurationsCmd)
	var subaccount string
	var resFile string
	var configDir string
	exportTrustConfigurationsCmd.Flags().StringVarP(&subaccount, "subaccount", "s", "", "Id of the subaccount")
	_ = exportTrustConfigurationsCmd.MarkFlagRequired("subaccount")
	exportTrustConfigurationsCmd.Flags().StringVarP(&resFile, "resourceFileName", "f", "resources.tf", "filename for resource config generation")
	exportTrustConfigurationsCmd.Flags().StringVarP(&configDir, "config-output-dir", "o", "generated_configurations", "folder for config generation")
}
