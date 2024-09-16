package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"
)

type Resource struct {
	Name   string
	Values []string
}

type ResourcesArr struct {
	Btp_resources []Resource
}

func exportFromFile(subaccount string, jsonfile string, resourceFile string, configDir string) {
	jsonFile, err := os.Open(jsonfile)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	var resources ResourcesArr

	err = json.Unmarshal(byteValue, &resources)

	if err != nil {
		log.Fatalf("error in unmarshall: %v", err)
		return
	}

	var resNames []string

	for i := 0; i < len(resources.Btp_resources); i++ {
		resNames = append(resNames, resources.Btp_resources[i].Name)
	}
	if len(resNames) == 0 {
		fmt.Println("No resource needs to be export")
		return
	}

	setupConfigDir(configDir)

	for _, resName := range resNames {
		var value []string
		for _, temp := range resources.Btp_resources {
			if temp.Name == resName {
				value = temp.Values
			}
		}
		if len(value) != 0 {
			generateConfigForResource(resName, value, subaccount, configDir, resourceFile)
		}
	}

	finalizeTfConfig(configDir)
}

func generateConfigForResource(resource string, values []string, subaccount string, configDir string, resourceFileName string) {
	tempConfigDir := resource + "-config"
	techResourceNameLong := strings.ToUpper(string(translateResourceParamToTechnicalName(resource)))

	execPreExportSteps(tempConfigDir)
	// Export must be done for each resource individually
	switch resource {
	case CmdSubaccountParameter:
		exportSubaccount(subaccount, tempConfigDir, values)
	case CmdEntitlementParameter:
		exportSubaccountEntitlements(subaccount, tempConfigDir, values)
	case CmdEnvironmentInstanceParameter:
		exportSubaccountEnvironmentInstances(subaccount, tempConfigDir, values)
	case CmdSubscriptionParameter:
		exportSubaccountSubscriptions(subaccount, tempConfigDir, values)
	case CmdTrustConfigurationParameter:
		exportSubaccountTrustConfigurations(subaccount, tempConfigDir, values)
	}

	execPostExportSteps(tempConfigDir, configDir, resourceFileName, techResourceNameLong)
}

func isSubset(superSet []string, subset []string) (string, bool) {
	for _, value := range subset {
		if !slices.Contains(superSet, value) {
			return value, false
		}
	}
	return "", true
}
