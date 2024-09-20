package output

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/viper"
	"github.com/theckman/yacspin"
)

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

func StartSpinner(message string) (*yacspin.Spinner, error) {

	// No spinner execution during debug mode
	debug := viper.GetViper().GetBool("debug")
	if debug {
		return nil, nil
	}

	spinner, err := createSpinner(message)
	if err != nil {
		fmt.Printf("failed to make spinner from config struct: %v\n", err)
		os.Exit(1)
	}

	stopOnSignal(spinner)

	err = renderSpinner(spinner)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return spinner, nil
}

func StopSpinner(spinner *yacspin.Spinner) error {

	// No spinner execution during debug mode
	debug := viper.GetViper().GetBool("debug")
	if debug {
		return nil
	}

	if err := spinner.Stop(); err != nil {
		return fmt.Errorf("failed to stop spinner: %w", err)
	}

	return nil
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

func AddNewLine() {
	fmt.Println("")
}
