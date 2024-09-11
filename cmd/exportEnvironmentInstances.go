package cmd

import (
	"github.com/spf13/cobra"
)

// exportEnvironmentInstancesCmd represents the exportEnvironmentInstances command
var exportEnvironmentInstancesCmd = &cobra.Command{
	Use:   "environment-instances",
	Short: "export environment instance of a subaccount",
	Long:  `exportEnvironmentInstance will export all the environment instance of the given subaccount and generate resource configuration for it`,
	Run: func(cmd *cobra.Command, args []string) {
		subaccount, _ := cmd.Flags().GetString("subaccount")
		resourceFileName, _ := cmd.Flags().GetString("resourceFileName")
		configDir, _ := cmd.Flags().GetString("config-output-dir")
		setupConfigDir(configDir)
		exportEnvironmentInstances(subaccount, configDir)
		generateConfig(resourceFileName, configDir)
	},
}

func init() {
	exportCmd.AddCommand(exportEnvironmentInstancesCmd)
	var subaccount string
	var resFile string
	var configDir string
	exportEnvironmentInstancesCmd.Flags().StringVarP(&subaccount, "subaccount", "s", "", "Id of the subaccount")
	exportEnvironmentInstancesCmd.MarkFlagRequired("subaccount")
	exportEnvironmentInstancesCmd.Flags().StringVarP(&resFile, "resourceFileName", "f", "resources.tf", "filename for resource config generation")
	exportEnvironmentInstancesCmd.Flags().StringVarP(&configDir, "config-output-dir", "o", "generated_configurations", "folder for config generation")
}
