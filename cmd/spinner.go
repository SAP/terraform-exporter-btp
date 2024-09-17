package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

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

func startSpinner(message string) (*yacspin.Spinner, error) {

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

func stopSpinner(spinner *yacspin.Spinner) error {
	if err := spinner.Stop(); err != nil {
		return fmt.Errorf("failed to stop spinner: %w", err)
	}

	return nil
}

func printExportStartMessage() {
	fmt.Println("")
	fmt.Println("🚀 Terraform configuration export started ...")
	fmt.Println("")
}

func printExportSuccessMessage() {
	fmt.Println("")
	fmt.Println("🎉 Terraform configuration successfully created")
	fmt.Println("")
}
