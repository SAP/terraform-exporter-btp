package tfutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

const minimumTerraformQueryMinor = 14

var getIaCTool = func() (string, error) {
	if _, err := exec.LookPath("terraform"); err != nil {
		return "", fmt.Errorf("terraform executable not found: %w", err)
	}
	return "terraform", nil
}

func runTfCmdGeneric(args ...string) (string, error) {
	tool, err := getIaCTool()
	if err != nil {
		return "", err
	}
	args = append(args, "-no-color")
	cmd := exec.Command(tool, args...)
	var stderr bytes.Buffer
	if viper.GetViper().GetBool("verbose") {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stderr = &stderr
	}
	if err := cmd.Run(); err != nil {
		return stderr.String(), err
	}
	return "", nil
}

func ValidateTerraformQuerySupport() error {
	if _, err := getIaCTool(); err != nil {
		return err
	}
	version, err := getTerraformVersion()
	if err != nil {
		return err
	}
	parts := strings.Split(version, ".")
	if len(parts) < 2 {
		return fmt.Errorf("invalid Terraform version %q", version)
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("invalid Terraform version %q: %w", version, err)
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid Terraform version %q: %w", version, err)
	}
	if major < 1 || (major == 1 && minor < minimumTerraformQueryMinor) {
		return fmt.Errorf("terraform 1.%d or newer is required; found %s", minimumTerraformQueryMinor, version)
	}
	return nil
}

func getTerraformVersion() (string, error) {
	output, err := exec.Command("terraform", "version", "-json").Output()
	if err != nil {
		return "", fmt.Errorf("failed to read Terraform version: %w", err)
	}
	var response struct {
		TerraformVersion string `json:"terraform_version"`
	}
	if err := json.Unmarshal(output, &response); err != nil {
		return "", fmt.Errorf("failed to parse Terraform version: %w", err)
	}
	return response.TerraformVersion, nil
}
