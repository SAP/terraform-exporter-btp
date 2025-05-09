package output

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/SAP/terraform-exporter-btp/pkg/files"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/viper"
	"github.com/theckman/yacspin"
)

type NextStepTemplateData struct {
	ConfigDir string
	UUID      string
	Level     string
}

func createSpinner(message string) (*yacspin.Spinner, error) {
	cfg := yacspin.Config{
		Frequency:         100 * time.Millisecond,
		CharSet:           yacspin.CharSets[11],
		Suffix:            "  ", // puts a least one space between the animating spinner and the Message
		Message:           message,
		SuffixAutoColon:   true,
		ColorAll:          true,
		Colors:            []string{"fgYellow"},
		StopCharacter:     "✓",
		StopColors:        []string{"fgGreen"},
		StopMessage:       "done " + message,
		StopFailCharacter: "✗",
		StopFailColors:    []string{"fgRed"},
		StopFailMessage:   "failed",
	}

	s, err := yacspin.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to make spinner from struct: %w", err)
	}

	return s, nil
}

func stopOnSignal(spinner *yacspin.Spinner) {
	// ensure we stop the spinner before exiting, otherwise cursor will remain
	// hidden and terminal will require a `reset`
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh

		spinner.StopFailMessage("interrupted")

		// ignoring error intentionally
		_ = spinner.StopFail()

		os.Exit(0)
	}()
}

func renderSpinner(spinner *yacspin.Spinner) error {
	// start the spinner animation
	if err := spinner.Start(); err != nil {
		return fmt.Errorf("failed to start spinner: %w", err)
	}

	return nil
}

func StartSpinner(message string) *yacspin.Spinner {

	// No spinner execution during verbose mode
	verbose := viper.GetViper().GetBool("verbose")
	if verbose {
		return nil
	}

	spinner, err := createSpinner(message)
	if err != nil {
		slog.Warn(fmt.Sprintf("failed to make spinner from config struct: %v", err))
		return nil
	}

	stopOnSignal(spinner)

	err = renderSpinner(spinner)
	if err != nil {
		slog.Warn(err.Error())
		return nil
	}
	return spinner
}

func StopSpinner(spinner *yacspin.Spinner) {

	// No spinner execution during verbose mode
	verbose := viper.GetViper().GetBool("verbose")
	if verbose {
		return
	}

	if spinner == nil {
		return
	}

	if err := spinner.Stop(); err != nil {
		slog.Warn(fmt.Errorf("failed to stop spinner: %w", err).Error())
	}
}

func PrintExportStartMessage() {
	fmt.Println("")
	fmt.Println("🚀 Terraform configuration export started ...")
	fmt.Println("")
}

func PrintExportSuccessMessage(configDir string) {
	path2Config := files.GetFullPath(configDir)
	fileLink := makeFileLink(path2Config)

	fmt.Println("")
	fmt.Printf("🎉 Terraform configuration successfully created at %s\n", BoldString(configDir))
	fmt.Println("")
	fmt.Printf("Click here to navigate to the folder %s\n", AsLink(fileLink))
	fmt.Println("")
}

func PrintInventoryCreationStartMessage() {
	fmt.Println("")
	fmt.Println("🚀 Creation of resource list started ...")
	fmt.Println("")
}

func PrintInventoryCreationSuccessMessage(file string) {

	path2File := files.GetFullPath(file)
	folderPath := filepath.Dir(path2File)
	fileLink := makeFileLink(folderPath)

	fmt.Println("")
	fmt.Printf("📋 Resource list successfully created: %s\n", BoldString(file))
	fmt.Println("")
	fmt.Printf("Click here to navigate to the folder %s\n", AsLink(fileLink))
	fmt.Println("")
}

func RenderSummaryTable(data map[string]int) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"Resource Name", "Number of exported resources"})

	for key, value := range data {

		t.AppendRow(table.Row{key, value})

	}
	t.AppendSeparator()
	fmt.Println("")
	fmt.Println("📋 Export Summary")
	t.Render()
}

func ColorStringGrey(s string) string {
	return color.HiBlackString(s)
}

func ColorStringCyan(s string) string {
	return color.CyanString(s)
}

func ColorStringLightGreen(s string) string {
	return color.HiGreenString(s)
}

func ColorStringLightRed(s string) string {
	return color.HiRedString(s)
}

func ColorStringYellow(s string) string {
	return color.YellowString(s)
}

func AddNewLine() {
	fmt.Println("")
}

func BoldString(s string) string {
	return color.New(color.Bold).Sprint(s)
}

func AsLink(s string) string {
	return color.HiCyanString(s)
}

func GetNextStepsTemplate(input NextStepTemplateData) string {
	return fmt.Sprintf(`# How to Work With the Exported Configuration Files

You've successfully exported resources from a %s on SAP BTP using the Terraform exporter for SAP BTP (btptf CLI).

This export created Terraform configuration files and import blocks for your %s with ID %s in the %s folder. You'll need these files to run '*terraform apply*' and import the state.

At export, the generated code was refined by the btptf CLI as outlined in the [documentation](https://sap.github.io/terraform-exporter-btp/tfcodeimprovements/).

However, we strongly recommend that you review the code before you execute the state import.

Here are some points to consider:

1. **Check provider version constraints**: Check the version constraint in the provider configuration (*provider.tf*) i.e. make sure that the constraints are compliant with the rules of your company like cherry-picking one explicit version. We recommend to always use the latest version independent of the constraints you add.

2. **Cleanup configuration of resources**: The configuration (*btp_resources.tf*) is generated based on the information about the resources available from the provider plugin. All data including optional data that got defaulted (e.g. usage in the btp_subaccount resource) is added to the configuration. To reduce the amount of data you could remove optional attributes that you don't want to have set explicitly.

3. **Declare variables**: The generated code already contains some variables in the *variables.tf* file. Depending on your requirements you might want to add further parameters to the variable list. For example, the name of a subaccount.

4. **Configure backend**: The state of your configuration should be stored in a remote state backend. If you have not injected an existing remote state at export (see [How to Add a Remote Backend Configuration](https://sap.github.io/terraform-exporter-btp/remotebackend/)), make sure to add the corresponding configuration in the *provider.tf* file manually. You find more details in the [Terraform documentation](https://developer.hashicorp.com/terraform/language/backend).

5. **Validate the import**: Validate that the import is possible by executing '*terraform plan*'. Depending on the number of resources, the planing should return a message like this:

    Plan: n to import, 0 to add, 0 to change, 0 to destroy.

Now you're all set to run '*terraform apply*', which will import the state and thus bring your SAP BTP resources under the management of Terraform. Congrats!

`, input.Level, input.Level, input.UUID, input.ConfigDir)
}

func makeFileLink(path string) (fileLink string) {
	if runtime.GOOS == "windows" {
		// Use double backslashes for Windows paths
		fileLink = fmt.Sprintf("file:///%s", filepath.ToSlash(path))
	} else {
		// Use standard file URL for Unix-based systems
		fileLink = fmt.Sprintf("file://%s", path)
	}
	return fileLink
}
