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
	"github.com/SAP/terraform-exporter-btp/pkg/resume"
	tfcleantypes "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	tfcleanorchestrator "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/orchestrator"
	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
)

// maxJSONInventorySize bounds how many bytes are read from the inventory file
// as a safeguard in addition to the up-front file-size check.
const maxJSONInventorySize = 5 * 1024 * 1024 // 5 MB

func exportByJson(subaccount string, directory string, organization string, jsonfile string, resourceFile string, configDir string, backendConfig tfutils.BackendConfig) {
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

	defer func() {
		if tempErr := jsonFile.Close(); tempErr != nil {
			fmt.Print("\r\n")
			log.Fatalf("error closing JSON file with resources: %v", tempErr)
		}
	}()

	byteValue, err := io.ReadAll(io.LimitReader(jsonFile, maxJSONInventorySize))
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

	level, _ := tfutils.GetExecutionLevelAndId(subaccount, directory, organization, "")

	allowedResources := tfutils.GetValidResourcesByLevel(level)

	for i := range resources.BtpResources {

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

	exportLog, _ := resume.GetExistingExportLog(configDir)
	fullExportLog, _ := resume.GetExistingExportLogComplete(configDir)

	for _, resName := range resNames {
		// Check if the resource is already exported
		if len(exportLog) > 0 && slices.Contains(exportLog, resName) {
			// Skip the resource if it is already exported
			continue
		}

		var value []string
		for _, temp := range resources.BtpResources {
			if temp.Name == resName {
				value = temp.Values
			}
		}

		if len(value) != 0 {
			resourceType, count := generateConfigForResource(resName, value, subaccount, directory, organization, "", configDir, resourceFile)
			resultStore[resourceType] = count
			_ = resume.WriteExportLog(configDir, resName, resourceType, count)
		}
	}

	levelIds := tfcleantypes.LevelIds{
		SubaccountId: subaccount,
		DirectoryId:  directory,
		CfOrgId:      organization,
	}

	tfcleanorchestrator.CleanUpGeneratedCode(configDir, level, levelIds, &resultStore, backendConfig)
	tfutils.FinalizeTfConfig(configDir)
	generateNextStepsDocument(configDir, subaccount, directory, organization, "")
	_ = resume.RemoveExportLog(configDir)
	resultStoreNew := resume.MergeSummaryTable(resultStore, fullExportLog)
	output.RenderSummaryTable(resultStoreNew)
	tfutils.CleanupProviderConfig()
}
