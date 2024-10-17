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
	return fmt.Sprintf(`# Next Steps

## Status Quo

After executing the *btptf* CLI you have now created the Terraform configuration in the directory %s as well as the Terraform [import blocks](https://developer.hashicorp.com/terraform/language/import) needed to import the state of the configuration of the %s with ID %s.

With this in place you can further proceed to bring the resources under the management of Terraform. The following section will give you some advice on the next steps.

## Next Steps

### Review of Configuration

You should in any case review the configuration. While being technically valid, the configuration might not reflect the way you would want to have the configuration be done. Therefore the first step should always be a review by one of the experts to validate the generated code. This also comprises the provider configuration especially the version constraints defined there.

### CleanUp of Configuration

The configuration is generated based on the information available from the provider plugin. All data including default data is added to the configuration. To reduce the amount of data you could remove optional attributes that are optional and you do not want to have set explicitly.

### Addition of Variables

The generated code does not contain any variables. It makes sense to put at least the value of the subdomain of the global account in the *provider.tf* file as well as the ID %s as variables. Depending on your requirements you would also want to add further parameters to the variable list.

### Adding Dependencies

The export is not capable of detecting explicit dependencies. This must be manually added. One typical scenario is the dependency between entitlements and the services/subscriptions that are defined in your configuration.

### Adding State Backend

The state of your configuration should be stored in a remote state backend. Make sure to add the corresponding configuration. You find more details in the [Terraform documentation](https://developer.hashicorp.com/terraform/language/backend)

### State Import

If you have finished the refinement of the configuration you can trigger the import of the state via *terraform apply*
`, input.ConfigDir, input.Level, input.UUID, input.UUID)
}
