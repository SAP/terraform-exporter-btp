package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "btptfexporter",
	Short: "Terraform exporter for BTP",
	Long: `btptfexporter is a utility to generate configuration for existing btp resources that are created manually and not managed by terraform. btptfexporter help to generate configuration which then can be used by Terraform to bring that resource under terraform state.
	`,
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

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(docCmd)
}

var docCmd = &cobra.Command{
	Use:    "gendoc",
	Short:  "Generate markdown documentation",
	Hidden: true, // Hide the command from the official CLI
	Run: func(cmd *cobra.Command, args []string) {
		docsDir := "./docs"
		if err := os.MkdirAll(docsDir, os.ModePerm); err != nil {
			fmt.Println("Error creating docs directory:", err)
			os.Exit(1)
		}

		err := doc.GenMarkdownTree(rootCmd, docsDir)
		if err != nil {
			fmt.Println("Error generating documentation:", err)
			os.Exit(1)
		}

		fmt.Println("Documentation generated successfully in", docsDir)

	},
}
