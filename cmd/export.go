/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/SAP/terraform-exporter-btp/output"
	"github.com/spf13/cobra"
)

var configDir string
var subaccount string

const tfConfigFileName = "btp_resources.tf"
const configDirDefault = "generated_configurations"

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export specific resources from an SAP BTP subaccount",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Invalid command\n\nUse 'btptf resource --help' for syntax instructions.")
	},
}

func init() {
	templateOptions := generateCmdHelpOptions{
		Description: getExportCmdDescription,
		Usage:       getExportCmdUsage,
	}

	exportCmd.PersistentFlags().StringVarP(&subaccount, "subaccount", "s", "", "Id of the subaccount")
	_ = exportCmd.MarkPersistentFlagRequired("subaccount")
	exportCmd.PersistentFlags().StringVarP(&configDir, "config-dir", "o", configDirDefault, "folder for config generation")
	exportCmd.SetUsageTemplate(generateCmdHelp(exportCmd, templateOptions))
	exportCmd.SetHelpTemplate(generateCmdHelp(exportCmd, templateOptions))
	rootCmd.AddCommand(exportCmd)
}

func getExportCmdDescription(c *cobra.Command) string {
	return generateCmdHelpDescription(c.Short,
		[]string{
			formatHelpNote(
				fmt.Sprintf("Use the %s commands to export resources specified by type",
					output.ColorStringCyan("by-resource"),
				)),
			formatHelpNote(
				fmt.Sprintf("Use the %s commands to export resources specified by type",
					output.ColorStringCyan("by-json"),
				)),
			formatHelpNote(
				fmt.Sprintf("The export creates a new directory (%s). "+
					"This directory contains the terraform configuration file (%s) and import files for each resource.",
					output.ColorStringYellow(fmt.Sprintf("%s_<UUID of parent>", configDirDefault)),
					output.ColorStringYellow(tfConfigFileName),
				)),
			formatHelpNote(
				fmt.Sprintf("You can change the default values for the directory by using the flag %s.",
					output.ColorStringCyan("--config-dir"),
				)),
		})
}

func getExportCmdUsage(*cobra.Command) string {
	return fmt.Sprintf("%s\n  %s\n\n",
		output.BoldString("Usage"), "btptf export by-json|by-resource [flags]")
}
