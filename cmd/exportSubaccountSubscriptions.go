package cmd

import (
	"github.com/spf13/cobra"
)

// exportSubaccountSubscriptionCmd represents the exportSubaccountSubscription command
var exportSubaccountSubscriptionsCmd = &cobra.Command{
	Use:   "subscriptions",
	Short: "export subscriptions of a subaccount",
	Long:  `exportSubaccountSubscriptions will export subscriptions of the given subaccount and gererate resource configuration for it`,
	Run: func(cmd *cobra.Command, args []string) {
		subaccount, _ := cmd.Flags().GetString("subaccount")
		resourceFileName, _ := cmd.Flags().GetString("resourceFileName")
		configDir, _ := cmd.Flags().GetString("config-output-dir")
		setupConfigDir(configDir)
		exportSubaccountSubscriptions(subaccount, configDir)
		generateConfig(resourceFileName, configDir)
	},
}

func init() {
	exportCmd.AddCommand(exportSubaccountSubscriptionsCmd)
	var subaccount string
	var resFile string
	var configDir string
	exportSubaccountSubscriptionsCmd.Flags().StringVarP(&subaccount, "subaccount", "s", "", "Id of the subaccount")
	exportSubaccountSubscriptionsCmd.MarkFlagRequired("subaccount")
	exportSubaccountSubscriptionsCmd.Flags().StringVarP(&resFile, "resourceFileName", "f", "resources.tf", "filename for resource config generation")
	exportSubaccountSubscriptionsCmd.Flags().StringVarP(&configDir, "config-output-dir", "o", "generated_configurations", "folder for config generation")
}
