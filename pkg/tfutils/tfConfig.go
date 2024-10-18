package tfutils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	files "github.com/SAP/terraform-exporter-btp/pkg/files"
	output "github.com/SAP/terraform-exporter-btp/pkg/output"
	"github.com/spf13/viper"
	"github.com/theckman/yacspin"
)

var TmpFolder string

var AllowedResourcesSubaccount = []string{
	CmdSubaccountParameter,
	CmdEntitlementParameter,
	CmdEnvironmentInstanceParameter,
	CmdSubscriptionParameter,
	CmdTrustConfigurationParameter,
	CmdRoleParameter,
	CmdRoleCollectionParameter,
	CmdServiceBindingParameter,
	CmdServiceInstanceParameter,
	CmdSecuritySettingParameter,
}

var AllowedResourcesDirectory = []string{
	CmdDirectoryParameter,
	CmdEntitlementParameter,
	CmdRoleParameter,
	CmdRoleCollectionParameter,
}

func GenerateConfig(resourceFileName string, configFolder string, isMainCmd bool, resourceNameLong string) error {

	var spinner *yacspin.Spinner
	var err error

	if isMainCmd {
		// We must distinguish if the command is run from a main command or via delegation from helper functions
		spinner = output.StartSpinner("generating Terraform configuration for " + resourceNameLong)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %v", err)
	}

	terraformConfigPath := filepath.Join(currentDir, configFolder)
	err = os.Chdir(terraformConfigPath)
	if err != nil {
		return fmt.Errorf("error changing directory to %s: %v", terraformConfigPath, err)
	}

	if err := runTerraformCommand("init"); err != nil {
		return fmt.Errorf("error running Terraform init: %v", err)
	}

	planOption := "--generate-config-out=" + resourceFileName
	if err := runTerraformCommand("plan", planOption); err != nil {
		return fmt.Errorf("error running Terraform plan: %v", err)
	}

	if err := runTerraformCommand("fmt", "-recursive", "-list=false"); err != nil {
		return fmt.Errorf("error running Terraform fmt: %v", err)
	}

	//Switch back to the original directory
	err = os.Chdir(currentDir)
	if err != nil {
		return fmt.Errorf("error changing directory to %s: %v", currentDir, err)
	}

	cleanup()

	if isMainCmd {
		output.StopSpinner(spinner)
		fmt.Println(output.ColorStringGrey("   temporary files deleted"))
	}

	return nil
}

func ConfigureProvider() {
	tmpdir, err := os.MkdirTemp("", "provider.tf")
	if err != nil {
		panic(err)
	}

	TmpFolder = tmpdir
	abspath := filepath.Join(tmpdir, "provider.tf")

	verbose := viper.GetViper().GetBool("verbose")

	if verbose {
		fmt.Printf("temp file created at %s\n", abspath)
	}

	username := os.Getenv("BTP_USERNAME")
	password := os.Getenv("BTP_PASSWORD")
	enableSSO := os.Getenv("BTP_ENABLE_SSO")
	cliServerUrl := os.Getenv("BTP_CLI_SERVER_URL")
	globalAccount := os.Getenv("BTP_GLOBALACCOUNT")
	idp := os.Getenv("BTP_IDP")
	tlsClientCertificate := os.Getenv("BTP_TLS_CLIENT_CERTIFICATE")
	tlsClientKey := os.Getenv("BTP_TLS_CLIENT_KEY")
	tlsIdpURL := os.Getenv("BTP_TLS_IDP_URL")

	providerContent := "terraform {\nrequired_providers {\nbtp = {\nsource  = \"SAP/btp\"\nversion = \"" + BtpProviderVersion[1:] + "\"\n}\n}\n}\n\nprovider \"btp\" {\n"

	if !(len(strings.TrimSpace(username)) != 0 && len(strings.TrimSpace(password)) != 0) {
		if len(strings.TrimSpace(enableSSO)) == 0 {
			cleanup()
			fmt.Print("\r\n")
			log.Fatalf("set BTP_USERNAME and BTP_PASSWORD environment variable or enable SSO for login.")
		}
	}

	if len(strings.TrimSpace(globalAccount)) == 0 {
		cleanup()
		fmt.Print("\r\n")
		log.Fatalf("global account not set. set BTP_GLOBALACCOUNT environment variable to set global account")
	} else {
		providerContent = providerContent + "globalaccount = \"" + globalAccount + "\"\n"
	}

	if len(strings.TrimSpace(cliServerUrl)) != 0 {
		providerContent = providerContent + "cli_server_url=\"" + cliServerUrl + "\"\n"
	}

	if len(strings.TrimSpace(idp)) != 0 {
		providerContent = providerContent + "idp=\"" + idp + "\"\n"
	}

	if len(strings.TrimSpace(tlsClientCertificate)) != 0 {
		providerContent = providerContent + "tls_client_certificate =\"" + tlsClientCertificate + "\"\n"
	}

	if len(strings.TrimSpace(tlsClientKey)) != 0 {
		providerContent = providerContent + "tls_client_key =\"" + tlsClientKey + "\"\n"
	}

	if len(strings.TrimSpace(tlsIdpURL)) != 0 {
		providerContent = providerContent + "tls_idp_url =\"" + tlsIdpURL + "\"\n"
	}

	providerContent = providerContent + "}"

	err = files.CreateFileWithContent(abspath, providerContent)
	if err != nil {
		cleanup()
		fmt.Print("\r\n")
		log.Fatalf("create file %s failed!", abspath)
	}

}

func SetupConfigDir(configFolder string, isMainCmd bool) {

	if isMainCmd {
		message := "set up config directory \"" + configFolder + "\""
		fmt.Println(output.ColorStringGrey(message))
	}

	if len(TmpFolder) == 0 {
		ConfigureProvider()
	}
	curWd, err := os.Getwd()
	if err != nil {
		CleanupProviderConfig()
		fmt.Print("\r\n")
		log.Fatalf("error getting current working directory: %v", err)
	}

	configFilepath := filepath.Join(curWd, configFolder)

	exist, err := files.Exists(configFilepath)
	if err != nil {
		CleanupProviderConfig()
		fmt.Print("\r\n")
		log.Fatalf("error reading configuration folder %s: %v", configFolder, err)
	}

	if !exist {
		err = os.Mkdir(configFilepath, 0700)
		if err != nil {
			CleanupProviderConfig()
			fmt.Print("\r\n")
			log.Fatalf("error creating configuration folder %s at %s: %v", configFolder, curWd, err)
		}
	} else {
		fmt.Print("the configuration directory already exist. Do you want to continue? If yes then the directory will be overwritten (y/N): ")
		var choice string

		_, err = fmt.Scanln(&choice)
		if err != nil {
			if err.Error() == "unexpected newline" {
				choice = "N"
			} else {
				CleanupProviderConfig()
				fmt.Print("\r\n")
				log.Fatalf("error reading input: %v", err)
			}
		}

		choice = strings.TrimSpace(choice)
		if choice == "" {
			choice = "N"
		}

		// We acccept "Yes" or "No" as entry and take the first letter only
		if strings.ToUpper(choice[:1]) == "N" {
			CleanupProviderConfig()
			os.Exit(0)
		} else if strings.ToUpper(choice[:1]) == "Y" {
			fmt.Println(output.ColorStringCyan("existing files will be overwritten"))

			// Configuration folder must be re-created, otherwiese the Terraform commands will fail
			err := recreateExistingConfigDir(configFilepath)
			if err != nil {
				CleanupProviderConfig()
				fmt.Print("\r\n")
				log.Fatalf("error recreating configuration folder %s at %s: %v", configFolder, curWd, err)
			}
		} else {
			CleanupProviderConfig()
			fmt.Print("\r\n")
			log.Fatalf("invalid input. exiting the process")
		}
	}

	sourceFile, err := os.Open(TmpFolder + "/provider.tf")
	if err != nil {
		CleanupProviderConfig()
		fmt.Print("\r\n")
		log.Fatalf("failed to open file 'provider.tf' at %s: %v", TmpFolder, err)
	}
	defer sourceFile.Close()

	fullpath := filepath.Join(curWd, configFolder)

	destinationFile, err := os.Create(fullpath + "/provider.tf")
	if err != nil {
		CleanupProviderConfig()
		fmt.Print("\r\n")
		log.Fatalf("failed to create file 'provider.tf' at %s: %v", fullpath, err)
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		CleanupProviderConfig(fullpath)
		fmt.Print("\r\n")
		log.Fatalf("failed to copy file from temporary (%s) to final configuration directory (%s): %v", TmpFolder, fullpath, err)
	}
}

func GetResourcesList(resourcesString string, level string) []string {

	var resources []string

	allowedResources := GetValidResourcesByLevel(level)

	if resourcesString == "all" {
		resources = allowedResources
	} else {
		resources = strings.Split(resourcesString, ",")

		for _, resource := range resources {
			if !(slices.Contains(allowedResources, resource)) {

				allowedResourceList := strings.Join(allowedResources, ", ")
				fmt.Print("\r\n")
				log.Fatal("please check the resource provided. Currently supported resources are " + allowedResourceList + ". Provide 'all' to check for all resources")
			}
		}
	}

	return resources
}

func CleanupProviderConfig(directory ...string) {
	cleanup()

	for _, dir := range directory {
		CleanupTempFiles(dir)
	}
}

func GetValidResourcesByLevel(level string) []string {
	if level == SubaccountLevel {
		return AllowedResourcesSubaccount
	}

	return AllowedResourcesDirectory
}

func cleanup() {
	// Cleanup temporary folder variable
	TmpFolder = ""

	// Cleanup temporary files on disk
	err := os.RemoveAll(TmpFolder)
	if err != nil {
		fmt.Print("\r\n")
		log.Fatalf("error deleting temp files: %v", err)
	}
}

func FinalizeTfConfig(configFolder string) {

	output.AddNewLine()
	spinner := output.StartSpinner("finalizing Terraform configuration")

	currentDir, err := os.Getwd()
	if err != nil {
		CleanupProviderConfig()
		fmt.Print("\r\n")
		log.Fatalf("error getting current directory: %v", err)
	}

	terraformConfigPath := filepath.Join(currentDir, configFolder)

	err = os.Chdir(terraformConfigPath)
	if err != nil {
		CleanupProviderConfig()
		fmt.Print("\r\n")
		log.Fatalf("error changing directory to %s: %v \n", terraformConfigPath, err)
	}

	if err := runTerraformCommand("init"); err != nil {
		CleanupProviderConfig()
		fmt.Print("\r\n")
		log.Fatalf("error initializing Terraform: %v", err)
	}

	if err := runTerraformCommand("fmt", "-recursive", "-list=false"); err != nil {
		CleanupProviderConfig()
		fmt.Print("\r\n")
		log.Fatalf("error running Terraform fmt: %v", err)
	}

	//Switch back to the original directory
	err = os.Chdir(currentDir)
	if err != nil {
		CleanupProviderConfig()
		fmt.Print("\r\n")
		log.Fatalf("error changing directory to %s: %v \n", currentDir, err)
	}

	output.StopSpinner(spinner)
}

// Convenience functions that wrap repetitive steps
func ExecPreExportSteps(tempConfigDir string) {
	SetupConfigDir(tempConfigDir, false)
}

func ExecPostExportSteps(tempConfigDir string, targetConfigDir string, targetResourceFileName string, resourceNameLong string) {

	spinner := output.StartSpinner("generating Terraform configuration for " + resourceNameLong)

	err := GenerateConfig(targetResourceFileName, tempConfigDir, false, resourceNameLong)
	if err != nil {
		CleanupTempFiles(tempConfigDir)
		fmt.Print("\r\n")
		log.Fatalf("error generating Terraform configuration for %s: %v", resourceNameLong, err)
	}

	err = mergeTfConfig(targetConfigDir, targetResourceFileName, tempConfigDir, resourceNameLong)
	if err != nil {
		CleanupTempFiles(tempConfigDir)
		fmt.Print("\r\n")
		log.Fatalf("error merging Terraform configuration for %s: %v", resourceNameLong, err)
	}

	CleanupTempFiles(tempConfigDir)

	output.StopSpinner(spinner)

	fmt.Println(output.ColorStringGrey("   temporary files deleted"))
}

func CleanupTempFiles(tempConfigDir string) {
	files.DeleteSourceFolder(tempConfigDir)
}

func mergeTfConfig(configFolder string, fileName string, resourceConfigFolder string, resourceName string) error {

	currentDir, err := os.Getwd()

	if err != nil {
		return fmt.Errorf("error getting current directory: %v", err)
	}

	sourceConfigPath := filepath.Join(currentDir, resourceConfigFolder, fileName)

	// Check if the source file exists
	exist, err := files.Exists(sourceConfigPath)
	if err != nil {
		return fmt.Errorf("error checking if source directory exists: %v", err)
	}

	if !exist {
		// Nothing to do as the source file does not exist
		return nil
	}

	sourceFile, err := os.Open(sourceConfigPath)
	if err != nil {
		return fmt.Errorf("error opening resource config file: %v", err)
	}
	defer sourceFile.Close()

	targetConfigPath := filepath.Join(currentDir, configFolder, fileName)

	exist, err = files.Exists(targetConfigPath)
	if err != nil {
		return fmt.Errorf("error checking if target directory exists: %v", err)
	}

	if !exist {
		// In the first run we must create the file if it does not exist
		_, err := os.Create(targetConfigPath)
		if err != nil {
			return fmt.Errorf("error creating target configuration file: %v", err)
		}
	}

	targetFile, err := os.OpenFile(targetConfigPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening target configuration file: %v", err)
	}
	defer targetFile.Close()

	headerTemplate := `
###
# Resource: ` + resourceName + `
###
`
	if _, err := targetFile.Write([]byte(headerTemplate)); err != nil {
		return fmt.Errorf("error adding header line to target file: %v", err)
	}

	if _, err := io.Copy(targetFile, sourceFile); err != nil {
		return fmt.Errorf("error copying resource file to target file: %v", err)
	}

	err = files.CopyImportFiles(resourceConfigFolder, configFolder)
	if err != nil {
		return fmt.Errorf("error copying import files: %v", err)
	}
	return nil
}

func recreateExistingConfigDir(filepath string) error {
	err := os.RemoveAll(filepath)
	if err != nil {
		return fmt.Errorf("error recreating existing configuration folder %s: %v", filepath, err)
	}

	err = os.Mkdir(filepath, 0700)
	if err != nil {
		return fmt.Errorf("error recreating configuration folder %s: %v", filepath, err)
	}

	return nil
}
