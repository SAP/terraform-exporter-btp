package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSetupConfigDir(t *testing.T) {

	t.Setenv("BTP_USERNAME", "username")
	t.Setenv("BTP_PASSWORD", "password")
	t.Setenv("BTP_GLOBALACCOUNT", "ga")

	curWd, _ := os.Getwd()

	setupConfigDir("configFolder", true)
	if _, err := os.Stat(filepath.Join(curWd, "configFolder")); os.IsNotExist(err) {
		t.Errorf("Directory should have been created")
	}

	if _, err := os.Stat(filepath.Join(curWd, "configFolder", "provider.tf")); os.IsNotExist(err) {
		t.Errorf("File should have been copied to existing directory")
	}

	os.RemoveAll(filepath.Join(curWd, "configFolder"))
	cleanup()
}

func setupTestEnvironment() (string, func()) {
	tempDir, err := os.MkdirTemp("", "providerTest")
	if err != nil {
		panic(err)
	}

	return tempDir, func() {
		os.RemoveAll(tempDir)
	}
}

func TestConfigureProvider(t *testing.T) {
	tempDir, cleanup := setupTestEnvironment()
	defer cleanup()

	TmpFolder = tempDir

	t.Setenv("BTP_USERNAME", "testuser")
	t.Setenv("BTP_PASSWORD", "testpass")
	t.Setenv("BTP_GLOBALACCOUNT", "testaccount")
	t.Setenv("BTP_CLI_SERVER_URL", "https://test.com")

	configureProvider()
	expectedFilePath := filepath.Join(TmpFolder, "provider.tf")
	if _, err := os.Stat(expectedFilePath); os.IsNotExist(err) {
		t.Errorf("Expected file %s does not exist", expectedFilePath)
	}

	expectedContent := `terraform {
required_providers {
btp = {
source  = "SAP/btp"
version = "[VERSION]"
}
}
}

provider "btp" {
globalaccount = "testaccount"
cli_server_url="https://test.com"
}`

	expectedContent = strings.Replace(expectedContent, "[VERSION]", BtpProviderVersion[1:], -1)
	content, err := os.ReadFile(expectedFilePath)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", expectedFilePath, err)
	}

	if string(content) != expectedContent {
		t.Errorf("Content of the file %s does not match expected content.\nGot:\n%s\nExpected:\n%s", expectedFilePath, string(content), expectedContent)
	}
}
