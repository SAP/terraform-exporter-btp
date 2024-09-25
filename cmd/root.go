package cmd

import (
	"fmt"
	"os"

	output "github.com/SAP/terraform-exporter-btp/output"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
)

var debug bool
var subaccount string

var rootCmd = &cobra.Command{
	Use:               "btptf",
	Short:             "Terraform Exporter for SAP BTP",
	Long:              `btptf is a utility to generate Terraform configurations for existing SAP BTP resources that have been created manually and are not managed by Terraform.`,
	DisableAutoGenTag: true,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Display debugging output in the console.")
	rootCmd.PersistentFlags().StringVarP(&subaccount, "subaccount", "s", "", "Id of the subaccount")
	_ = rootCmd.MarkPersistentFlagRequired("subaccount")
	_ = viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	rootCmd.AddCommand(docCmd)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var docCmd = &cobra.Command{
	Use:    "gendoc",
	Short:  "Generate markdown documentation",
	Hidden: true, // Hide the command from the official CLI
	Run: func(cmd *cobra.Command, args []string) {
		docsDir := "./docs"
		if err := os.MkdirAll(docsDir, os.ModePerm); err != nil {
			errorMsg := output.ColorStringLightRed("error creating docs directory:")
			fmt.Println(errorMsg, err)
			os.Exit(1)
		}

		err := doc.GenMarkdownTree(rootCmd, docsDir)
		if err != nil {
			errorMsg := output.ColorStringLightRed("error generating documentation:")
			fmt.Println(errorMsg, err)
			os.Exit(1)
		}

		successMsg := output.ColorStringLightGreen("Documentation generated successfully in:")
		fmt.Println(successMsg, docsDir)
	},
}
