package cmd

import (
	"fmt"
	"os"

	output "github.com/SAP/terraform-exporter-btp/pkg/output"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var verbose bool

const version = "1.0.0-rc1"

var rootCmd = &cobra.Command{
	Use:               "btptf",
	Short:             "Terraform Exporter for SAP BTP",
	DisableAutoGenTag: true,
	Version:           version,
}

func init() {

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	templateOptions := generateCmdHelpOptions{
		Description:     getRootCmdDescription,
		DescriptionNote: getRootCmdDescriptionNote,
	}

	rootCmd.SetVersionTemplate(fmt.Sprintf("Version: %s", version))
	rootCmd.Flags().BoolP("version", "v", false, "Print the version number of btptf")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Enable verbose output for debugging")
	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	rootCmd.AddCommand(docCmd)
	rootCmd.SetHelpTemplate(generateCmdHelp(rootCmd, templateOptions))
	rootCmd.SetUsageTemplate(generateCmdHelp(rootCmd, templateOptions))
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

		generateMarkdown(rootCmd)
	},
}

func getRootCmdDescription(c *cobra.Command) string {

	point1 := formatHelpNote("Directories")
	point2 := formatHelpNote("Subaccounts")
	point3 := formatHelpNote("Cloud Foundry environment instances")

	list := fmt.Sprintf("%s\n%s\n%s", point1, point2, point3)

	btptf := output.BoldString("btptf")

	description := `The Terraform Exporter for SAP BTP, or ` + btptf + `, exports existing SAP BTP resources as Terraform code, so you can start adopting Infrastructure-as-Code with Terraform.

The following SAP BTP account levels can be exported:
` + list + `

We recommend to start with 'btptf create-json', and then do the export with 'btptf export-by-json', as this lets you check and edit the resources before you export them.`

	return generateCmdHelpDescription(description, nil)
}

func getRootCmdDescriptionNote(c *cobra.Command) string {

	linkToRepo := output.AsLink("https://sap.github.io/terraform-exporter-btp/prerequisites/")

	point1 := formatHelpNote("To work with the btptf CLI, you need to configure authentication to access your global account on SAP BTP. For instructions, see " + linkToRepo + ".")
	point2 := formatHelpNote("To export resources, you need global account administrator permissions.")

	content := fmt.Sprintf("%s\n%s", point1, point2)

	return getSectionWithHeader("Prerequisites", content)
}
