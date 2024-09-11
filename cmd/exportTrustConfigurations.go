package cmd

import (
	"github.com/spf13/cobra"
)

// exportTrustConfigurationsCmd represents the exportTrustConfigurations command
var exportTrustConfigurationsCmd = &cobra.Command{
	Use:   "trust-configurations",
	Short: "export trust configurations of a subaccount",
	Long:  `export trust-configurations will export trust configurations of the given subaccount and generate resource configuration for it`,
	Run: func(cmd *cobra.Command, args []string) {
		subaccount, _ := cmd.Flags().GetString("subaccount")
		resourceFileName, _ := cmd.Flags().GetString("resourceFileName")
		configDir, _ := cmd.Flags().GetString("config-output-dir")
		setupConfigDir(configDir)
		exportTrustConfigurations(subaccount, configDir)
		generateConfig(resourceFileName, configDir)
	},
}

func init() {
	exportCmd.AddCommand(exportTrustConfigurationsCmd)
	var subaccount string
	var resFile string
	var configDir string
	exportTrustConfigurationsCmd.Flags().StringVarP(&subaccount, "subaccount", "s", "", "Id of the subaccount")
	exportTrustConfigurationsCmd.MarkFlagRequired("subaccount")
	exportTrustConfigurationsCmd.Flags().StringVarP(&resFile, "resourceFileName", "f", "resources.tf", "filename for resource config generation")
	exportTrustConfigurationsCmd.Flags().StringVarP(&configDir, "config-output-dir", "o", "generated_configurations", "folder for config generation")
}
