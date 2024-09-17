package cmd

import (
	"btptfexport/tfutils"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/viper"
	"github.com/theckman/yacspin"
)

var TmpFolder string

func runTerraformCommand(args ...string) error {

	debug := viper.GetViper().GetBool("debug")
	cmd := exec.Command("terraform", args...)
	if debug {
		cmd.Stdout = os.Stdout
	} else {
		cmd.Stdout = nil
	}

	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func generateConfig(resourceFileName string, configFolder string, isMainCmd bool, resourceNameLong string) {

	var spinner *yacspin.Spinner
	var err error

	if isMainCmd {
		// We must distinguish if the command is run from a main command or via delegation from helper functions
		spinner, err = startSpinner("create Terraform configuration for " + resourceNameLong)
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
	// Initialize Terraform
	if err := runTerraformCommand("init"); err != nil {
		log.Fatalf("error initializing Terraform: %v", err)
		return
	}

	// Execute Terraform plan
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
		err = stopSpinner(spinner)
		if err != nil {
			log.Fatalf("error: %v", err)
			return
		}
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

func configureProvider() {
	tmpdir, err := os.MkdirTemp("", "provider.tf")
	if err != nil {
		panic(err)
	}
	TmpFolder = tmpdir
	abspath := filepath.Join(tmpdir, "provider.tf")
	fmt.Println(abspath)
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

	err = tfutils.CreateFileWithContent(abspath, providerContent)
	if err != nil {
		log.Fatalf("create file %s failed!", abspath)
		return
	}

}

func setupConfigDir(configFolder string, isMainCmd bool) {

	if isMainCmd {
		message := "set up config directory \"" + configFolder + "\""
		fmt.Println(color.HiBlackString(message))
	}

	if len(TmpFolder) == 0 {
		configureProvider()
	}
	curWd, err := os.Getwd()
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	exist, err := exists(filepath.Join(curWd, configFolder))
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

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func mergeTfConfig(configFolder string, fileName string, resourceConfigFolder string, resourceName string) {

	currentDir, err := os.Getwd()

	if err != nil {
		log.Fatalf("error getting current directory: %v", err)
		return
	}

	sourceConfigPath := filepath.Join(currentDir, resourceConfigFolder, fileName)

	// Check if the source file exists
	exist, err := exists(sourceConfigPath)
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

	exist, err = exists(targetConfigPath)
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

	copyImportFiles(resourceConfigFolder, configFolder)
}

func copyImportFiles(srcDir, destDir string) {
	// Find all files ending with "_import.tf" in the source directory
	files, err := filepath.Glob(filepath.Join(srcDir, "*_import.tf"))
	if err != nil {
		log.Fatalf("error finding files: %v", err)
	}

	// Copy each file to the destination directory
	for _, srcFile := range files {
		destFile := filepath.Join(destDir, filepath.Base(srcFile))

		err := copyFile(srcFile, destFile)
		if err != nil {
			log.Printf("error copying file %s to %s: %v", srcFile, destFile, err)
		}
	}
}

func copyFile(src, dest string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func deleteSourceFolder(srcDir string) {
	err := os.RemoveAll(srcDir)
	if err != nil {
		log.Fatalf("error deleting source folder %s: %v", srcDir, err)
	}
}

func finalizeTfConfig(configFolder string) {

	spinner, err := startSpinner("finalizing Terraform configuration")
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

	err = stopSpinner(spinner)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}
}

// Conveniecce functions that wrap repetitive steps
func execPreExportSteps(tempConfigDir string) {
	setupConfigDir(tempConfigDir, false)
}

func execPostExportSteps(tempConfigDir string, targetConfigDir string, targetResourceFileName string, resourceNameLong string) {

	spinner, err := startSpinner("create Terraform configuration for " + resourceNameLong)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	generateConfig(targetResourceFileName, tempConfigDir, false, resourceNameLong)
	mergeTfConfig(targetConfigDir, targetResourceFileName, tempConfigDir, resourceNameLong)
	deleteSourceFolder(tempConfigDir)

	err = stopSpinner(spinner)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}
}
