package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	output "github.com/SAP/terraform-exporter-btp/pkg/output"
	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
)

func exportByJson(subaccount string, directory string, jsonfile string, resourceFile string, configDir string) {
	jsonFile, err := os.Open(jsonfile)
	if err != nil {
		log.Fatalf("error opening JSON file with resources: %v", err)
	}

	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf("error reading JSON file: %v", err)
	}

	var resources tfutils.BtpResources

	err = json.Unmarshal(byteValue, &resources)
	if err != nil {
		log.Fatalf("error unmarshalling JSON file: %v", err)
	}

	var resNames []string

	for i := 0; i < len(resources.BtpResources); i++ {
		resNames = append(resNames, resources.BtpResources[i].Name)
	}

	if len(resNames) == 0 {
		fmt.Println(output.ColorStringCyan("no resource needs to be exported"))
		return
	}

	tfutils.SetupConfigDir(configDir, true)

	for _, resName := range resNames {
		var value []string
		for _, temp := range resources.BtpResources {
			if temp.Name == resName {
				value = temp.Values
			}
		}
		if len(value) != 0 {
			generateConfigForResource(resName, value, subaccount, configDir, resourceFile)
		}
	}

	tfutils.FinalizeTfConfig(configDir)
	tfutils.CleanupProviderConfig()
}
