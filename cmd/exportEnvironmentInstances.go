package cmd

import (
	"btptfexport/tfutils"
	"strings"

	"github.com/spf13/cobra"
)

// exportEnvironmentInstancesCmd represents the exportSubaccountEnvironmentInstances command
var exportEnvironmentInstancesCmd = &cobra.Command{
	Use:               "environment-instances",
	Short:             "export environment instance of a subaccount",
	Long:              `export environment-instance will export all the environment instance of the given subaccount and generate resource configuration for it`,
	DisableAutoGenTag: true,
	Run: func(cmd *cobra.Command, args []string) {
		subaccount, _ := cmd.Flags().GetString("subaccount")
		resourceFileName, _ := cmd.Flags().GetString("resourceFileName")
		configDir, _ := cmd.Flags().GetString("config-output-dir")
		tfutils.SetupConfigDir(configDir, true)
		exportSubaccountEnvironmentInstances(subaccount, configDir, nil)
		tfutils.GenerateConfig(resourceFileName, configDir, true, strings.ToUpper(tfutils.SubaccountEnvironmentInstanceType))
	},
}

func init() {
	exportCmd.AddCommand(exportEnvironmentInstancesCmd)
	var subaccount string
	var resFile string
	var configDir string
	exportEnvironmentInstancesCmd.Flags().StringVarP(&subaccount, "subaccount", "s", "", "Id of the subaccount")
	_ = exportEnvironmentInstancesCmd.MarkFlagRequired("subaccount")
	exportEnvironmentInstancesCmd.Flags().StringVarP(&resFile, "resourceFileName", "f", "resources.tf", "filename for resource config generation")
	exportEnvironmentInstancesCmd.Flags().StringVarP(&configDir, "config-output-dir", "o", "generated_configurations", "folder for config generation")
}
