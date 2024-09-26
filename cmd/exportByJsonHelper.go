package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	output "github.com/SAP/terraform-exporter-btp/output"
	tfutils "github.com/SAP/terraform-exporter-btp/tfutils"
)

func exportByJson(subaccount string, jsonfile string, resourceFile string, configDir string) {
	jsonFile, err := os.Open(jsonfile)

	if err != nil {
		log.Fatalf("Error: %v", err)
		return
	}

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	var resources tfutils.BtpResources

	err = json.Unmarshal(byteValue, &resources)
	if err != nil {
		log.Fatalf("error in unmarshall: %v", err)
		return
	}

	var resNames []string

	for i := 0; i < len(resources.BtpResources); i++ {
		resNames = append(resNames, resources.BtpResources[i].Name)
	}
	if len(resNames) == 0 {
		fmt.Println(output.ColorStringCyan("No resource needs to be exported"))
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
}