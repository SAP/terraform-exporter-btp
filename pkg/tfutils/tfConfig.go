package tfutils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	files "github.com/SAP/terraform-exporter-btp/pkg/files"
	output "github.com/SAP/terraform-exporter-btp/pkg/output"
	"github.com/SAP/terraform-exporter-btp/pkg/resume"
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

var AllowedResourcesOrganization = []string{
	CmdCfSpaceParameter,
	CmdCfUserParameter,
	CmdCfDomainParamater,
	CmdCfOrgRoleParameter,
	CmdCfRouteParameter,
	CmdCfSpaceQuotaParameter,
	CmdCfServiceInstanceParameter,
	CmdCfSpaceRoleParameter,
}

const selectionOverwrite = "Overwrite the existing directory and continue"
const selectionResume = "Resume the existing export"
const selectionAbort = "Abort the processing"
const selectionInvalid = "Invalid selection"

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

	if err := runTfCmdGeneric("init"); err != nil {
		return fmt.Errorf("error running Terraform init: %v", err)
	}

	planOption := "--generate-config-out=" + resourceFileName
	if err := runTfCmdGeneric("plan", planOption); err != nil {
		return fmt.Errorf("error running Terraform plan: %v", err)
	}

	if err := runTfCmdGeneric("fmt", "-recursive", "-list=false"); err != nil {
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

func ConfigureProvider(level string) {
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

	var providerContent string

	switch level {
	case SubaccountLevel, DirectoryLevel:

		username := os.Getenv("BTP_USERNAME")
		password := os.Getenv("BTP_PASSWORD")
		cliServerUrl := os.Getenv("BTP_CLI_SERVER_URL")
		globalAccount := os.Getenv("BTP_GLOBALACCOUNT")
		idp := os.Getenv("BTP_IDP")
		tlsClientCertificate := os.Getenv("BTP_TLS_CLIENT_CERTIFICATE")
		tlsClientKey := os.Getenv("BTP_TLS_CLIENT_KEY")
		tlsIdpURL := os.Getenv("BTP_TLS_IDP_URL")
		ssoEnabled := os.Getenv("BTP_ENABLE_SSO")

		validateBtpAuthenticationData(username, password, tlsClientCertificate, tlsClientKey, tlsIdpURL, ssoEnabled)
		validateGlobalAccount(globalAccount)

		providerContent = "terraform {\nrequired_providers {\nbtp = {\nsource  = \"SAP/btp\"\nversion = \"" + BtpProviderVersion[1:] + "\"\n}\n}\n}\n\nprovider \"btp\" {\n"
		providerContent = providerContent + "globalaccount = \"" + globalAccount + "\"\n"

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

	case OrganizationLevel, SpaceLevel:

		username := os.Getenv("CF_USER")
		password := os.Getenv("CF_PASSWORD")
		apiUrl := os.Getenv("CF_API_URL")
		cfOrigin := os.Getenv("CF_ORIGIN")
		cfClientId := os.Getenv("CF_CLIENT_ID")
		cfClientSecret := os.Getenv("CF_CLIENT_SECRET")
		cfAccessToken := os.Getenv("CF_ACCESS_TOKEN")
		cfRefreshToken := os.Getenv("CF_REFRESH_TOKEN")

		validateCfAuthenticationData(username, password, cfAccessToken, cfRefreshToken, cfClientId, cfClientSecret)
		validateCfApiUrl(apiUrl)

		providerContent = "terraform {\nrequired_providers {\ncloudfoundry = {\nsource  = \"cloudfoundry/cloudfoundry\"\nversion = \"" + CfProviderVersion[1:] + "\"\n}\n}\n}\n\nprovider \"cloudfoundry\" {\n"
		providerContent = providerContent + "api_url = \"" + apiUrl + "\"\n"

		if len(strings.TrimSpace(cfOrigin)) != 0 {
			providerContent = providerContent + "origin=\"" + cfOrigin + "\"\n"
		}

		if len(strings.TrimSpace(cfClientId)) != 0 {
			providerContent = providerContent + "cf_client_id =\"" + cfClientId + "\"\n"
		}

		if len(strings.TrimSpace(cfClientSecret)) != 0 {
			providerContent = providerContent + "cf_client_secret =\"" + cfClientSecret + "\"\n"
		}

		if len(strings.TrimSpace(cfAccessToken)) != 0 {
			providerContent = providerContent + "cf_access_token =\"" + cfAccessToken + "\"\n"
		}

		if len(strings.TrimSpace(cfRefreshToken)) != 0 {
			providerContent = providerContent + "cf_refresh_token =\"" + cfRefreshToken + "\"\n"
		}

		providerContent = providerContent + "}"

	}

	err = files.CreateFileWithContent(abspath, providerContent)
	if err != nil {
		cleanup()
		fmt.Print("\r\n")
		log.Fatalf("create file %s failed!", abspath)
	}

}

func validateCfApiUrl(apiUrl string) {
	if len(strings.TrimSpace(apiUrl)) == 0 {
		cleanup()
		fmt.Print("\r\n")
		log.Fatalf("cf api URL not set. set CF_API_URL environment variable to set CF API endpoint")
	}
}

func validateCfAuthenticationData(username string, password string, cfAccessToken string, cfRefreshToken string, cfClientId string, cfClientSecret string) {
	if allStringsEmpty(username, password, cfAccessToken, cfRefreshToken, cfClientId, cfClientSecret) {
		cleanup()
		fmt.Print("\r\n")
		log.Fatalf("set Cloud Foundry environment variables for login.")
	}
}

func validateGlobalAccount(globalAccount string) {
	if allStringsEmpty(globalAccount) {
		cleanup()
		fmt.Print("\r\n")
		log.Fatalf("global account not set. set BTP_GLOBALACCOUNT environment variable to set global account")
	}
}

func validateBtpAuthenticationData(username string, password string, tlsClientCertificate string, tlsClientKey string, tlsIdpURL string, ssoEnabled string) {

	if ssoEnabled != "" {
		fmt.Print("\r\n")
		log.Fatalf("SSO is not supported for this command. Please remove the BTP_ENABLE_SSO environment variable and re-run the command.")
	}

	// Check if any of the authentication data is set (username and password or TLS client certificate and key)
	if allStringsEmpty(username, password) && allStringsEmpty(tlsClientCertificate, tlsClientKey, tlsIdpURL) {
		cleanup()
		fmt.Print("\r\n")
		log.Fatalf("set valid authentication data for login e.g. BTP_USERNAME and BTP_PASSWORD environment variables.")
	}
}

func allStringsEmpty(stringsToCheck ...string) bool {

	for _, str := range stringsToCheck {
		if len(strings.TrimSpace(str)) != 0 {
			return false
		}
	}

	return true
}

func SetupConfigDir(configFolder string, isMainCmd bool, level string) {

	if isMainCmd {
		message := "set up config directory \"" + configFolder + "\""
		fmt.Println(output.ColorStringGrey(message))
	}

	if len(TmpFolder) == 0 {
		ConfigureProvider(level)
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
		createNewConfigDir(configFilepath, configFolder, curWd)
	} else {
		handleExistingDir(isMainCmd, configFilepath, configFolder, curWd)
	}

	sourceFile, err := os.Open(TmpFolder + "/provider.tf")
	if err != nil {
		CleanupProviderConfig()
		fmt.Print("\r\n")
		log.Fatalf("failed to open file 'provider.tf' at %s: %v", TmpFolder, err)
	}

	defer func() {
		if tempErrSrc := sourceFile.Close(); tempErrSrc != nil {
			CleanupProviderConfig()
			fmt.Print("\r\n")
			log.Fatalf("failed to close file 'provider.tf' at %s: %v", TmpFolder, tempErrSrc)
		}
	}()

	fullpath := filepath.Join(curWd, configFolder)

	destinationFile, err := os.Create(fullpath + "/provider.tf")
	if err != nil {
		CleanupProviderConfig()
		fmt.Print("\r\n")
		log.Fatalf("failed to create file 'provider.tf' at %s: %v", fullpath, err)
	}

	defer func() {
		if tempErrDest := destinationFile.Close(); tempErrDest != nil {
			CleanupProviderConfig(fullpath)
			fmt.Print("\r\n")
			log.Fatalf("failed to close file 'provider.tf' at %s: %v", fullpath, tempErrDest)
		}
	}()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		CleanupProviderConfig(fullpath)
		fmt.Print("\r\n")
		log.Fatalf("failed to copy file from temporary (%s) to final configuration directory (%s): %v", TmpFolder, fullpath, err)
	}
}

func handleInputExistingDir(choice string, configFilepath string, configFolder string, curWd string) {
	switch choice {
	case selectionAbort:
		CleanupProviderConfig()
		os.Exit(0)
	case selectionOverwrite:
		output.AddNewLine()
		fmt.Println(output.ColorStringCyan("existing files will be overwritten"))
		output.AddNewLine()

		err := recreateExistingConfigDir(configFilepath)
		if err != nil {
			CleanupProviderConfig()
			output.AddNewLine()
			log.Fatalf("error recreating configuration folder %s at %s: %v", configFolder, curWd, err)
		}
	case selectionResume:
		// Can only happen if we are in the main command
		// Do nothing, the processing will be resumed with the existing directory
		output.AddNewLine()
		fmt.Println(output.ColorStringCyan("export will be resumed"))
		output.AddNewLine()
	default:
		CleanupProviderConfig()
		output.AddNewLine()
		log.Fatalf("invalid input. exiting the process")
	}
}

func handleReturnWoInput(err error) (choice string) {
	if err.Error() == "unexpected newline" {
		choice = "N"
	} else {
		CleanupProviderConfig()
		fmt.Print("\r\n")
		log.Fatalf("error reading input: %v", err)
	}
	return choice
}

func createNewConfigDir(configFilepath string, configFolder string, curWd string) {
	err := os.Mkdir(configFilepath, 0700)
	if err != nil {
		CleanupProviderConfig()
		fmt.Print("\r\n")
		log.Fatalf("error creating configuration folder %s at %s: %v", configFolder, curWd, err)
	}
}

func handleExistingDir(isMainCmd bool, configFilepath string, configFolder string, curWd string) {
	var err error
	var importLog []string
	if isMainCmd {
		// check if an import log exists to decide how to proceed
		importLog, _ = resume.GetExistingExportLog(configFolder)
	}

	surveyName := "existingdir"
	surveyMessage := ""
	surveyDefault := selectionOverwrite
	answer := struct {
		ExistingDir string `survey:"existingdir"`
	}{}

	surveyOptions := []string{selectionOverwrite, selectionAbort}

	if len(importLog) > 0 {
		surveyMessage = fmt.Sprintf("the configuration directory '%s' with an import log exists. How do you want to continue?", configFolder)
		surveyOptions = append([]string{selectionResume}, surveyOptions...)
		surveyDefault = selectionResume
	} else {
		surveyMessage = fmt.Sprintf("the configuration directory '%s' already exists. How do you want to continue?", configFolder)
	}

	qs := []*survey.Question{
		{
			Name: surveyName,
			Prompt: &survey.Select{
				Message: surveyMessage,
				Options: surveyOptions,
				Default: surveyDefault,
			},
		},
	}

	err = survey.Ask(qs, &answer)
	if err != nil {
		answer.ExistingDir = selectionInvalid
	}
	handleInputExistingDir(answer.ExistingDir, configFilepath, configFolder, curWd)

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

func FilterResourcesListByLog(resources []string, exportLog []string) []string {

	var filteredResources []string

	for _, resource := range resources {
		if !slices.Contains(exportLog, resource) {
			filteredResources = append(filteredResources, resource)
		}
	}

	return filteredResources
}

func CleanupProviderConfig(directory ...string) {
	cleanup()

	for _, dir := range directory {
		CleanupTempFiles(dir)
	}
}

func GetValidResourcesByLevel(level string) []string {
	switch level {
	case SubaccountLevel:
		return AllowedResourcesSubaccount
	case DirectoryLevel:
		return AllowedResourcesDirectory
	case OrganizationLevel:
		return AllowedResourcesOrganization
	default:
		return []string{}
	}
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

	if err := runTfCmdGeneric("init", "-backend=false"); err != nil {
		CleanupProviderConfig()
		fmt.Print("\r\n")
		log.Fatalf("error initializing Terraform: %v", err)
	}

	if err := runTfCmdGeneric("fmt", "-recursive", "-list=false"); err != nil {
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
func ExecPreExportSteps(tempConfigDir string, level string) {
	SetupConfigDir(tempConfigDir, false, level)
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
	defer func() {
		if tempErr := sourceFile.Close(); tempErr != nil {
			err = tempErr
		}
	}()

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

	defer func() {
		if tempErr := targetFile.Close(); tempErr != nil {
			err = tempErr
		}
	}()

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
