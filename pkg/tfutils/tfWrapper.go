package tfutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/SAP/terraform-exporter-btp/pkg/toggles"
	"github.com/spf13/viper"
)

func getIaCTool() (tool string, err error) {

	//For TESTING purposes, we can set the tool to be used
	tool = toggles.GetIacTool()
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

// function return true if resource identity is supported in the installed terraform version, false otherwise
func isResourceIdentitySupported() (bool, error) {
	tool, err := getIaCTool()
	if err != nil {
		return false, err
	}

	if tool == "terraform" {
		version, err := getTerraformVersion()
		if err != nil {
			return false, fmt.Errorf("failed to get Terraform version: %w", err)
		}

		terraformVersion := strings.Split(version, ".")
		majorVersion := terraformVersion[0]
		minorVersion := terraformVersion[1]
		version = majorVersion + "." + minorVersion

		floatVersion, err := strconv.ParseFloat(version, 64)
		if err != nil {
			return false, fmt.Errorf("failed to parse Terraform version %s: %w", version, err)
		}

		if floatVersion >= 1.12 {
			return true, nil
		}

		return false, nil
	} else {
		return false, nil
	}
}

// This function returns terraform version.
func getTerraformVersion() (string, error) {

	cmd := exec.Command("terraform", "version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute %s version command: %w", "terraform", err)
	}

	versionOutput := strings.Split(string(output), "\n")[0]
	versionParts := strings.Fields(versionOutput)
	version := strings.TrimPrefix(versionParts[1], "v")

	return version, nil
}
