package cmd

import (
	"github.com/spf13/cobra"
)

// exportSubaccountCmd represents the exportSubaccount command
var subaccountCmd = &cobra.Command{
	Use:               "subaccount",
	Short:             "export subaccount",
	Long:              `export subaccount will export the given subaccount and generate resource configuration for it`,
	DisableAutoGenTag: true,
	Run: func(cmd *cobra.Command, args []string) {
		subaccount, _ := cmd.Flags().GetString("subaccount")
		resourceFileName, _ := cmd.Flags().GetString("resourceFileName")
		configDir, _ := cmd.Flags().GetString("config-output-dir")
		setupConfigDir(configDir)
		exportSubaccount(subaccount, configDir)
		generateConfig(resourceFileName, configDir)
	},
}

func init() {
	exportCmd.AddCommand(subaccountCmd)
	var subaccount string
	var resFile string
	var configDir string
	subaccountCmd.Flags().StringVarP(&subaccount, "subaccount", "s", "", "Id of the subaccount")
	_ = subaccountCmd.MarkFlagRequired("subaccount")
	subaccountCmd.Flags().StringVarP(&resFile, "resourceFileName", "f", "resources.tf", "filename for resource config generation")
	subaccountCmd.Flags().StringVarP(&configDir, "config-output-dir", "o", "generated_configurations", "folder for config generation")
}
