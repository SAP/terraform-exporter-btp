package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/SAP/terraform-exporter-btp/pkg/files"
	output "github.com/SAP/terraform-exporter-btp/pkg/output"
	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
)

func exportByJson(subaccount string, directory string, organization string, jsonfile string, resourceFile string, configDir string) {
	// Check if file size is valid
	validFileSize, err := files.IsFileSizeValid(jsonfile)

	if err != nil {
		fmt.Print("\r\n")
		log.Fatalf("error checking file size: %v", err)
	}

	if !validFileSize {
		fmt.Print("\r\n")
		log.Fatalf("file size exceeds the limit of 5MB")
	}

	jsonFile, err := os.Open(jsonfile)
	if err != nil {
		fmt.Print("\r\n")
		log.Fatalf("error opening JSON file with resources: %v", err)
	}

	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Print("\r\n")
		log.Fatalf("error reading JSON file: %v", err)
	}

	var resources tfutils.BtpResources

	// Unmarshal validates basic JSON syntax and structure
	err = json.Unmarshal(byteValue, &resources)
	if err != nil {
		fmt.Print("\r\n")
		log.Fatalf("error unmarshalling JSON file: %v", err)
	}

	var resNames []string

	level, _ := tfutils.GetExecutionLevelAndId(subaccount, directory, organization)

	allowedResources := tfutils.GetValidResourcesByLevel(level)

	for i := 0; i < len(resources.BtpResources); i++ {

		if !(slices.Contains(allowedResources, resources.BtpResources[i].Name)) {

			allowedResourceList := strings.Join(allowedResources, ", ")
			fmt.Print("\r\n")
			log.Fatal("please check the resources provided in the JSON file. Currently supported resources are " + allowedResourceList + ".")
		}

		resNames = append(resNames, resources.BtpResources[i].Name)
	}

	if len(resNames) == 0 {
		fmt.Println(output.ColorStringCyan("no resource needs to be exported"))
		return
	}

	tfutils.SetupConfigDir(configDir, true, level)
	resultStore := make(map[string]int)

	for _, resName := range resNames {
		var value []string
		for _, temp := range resources.BtpResources {
			if temp.Name == resName {
				value = temp.Values
			}
		}
		if len(value) != 0 {
			resourceType, count := generateConfigForResource(resName, value, subaccount, directory, organization, configDir, resourceFile)
			resultStore[resourceType] = count
		}
	}

	tfutils.FinalizeTfConfig(configDir)
	generateNextStepsDocument(configDir, subaccount, directory, organization)
	output.RenderSummaryTable(resultStore)
	tfutils.CleanupProviderConfig()
}
