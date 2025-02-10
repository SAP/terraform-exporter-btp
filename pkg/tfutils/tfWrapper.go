package tfutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/viper"
)

func getIaCTool() (tool string, err error) {

	//For TESTING purposes, we can set the tool to be used
	tool = os.Getenv("BTPTF_IAC_TOOL")
	if tool != "" {
		return tool, nil
	}

	_, localerr := exec.LookPath("terraform")
	if localerr == nil {
		tool = "terraform"
		return tool, nil
	}

	_, localerr = exec.LookPath("tofu")
	if localerr == nil {
		tool = "tofu"
		return tool, nil
	}

	fmt.Print("\r\n")
	log.Fatalf("error finding Terraform or OpenTofu executable: %v", err)
	return "", err
}

func runTfCmdGeneric(args ...string) error {
	tool, err := getIaCTool()
	if err != nil {
		return err
	}

	verbose := viper.GetViper().GetBool("verbose")
	cmd := exec.Command(tool, args...)
	if verbose {
		cmd.Stdout = os.Stdout
	} else {
		cmd.Stdout = nil
	}

	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runTfShowJson(directory string) (*State, error) {
	chDir := fmt.Sprintf("-chdir=%s", directory)

	tool, err := getIaCTool()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(tool, chDir, "show", "-json")

	var outBuffer bytes.Buffer
	cmd.Stdout = &outBuffer

	err = cmd.Run()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return nil, err
	}

	var state State

	err = json.Unmarshal(outBuffer.Bytes(), &state)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil, err
	}

	return &state, nil
}
