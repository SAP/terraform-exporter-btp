/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"btptfexport/output"
	"btptfexport/tfutils"
	"strings"

	"github.com/spf13/cobra"
)

// exportSubaccountServiceBindingsCmd represents the exportSubaccountServiceBindings command
var exportSubaccountServiceBindingsCmd = &cobra.Command{
	Use:               "service-bindings",
	Short:             "export service bindings of a subaccount",
	Long:              `export service-bindings will export all the service bindings of the given subaccount and generate resource configuration for it`,
	DisableAutoGenTag: true,
	Run: func(cmd *cobra.Command, args []string) {
		subaccount, _ := cmd.Flags().GetString("subaccount")
		resourceFileName, _ := cmd.Flags().GetString("resourceFileName")
		configDir, _ := cmd.Flags().GetString("config-output-dir")

		output.PrintExportStartMessage()
		tfutils.SetupConfigDir(configDir, true)
		exportSubaccountServiceBindings(subaccount, configDir, nil)
		tfutils.GenerateConfig(resourceFileName, configDir, true, strings.ToUpper(tfutils.SubaccountServiceBindingType))
		output.PrintExportSuccessMessage()
	},
}

func init() {
	exportCmd.AddCommand(exportSubaccountServiceBindingsCmd)
	var subaccount string
	var resourceFile string
	var configDir string
	exportSubaccountServiceBindingsCmd.Flags().StringVarP(&subaccount, "subaccount", "s", "", "Id of the subaccount")
	_ = exportSubaccountServiceBindingsCmd.MarkFlagRequired("subaccount")
	exportSubaccountServiceBindingsCmd.Flags().StringVarP(&resourceFile, "resourceFileName", "f", "resources.tf", "filename for resource config generation")
	exportSubaccountServiceBindingsCmd.Flags().StringVarP(&configDir, "config-output-dir", "o", "generated_configurations", "folder for config generation")
}
