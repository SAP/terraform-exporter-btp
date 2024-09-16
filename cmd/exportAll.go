package cmd

import (
	"github.com/spf13/cobra"
)

// exportAllCmd represents the exportAll command
var exportAllCmd = &cobra.Command{
	Use:   "all",
	Short: "export all resources of a subaccount",
	Long: `export all will export all the resources from a subaccount. Currently only few resources are supported.

export all is a single command to export btp_subaccount, btp_subaccount_entitlements, btp_subaccount_instances, btp_subaccount_subscriptions,
btp_subaccount_trust_configurations `,
	DisableAutoGenTag: true,
	Run: func(cmd *cobra.Command, args []string) {
		subaccount, _ := cmd.Flags().GetString("subaccount")
		resourceFileName, _ := cmd.Flags().GetString("resourceFileName")
		configDir, _ := cmd.Flags().GetString("config-output-dir")

		setupConfigDir(configDir)

		for _, resourceToImport := range AllowedResources {
			generateConfigForResource(resourceToImport, nil, subaccount, configDir, resourceFileName)
		}

		finalizeTfConfig(configDir)
	},
}

func init() {
	exportCmd.AddCommand(exportAllCmd)
	var subaccount string
	var resFile string
	var configDir string
	exportAllCmd.Flags().StringVarP(&subaccount, "subaccount", "s", "", "Id of the subaccount")
	_ = exportAllCmd.MarkFlagRequired("subaccount")
	exportAllCmd.Flags().StringVarP(&resFile, "resourceFileName", "f", "resources.tf", "filename for resource config generation")
	exportAllCmd.Flags().StringVarP(&configDir, "config-output-dir", "o", "generated_configurations", "folder for config generation")
}
