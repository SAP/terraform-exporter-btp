package cmd

import (
	"btptfexport/output"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
)

var Debug bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "btptfexport",
	Short: "Terraform exporter for BTP",
	Long: `btptfexport is a utility to generate configuration for existing btp resources that are created manually and not managed by terraform. The CLI helps to generate configuration which then can be used by Terraform to bring that resource under terraform state.
	`,
	DisableAutoGenTag: true,
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
	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "Display debugging output in the console. (default: false)")
	_ = viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	rootCmd.AddCommand(docCmd)
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
