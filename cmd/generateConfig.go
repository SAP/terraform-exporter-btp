package cmd

import (
	"btptfexporter/tfutils"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var TmpFolder string

func runTerraformCommand(args ...string) error {
	cmd := exec.Command("terraform", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func generateConfig(resourceFileName string, configFolder string) {

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

	fmt.Println("Terraform config successfully created")
	cleanup()
}

func cleanup() {
	err := os.RemoveAll(TmpFolder)
	if err != nil {
		log.Fatalf("error deleting temp files: %v", err)
	}
}

func configureProvider() {
	//tmpdir, err := ioutil.TempDir("/tmp", "sampledir")
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
	//providerContent := "terraform {\nrequired_providers {\nbtp = {\nsource  = \"SAP/btp\"\nversion = \"1.4.0\"\n}\n}\n}\n\nprovider \"btp\" {\ncli_server_url=\"https://cpcli.cf.eu10.hana.ondemand.com\"\nglobalaccount = \"terraformintprod\"\n}"

	err = tfutils.CreateFileWithContent(abspath, providerContent)
	if err != nil {
		log.Fatalf("create file %s failed!", abspath)
		return
	}

}

func setupConfigDir(configFolder string) {
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
