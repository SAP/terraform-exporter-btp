package output

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"
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

func PrintExportSuccessMessage() {
	fmt.Println("")
	fmt.Println("🎉 Terraform configuration successfully created")
	fmt.Println("")
}

func PrintInventoryCreationStartMessage() {
	fmt.Println("")
	fmt.Println("🚀 Creation of resource list started ...")
	fmt.Println("")
}

func PrintInventoryCreationSuccessMessage() {
	fmt.Println("")
	fmt.Println("📋 Resource list successfully created")
	fmt.Println("")
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
	return fmt.Sprintf(`## How to Work With the Exported Configuration Files

You've successfully exported resources from a %s on SAP BTP using the btptf CLI.

This created Terraform configuration files and import blocks for your %s with ID %s in the %s folder. You'll need these files to run '*terraform apply*'.

But you should first review the generated code:

1. Check provider version constraints
   Check the version constraint in the provider configuration (*provider.tf*) i.e. make sure that the constraints are compliant with the rules of your company like cherry-picking one explicit version. We recommend to always use the latest version independent of the constraints you add.

2. Cleanup configuration of resources
   The configuration (*btp_resources.tf*) is generated based on the information about the resources available from the provider plugin. All data including optinal data that got defaulted (e.g. usage in the btp_subaccount resource) is added to the configuration. To reduce the amount of data you could remove optional attributes that are optional and you do not want to have set explicitly. --> like what for example?

3. Declare variables
   The generated code doesn't contain any variables. We recommend to move the following into the *provider.tf* file

   - subdomain of the global account
   - %s ID: %s

    Depending on your requirements you might want to add further parameters to the variable list like the region your subaccount is created in.

4. Add dependencies
   As the export process doesn't detect dependencies, we recommend to add these manually. A typical scenario is the dependency between entitlements and the services/subscriptions specified in your configuration. Any more details on this?

5. Define a place for the state
   The state of your configuration should be stored in a remote state backend. Make sure to add the corresponding configuration (e.g. the *provider.tf*). You find more details in the [Terraform documentation](https://developer.hashicorp.com/terraform/language/backend)

6. Validate the import
   Validate that the import is possible by executing '*terraform plan*'. Depending on the number of resources the planing should return a message like this:
   Plan: n to import, 0 to add, 0 to change, 0 to destroy.

Now you're all set to run '*terraform apply*', which will import the state and thus bring your SAP BTP resources under the management of Terraform. Congrats!

`, input.Level, input.Level, input.UUID, input.ConfigDir, input.Level, input.UUID)
}
