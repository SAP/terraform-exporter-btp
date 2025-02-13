package tfutils

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestSetupConfigDir(t *testing.T) {

	t.Setenv("BTP_USERNAME", "username")
	t.Setenv("BTP_PASSWORD", "password")
	t.Setenv("BTP_GLOBALACCOUNT", "ga")

	curWd, _ := os.Getwd()

	SetupConfigDir("configFolder", true, SubaccountLevel)
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

func TestConfigureProviderBtp(t *testing.T) {
	tempDir, cleanup := setupTestEnvironment()
	defer cleanup()

	TmpFolder = tempDir

	t.Setenv("BTP_USERNAME", "testuser")
	t.Setenv("BTP_PASSWORD", "testpass")
	t.Setenv("BTP_GLOBALACCOUNT", "testaccount")
	t.Setenv("BTP_CLI_SERVER_URL", "https://test.com")

	ConfigureProvider(SubaccountLevel)
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

func TestConfigureProviderCF(t *testing.T) {
	tempDir, cleanup := setupTestEnvironment()
	defer cleanup()

	TmpFolder = tempDir

	t.Setenv("CF_USER", "testuser")
	t.Setenv("CF_PASSWORD", "testpass")
	t.Setenv("CF_API_URL", "https://test.com")

	ConfigureProvider(OrganizationLevel)
	expectedFilePath := filepath.Join(TmpFolder, "provider.tf")
	if _, err := os.Stat(expectedFilePath); os.IsNotExist(err) {
		t.Errorf("Expected file %s does not exist", expectedFilePath)
	}

	expectedContent := `terraform {
required_providers {
cloudfoundry = {
source  = "cloudfoundry/cloudfoundry"
version = "1.3.0"
}
}
}

provider "cloudfoundry" {
api_url = "https://test.com"
}`

	expectedContent = strings.Replace(expectedContent, "[VERSION]", CfProviderVersion[1:], -1)
	content, err := os.ReadFile(expectedFilePath)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", expectedFilePath, err)
	}

	if string(content) != expectedContent {
		t.Errorf("Content of the file %s does not match expected content.\nGot:\n%s\nExpected:\n%s", expectedFilePath, string(content), expectedContent)
	}
}
func TestValidateCfApiUrl(t *testing.T) {
	// Test case where CF_API_URL is set
	t.Run("CF_API_URL is set", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("validateCfApiUrl() panicked when CF_API_URL was set")
			}
		}()
		validateCfApiUrl("https://api.example.com")
	})

}

func TestValidateCfAuthenticationData(t *testing.T) {
	// Test case where all authentication data is provided
	t.Run("All authentication data provided", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("validateCfAuthenticationData() panicked when all authentication data was provided")
			}
		}()
		validateCfAuthenticationData("username", "password", "accessToken", "refreshToken", "clientId", "clientSecret")
	})

	// Test case where only username and password are provided
	t.Run("Only username and password provided", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("validateCfAuthenticationData() panicked when only username and password were provided")
			}
		}()
		validateCfAuthenticationData("username", "password", "", "", "", "")
	})

	// Test case where only client ID and client secret are provided
	t.Run("Only client ID and client secret provided", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("validateCfAuthenticationData() panicked when only client ID and client secret were provided")
			}
		}()
		validateCfAuthenticationData("", "", "", "", "clientId", "clientSecret")
	})
}
func TestValidateGlobalAccount(t *testing.T) {
	// Test case where global account is provided
	t.Run("Global account provided", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("validateGlobalAccount() panicked when global account was provided")
			}
		}()
		validateGlobalAccount("testGlobalAccount")
	})

}
func TestValidateBtpAuthenticationData(t *testing.T) {
	// Test case where username and password are provided
	t.Run("Username and password provided", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("validateBtpAuthenticationData() panicked when username and password were provided")
			}
		}()
		validateBtpAuthenticationData("username", "password", "", "", "")
	})

	// Test case where TLS client certificate, key, and IDP URL are provided
	t.Run("TLS client certificate, key, and IDP URL provided", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("validateBtpAuthenticationData() panicked when TLS client certificate, key, and IDP URL were provided")
			}
		}()
		validateBtpAuthenticationData("", "", "tlsClientCertificate", "tlsClientKey", "tlsIdpURL")
	})

}
func TestAllStringsEmpty(t *testing.T) {
	// Test case where all strings are empty
	t.Run("All strings empty", func(t *testing.T) {
		result := allStringsEmpty("", " ", "   ")
		if !result {
			t.Errorf("Expected true, got false")
		}
	})

	// Test case where one string is not empty
	t.Run("One string not empty", func(t *testing.T) {
		result := allStringsEmpty("", "not empty", "   ")
		if result {
			t.Errorf("Expected false, got true")
		}
	})

	// Test case where all strings are not empty
	t.Run("All strings not empty", func(t *testing.T) {
		result := allStringsEmpty("a", "b", "c")
		if result {
			t.Errorf("Expected false, got true")
		}
	})

	// Test case where no strings are provided
	t.Run("No strings provided", func(t *testing.T) {
		result := allStringsEmpty()
		if !result {
			t.Errorf("Expected true, got false")
		}
	})
}

func TestHandleReturnWoInput(t *testing.T) {
	t.Run("Error is unexpected newline", func(t *testing.T) {
		err := fmt.Errorf("unexpected newline")
		choice := handleReturnWoInput(err)
		if choice != "N" {
			t.Errorf("Expected choice to be 'N', got '%s'", choice)
		}
	})

}

func TestGetResourcesList(t *testing.T) {
	t.Run("Resources string is 'all'", func(t *testing.T) {
		level := SubaccountLevel
		expectedResources := GetValidResourcesByLevel(level)
		resources := GetResourcesList("all", level)
		if !slices.Equal(resources, expectedResources) {
			t.Errorf("Expected resources %v, got %v", expectedResources, resources)
		}
	})

	t.Run("Resources string contains valid resources", func(t *testing.T) {
		level := SubaccountLevel
		expectedResources := []string{CmdSubaccountParameter, CmdEntitlementParameter}
		resourcesString := strings.Join(expectedResources, ",")
		resources := GetResourcesList(resourcesString, level)
		if !slices.Equal(resources, expectedResources) {
			t.Errorf("Expected resources %v, got %v", expectedResources, resources)
		}
	})
}

func TestGetValidResourcesByLevel(t *testing.T) {
	t.Run("Subaccount level", func(t *testing.T) {
		expectedResources := AllowedResourcesSubaccount
		resources := GetValidResourcesByLevel(SubaccountLevel)
		if !slices.Equal(resources, expectedResources) {
			t.Errorf("Expected resources %v, got %v", expectedResources, resources)
		}
	})

	t.Run("Directory level", func(t *testing.T) {
		expectedResources := AllowedResourcesDirectory
		resources := GetValidResourcesByLevel(DirectoryLevel)
		if !slices.Equal(resources, expectedResources) {
			t.Errorf("Expected resources %v, got %v", expectedResources, resources)
		}
	})

	t.Run("Organization level", func(t *testing.T) {
		expectedResources := AllowedResourcesOrganization
		resources := GetValidResourcesByLevel(OrganizationLevel)
		if !slices.Equal(resources, expectedResources) {
			t.Errorf("Expected resources %v, got %v", expectedResources, resources)
		}
	})

	t.Run("Invalid level", func(t *testing.T) {
		expectedResources := []string{}
		resources := GetValidResourcesByLevel("invalidLevel")
		if !slices.Equal(resources, expectedResources) {
			t.Errorf("Expected resources %v, got %v", expectedResources, resources)
		}
	})
}
