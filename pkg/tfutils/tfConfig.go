package tfutils

import (
	"errors"
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
)

var providerDirectory *files.TempDirectory

var AllowedResourcesSubaccount = []string{
	CmdSubaccountParameter,
	CmdEntitlementParameter,
	CmdEnvironmentInstanceParameter,
	CmdSubscriptionParameter,
	CmdTrustConfigurationParameter,
	CmdRoleParameter,
	CmdRoleCollectionParameter,
	CmdRoleCollectionAssignmentParameter,
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

func ConfigureProvider(level string) {
	tmpdir, err := files.NewTempDirectory("btptf-provider-")
	if err != nil {
		panic(err)
	}

	cleanup()
	providerDirectory = tmpdir
	abspath := filepath.Join(tmpdir.Root.Name(), "provider.tf")

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
		useBtpCLiSession := os.Getenv("USE_BTPCLI_SESSION")

		validateBtpAuthenticationData(username, password, tlsClientCertificate, tlsClientKey, tlsIdpURL, ssoEnabled, useBtpCLiSession)
		validateGlobalAccount(globalAccount)

		providerContent = "terraform {\nrequired_providers {\nbtp = {\nsource  = \"SAP/btp\"\nversion = \"" + BtpProviderVersion[1:] + "\"\n}\n}\n}\n\nprovider \"btp\" {\n"
		providerContent = providerContent + "globalaccount = \"" + EscapeHclString(globalAccount) + "\"\n"

		if len(strings.TrimSpace(cliServerUrl)) != 0 {
			providerContent = providerContent + "cli_server_url=\"" + EscapeHclString(cliServerUrl) + "\"\n"
		}

		if len(strings.TrimSpace(idp)) != 0 {
			providerContent = providerContent + "idp=\"" + EscapeHclString(idp) + "\"\n"
		}

		if len(strings.TrimSpace(tlsClientCertificate)) != 0 {
			providerContent = providerContent + "tls_client_certificate =\"" + EscapeHclString(tlsClientCertificate) + "\"\n"
		}

		if len(strings.TrimSpace(tlsClientKey)) != 0 {
			providerContent = providerContent + "tls_client_key =\"" + EscapeHclString(tlsClientKey) + "\"\n"
		}

		if len(strings.TrimSpace(tlsIdpURL)) != 0 {
			providerContent = providerContent + "tls_idp_url =\"" + EscapeHclString(tlsIdpURL) + "\"\n"
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

		// Sensitive credentials (cf_client_secret, cf_access_token, cf_refresh_token)
		// are intentionally not written to the generated provider configuration.
		// The Cloud Foundry provider sources them from the CF_CLIENT_SECRET,
		// CF_ACCESS_TOKEN and CF_REFRESH_TOKEN environment variables, which are
		// already required to be set for the export to run.
		providerContent = "terraform {\nrequired_providers {\ncloudfoundry = {\nsource  = \"cloudfoundry/cloudfoundry\"\nversion = \"" + CfProviderVersion[1:] + "\"\n}\n}\n}\n\nprovider \"cloudfoundry\" {\n"
		providerContent = providerContent + "api_url = \"" + EscapeHclString(apiUrl) + "\"\n"

		if len(strings.TrimSpace(cfOrigin)) != 0 {
			providerContent = providerContent + "origin=\"" + EscapeHclString(cfOrigin) + "\"\n"
		}

		if len(strings.TrimSpace(cfClientId)) != 0 {
			providerContent = providerContent + "cf_client_id =\"" + EscapeHclString(cfClientId) + "\"\n"
		}

		providerContent = providerContent + "}"

	}

	err = tmpdir.Root.WriteFile("provider.tf", []byte(providerContent), 0600)
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

func validateBtpAuthenticationData(username string, password string, tlsClientCertificate string, tlsClientKey string, tlsIdpURL string, ssoEnabled string, useBtpCLiSession string) {

	if ssoEnabled != "" {
		fmt.Print("\r\n")
		log.Fatalf("SSO is not supported for this command. Please remove the BTP_ENABLE_SSO environment variable and re-run the command.")
	}

	// Check if any of the authentication data is set (username and password or TLS client certificate and key)
	if allStringsEmpty(username, password) && allStringsEmpty(tlsClientCertificate, tlsClientKey, tlsIdpURL) && useBtpCLiSession != "true" {
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

// SetupConfigDir validates the CLI path before touching the destination and
// returns the handle used by all subsequent export operations.
func SetupConfigDir(configFolder string, isMainCmd bool, level string) (*os.Root, error) {
	return setupConfigDir(configFolder, isMainCmd, level, handleExistingDir)
}

func setupConfigDir(configFolder string, isMainCmd bool, level string, choose func(bool, *os.Root, string) (string, error)) (*os.Root, error) {
	if err := files.ValidateConfigDir(configFolder); err != nil {
		return nil, err
	}
	configFolder = filepath.Clean(configFolder)
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	work, err := os.OpenRoot(cwd)
	if err != nil {
		return nil, err
	}
	defer func() { _ = work.Close() }()
	if isMainCmd {
		fmt.Println(output.ColorStringGrey(fmt.Sprintf("set up config directory %q", configFolder)))
	}
	if _, err := work.Stat(configFolder); os.IsNotExist(err) {
		// Check the existing parent through the working root before configuring
		// credentials or creating anything at a user-supplied location.
		parent, err := work.OpenRoot(filepath.Dir(configFolder))
		if err != nil {
			return nil, err
		}
		_ = parent.Close()
		if providerDirectory == nil {
			ConfigureProvider(level)
		}
		if err := work.Mkdir(configFolder, 0700); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, fmt.Errorf("invalid configuration directory %q: %w", configFolder, err)
	} else {
		existing, err := openConfigDir(work, configFolder)
		if err != nil {
			return nil, err
		}
		if providerDirectory == nil {
			ConfigureProvider(level)
		}
		choice, err := choose(isMainCmd, existing, configFolder)
		_ = existing.Close()
		if err != nil {
			return nil, err
		}
		switch choice {
		case selectionAbort:
			CleanupProviderConfig()
			os.Exit(0)
		case selectionOverwrite:
			output.AddNewLine()
			fmt.Println(output.ColorStringCyan("existing files will be overwritten"))
			output.AddNewLine()
			if err := recreateExistingConfigDir(work, configFolder); err != nil {
				return nil, err
			}
		case selectionResume:
			fmt.Println(output.ColorStringCyan("export will be resumed"))
		default:
			return nil, fmt.Errorf("invalid input. exiting the process")
		}
	}
	root, err := openConfigDir(work, configFolder)
	if err != nil {
		return nil, err
	}
	if err := files.ValidateTerraformDirectory(root); err != nil {
		_ = root.Close()
		return nil, err
	}
	if err := files.CopyFile(providerDirectory.Root, root, "provider.tf", os.O_TRUNC); err != nil {
		_ = root.Close()
		return nil, fmt.Errorf("copy provider configuration: %w", err)
	}
	return root, nil
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

func handleExistingDir(isMainCmd bool, directory *os.Root, configFolder string) (string, error) {
	var err error
	var importLog []string
	if isMainCmd {
		// check if an import log exists to decide how to proceed
		importLog, _ = resume.GetExistingExportLog(directory)
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
	return answer.ExistingDir, err

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

func CleanupProviderConfig(directory ...*files.TempDirectory) {
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
	dir := providerDirectory
	providerDirectory = nil
	if err := dir.Cleanup(); err != nil {
		log.Printf("error deleting provider temporary files: %v", err)
	}
}

func FinalizeTfConfig(configFolder *os.Root, target ExecutionTarget, importsOnly bool) {

	output.AddNewLine()
	spinner := output.StartSpinner("finalizing Terraform configuration")
	if err := files.ValidateTerraformDirectory(configFolder); err != nil {
		CleanupProviderConfig()
		log.Fatal(err)
	}

	if err := applyExecutionTargetVariables(configFolder, target); err != nil {
		CleanupProviderConfig()
		fmt.Print("\r\n")
		log.Fatalf("error applying target variables: %v", err)
	}
	if err := normalizeGeneratedResourceAddresses(configFolder); err != nil {
		CleanupProviderConfig()
		fmt.Print("\r\n")
		log.Fatalf("error normalizing generated resource addresses: %v", err)
	}

	if importsOnly {
		if err := removeResourceConfiguration(configFolder); err != nil {
			CleanupProviderConfig()
			log.Fatalf("error preparing imports-only output: %v", err)
		}
	}

	if _, err := runTfCmdInDirectory(configFolder, "init", "-backend=false"); err != nil {
		CleanupProviderConfig()
		log.Fatalf("error initializing Terraform: %v", err)
	}
	if _, err := runTfCmdInDirectory(configFolder, "fmt", "-recursive", "-list=false"); err != nil {
		CleanupProviderConfig()
		log.Fatalf("error running Terraform fmt: %v", err)
	}

	output.StopSpinner(spinner)
}

// NewQueryDirectory allocates private scratch space with explicit cleanup
// ownership instead of reusing predictable resource-config directories.
func NewQueryDirectory(level string) (*files.TempDirectory, error) {
	if providerDirectory == nil {
		ConfigureProvider(level)
	}
	directory, err := files.NewTempDirectory("btptf-query-")
	if err != nil {
		return nil, err
	}
	if err := files.CopyFile(providerDirectory.Root, directory.Root, "provider.tf", os.O_TRUNC); err != nil {
		_ = directory.Cleanup()
		return nil, err
	}
	return directory, nil
}

func ExecPostQuerySteps(tempConfigDir *files.TempDirectory, targetConfigDir *os.Root, targetResourceFileName string, resourceNameLong string, importsOnly bool) {

	spinner := output.StartSpinner("generating Terraform configuration for " + resourceNameLong)

	if !importsOnly {
		if _, err := runTfCmdInDirectory(tempConfigDir.Root, "fmt", "-recursive", "-list=false"); err != nil {
			CleanupTempFiles(tempConfigDir)
			log.Fatalf("error formatting Terraform configuration for %s: %v", resourceNameLong, err)
		}
	}

	err := mergeTfConfig(targetConfigDir, targetResourceFileName, tempConfigDir.Root, resourceNameLong)
	if err != nil {
		CleanupTempFiles(tempConfigDir)
		fmt.Print("\r\n")
		log.Fatalf("error merging Terraform configuration for %s: %v", resourceNameLong, err)
	}

	CleanupTempFiles(tempConfigDir)

	output.StopSpinner(spinner)

	fmt.Println(output.ColorStringGrey("   temporary files deleted"))
}

func CleanupTempFiles(directory *files.TempDirectory) {
	if err := directory.Cleanup(); err != nil {
		log.Printf("error deleting temporary files: %v", err)
	}
}

func mergeTfConfig(configFolder *os.Root, fileName string, resourceConfigFolder *os.Root, resourceName string) error {
	if err := files.ValidateFilename(fileName); err != nil {
		return err
	}
	sourceFile, err := resourceConfigFolder.Open(fileName)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("open resource configuration: %w", err)
	}
	defer func() { _ = sourceFile.Close() }()
	targetFile, err := configFolder.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("open target configuration: %w", err)
	}
	header := "\n###\n# Resource: " + resourceName + "\n###\n"
	_, writeErr := io.WriteString(targetFile, header)
	if writeErr == nil {
		_, writeErr = io.Copy(targetFile, sourceFile)
	}
	closeErr := targetFile.Close()
	if err := errors.Join(writeErr, closeErr); err != nil {
		return err
	}
	return files.CopyImportFiles(resourceConfigFolder, configFolder)
}

func recreateExistingConfigDir(work *os.Root, name string) error {
	if err := files.ValidateConfigDir(name); err != nil {
		return err
	}
	existing, err := openConfigDir(work, name)
	if err != nil {
		return err
	}
	_ = existing.Close()
	if err := work.RemoveAll(name); err != nil {
		return fmt.Errorf("remove configuration directory: %w", err)
	}
	return work.Mkdir(name, 0700)
}

// A symlink to the working directory itself is not a strict descendant either.
func openConfigDir(work *os.Root, name string) (*os.Root, error) {
	root, err := work.OpenRoot(name)
	if err != nil {
		return nil, err
	}
	base, err := work.Stat(".")
	if err != nil {
		_ = root.Close()
		return nil, err
	}
	target, err := root.Stat(".")
	if err != nil {
		_ = root.Close()
		return nil, err
	}
	if os.SameFile(base, target) {
		_ = root.Close()
		return nil, fmt.Errorf("configuration directory must not resolve to the working directory")
	}
	return root, nil
}
