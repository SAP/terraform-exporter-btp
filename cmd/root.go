package cmd

import (
	"fmt"
	"os"

	output "github.com/SAP/terraform-exporter-btp/pkg/output"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
)

var verbose bool

var rootCmd = &cobra.Command{
	Use:               "btptf",
	Short:             "Terraform Exporter for SAP BTP",
	DisableAutoGenTag: true,
}

func init() {

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	templateOptions := generateCmdHelpOptions{
		Description: getRootCmdDescription,
	}

	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Display verbose output in the console for debugging.")
	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	rootCmd.AddCommand(docCmd)
	rootCmd.SetHelpTemplate(generateCmdHelp(rootCmd, templateOptions))
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
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

func getRootCmdDescription(c *cobra.Command) string {
	return generateCmdHelpDescription(c.Short,
		[]string{
			formatHelpNote(
				"Use the btptf command line tool (Terraform Exporter for SAP BTP) to generate Terraform configuration files for your SAP BTP resources. " +
					"These configuration files can then be used to manage SAP BTP resources with Terraform, adopting an Infrastructure-as-Code approach."),
			formatHelpNote(
				"The export is done based on subaccounts, so you need to specify the subaccount ID, and, optionally, a subset of resources inside the subaccount."),
			formatHelpNote(
				fmt.Sprintf("To help you specify the resources to be exported, you can first run the %s command to generate a .json file of all resources per subaccount. "+
					"This JSON file can then be used for the export with the %s command.",
					output.ColorStringCyan("create-json"),
					output.ColorStringCyan("export by-json"),
				)),
			formatHelpNote(
				fmt.Sprintf("Alternatively, you can directly specify the resources to be exported on the command line using %s.",
					output.ColorStringCyan("export by-resource"),
				)),
		})
}
