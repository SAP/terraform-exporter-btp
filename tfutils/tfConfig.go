package tfutils

import (
	"btptfexport/files"
	"btptfexport/output"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"github.com/theckman/yacspin"
)

var TmpFolder string

func GenerateConfig(resourceFileName string, configFolder string, isMainCmd bool, resourceNameLong string) {

	var spinner *yacspin.Spinner
	var err error

	if isMainCmd {
		// We must distinguish if the command is run from a main command or via delegation from helper functions
		spinner, err = output.StartSpinner("generating Terraform configuration for " + resourceNameLong)
		if err != nil {
			log.Fatalf("error: %v", err)
			return
		}
	}

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("error getting current directory: %v", err)
		return
	}

	terraformConfigPath := filepath.Join(currentDir, configFolder)
	err = os.Chdir(terraformConfigPath)
	if err != nil {
		log.Fatalf("error changing directory to %s: %v \n", terraformConfigPath, err)
		return
	}

	if err := runTerraformCommand("init"); err != nil {
		log.Fatalf("error running Terraform init: %v", err)
		return
	}

	planOption := "--generate-config-out=" + resourceFileName
	if err := runTerraformCommand("plan", planOption); err != nil {
		log.Fatalf("error running Terraform plan: %v", err)
		return
	}

	if err := runTerraformCommand("fmt", "-recursive", "-list=false"); err != nil {
		log.Fatalf("error running Terraform fmt: %v", err)
		return
	}

	cleanup()

	//Switch back to the original directory
	err = os.Chdir(currentDir)
	if err != nil {
		log.Fatalf("error changing directory to %s: %v \n", currentDir, err)
		return
	}

	if isMainCmd {
		err = output.StopSpinner(spinner)
		if err != nil {
			log.Fatalf("error: %v", err)
			return
		}
		fmt.Println(output.ColorStringGrey("   temporary files deleted"))
	}
}

func ConfigureProvider() {
	tmpdir, err := os.MkdirTemp("", "provider.tf")
	if err != nil {
		panic(err)
	}

	TmpFolder = tmpdir
	abspath := filepath.Join(tmpdir, "provider.tf")

	debug := viper.GetViper().GetBool("debug")

	if debug {
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
			log.Fatal("set BTP_USERNAME and BTP_PASSWORD environment variable or enable SSO for login.")
			os.Exit(0)
		}
	}

	if len(strings.TrimSpace(globalAccount)) == 0 {
		log.Fatal("global account not set. set BTP_GLOBALACCOUNT environment variable to set global account")
		os.Exit(0)
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
		log.Fatalf("create file %s failed!", abspath)
		return
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
		log.Fatalf("error: %v", err)
		return
	}

	exist, err := files.Exists(filepath.Join(curWd, configFolder))
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	if !exist {
		fullpath := filepath.Join(curWd, configFolder)
		err = os.Mkdir(fullpath, 0700)
		if err != nil {
			log.Fatalf("error: %v", err)
			return
		}
	} else {
		fmt.Print("config directory already exist. Do you want to continue? If yes then generated files will be overwritten if existing (Y/N): ")
		var choice string
		_, err = fmt.Scanln(&choice)
		if err != nil {
			log.Fatalf("error: %v", err)
			return
		}
		if strings.ToUpper(choice) == "N" {
			os.Exit(0)
		} else if strings.ToUpper(choice) == "Y" {
			fmt.Println("existing directory will be used. It can overwrite some files ")
		} else {
			fmt.Println("invalid input. Exiting the process")
			os.Exit(0)
		}
	}

	sourceFile, err := os.Open(TmpFolder + "/provider.tf")
	if err != nil {
		log.Fatalf("failed to open source file: %v", err)
		return
	}
	defer sourceFile.Close()

	fullpath := filepath.Join(curWd, configFolder)

	destinationFile, err := os.Create(fullpath + "/provider.tf")
	if err != nil {
		log.Fatalf("failed to create destination file: %v", err)
		return
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		log.Fatalf("failed to copy file: %v", err)
		return
	}
}

func cleanup() {
	err := os.RemoveAll(TmpFolder)
	if err != nil {
		log.Fatalf("error deleting temp files: %v", err)
	}

	// Cleanup temporary folder variable
	TmpFolder = ""
}

func FinalizeTfConfig(configFolder string) {

	spinner, err := output.StartSpinner("finalizing Terraform configuration")
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("error getting current directory: %v", err)
		return
	}

	terraformConfigPath := filepath.Join(currentDir, configFolder)

	err = os.Chdir(terraformConfigPath)
	if err != nil {
		log.Fatalf("error changing directory to %s: %v \n", terraformConfigPath, err)
		return
	}

	if err := runTerraformCommand("init"); err != nil {
		log.Fatalf("error initializing Terraform: %v", err)
		return
	}

	if err := runTerraformCommand("fmt", "-recursive", "-list=false"); err != nil {
		log.Fatalf("error running Terraform fmt: %v", err)
		return
	}

	//Switch back to the original directory
	err = os.Chdir(currentDir)
	if err != nil {
		log.Fatalf("error changing directory to %s: %v \n", currentDir, err)
		return
	}

	err = output.StopSpinner(spinner)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}
}

// Convenience functions that wrap repetitive steps
func ExecPreExportSteps(tempConfigDir string) {
	SetupConfigDir(tempConfigDir, false)
}

func ExecPostExportSteps(tempConfigDir string, targetConfigDir string, targetResourceFileName string, resourceNameLong string) {

	spinner, err := output.StartSpinner("generating Terraform configuration for " + resourceNameLong)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	GenerateConfig(targetResourceFileName, tempConfigDir, false, resourceNameLong)
	mergeTfConfig(targetConfigDir, targetResourceFileName, tempConfigDir, resourceNameLong)
	files.DeleteSourceFolder(tempConfigDir)

	err = output.StopSpinner(spinner)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	fmt.Println(output.ColorStringGrey("   temporary files deleted"))
}

func mergeTfConfig(configFolder string, fileName string, resourceConfigFolder string, resourceName string) {

	currentDir, err := os.Getwd()

	if err != nil {
		log.Fatalf("error getting current directory: %v", err)
		return
	}

	sourceConfigPath := filepath.Join(currentDir, resourceConfigFolder, fileName)

	// Check if the source file exists
	exist, err := files.Exists(sourceConfigPath)
	if err != nil {
		log.Fatalf("error checking if source directory exists: %v", err)
	}

	if !exist {
		// Nothing to do as the source file does not exist
		return
	}

	sourceFile, err := os.Open(sourceConfigPath)
	if err != nil {
		log.Fatalf("error opening resource config file: %v", err)
	}
	defer sourceFile.Close()

	targetConfigPath := filepath.Join(currentDir, configFolder, fileName)

	exist, err = files.Exists(targetConfigPath)
	if err != nil {
		log.Fatalf("error checking if target directory exists: %v", err)
	}

	if !exist {
		// In the first run we must create the file if it does not exist
		_, err := os.Create(targetConfigPath)
		if err != nil {
			log.Fatalf("error creating target configuration file: %v", err)
		}
	}

	targetFile, err := os.OpenFile(targetConfigPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("error opening target configuration file: %v", err)
	}
	defer targetFile.Close()

	headerTemplate := `
###
# Resource: ` + resourceName + `
###
`
	if _, err := targetFile.Write([]byte(headerTemplate)); err != nil {
		log.Fatalf("error adding header line to target file: %v", err)
	}

	if _, err := io.Copy(targetFile, sourceFile); err != nil {
		log.Fatalf("error copying resource file to target file: %v", err)
	}

	files.CopyImportFiles(resourceConfigFolder, configFolder)
}
