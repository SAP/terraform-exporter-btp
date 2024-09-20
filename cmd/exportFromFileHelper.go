package cmd

import (
	"btptfexport/output"
	"btptfexport/tfutils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"
)

func exportFromFile(subaccount string, jsonfile string, resourceFile string, configDir string) {
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

func generateConfigForResource(resource string, values []string, subaccount string, configDir string, resourceFileName string) {
	tempConfigDir := resource + "-config"
	techResourceNameLong := strings.ToUpper(tfutils.TranslateResourceParamToTechnicalName(resource))

	tfutils.ExecPreExportSteps(tempConfigDir)
	// Export must be done for each resource individually
	switch resource {
	case tfutils.CmdSubaccountParameter:
		exportSubaccount(subaccount, tempConfigDir, values)
	case tfutils.CmdEntitlementParameter:
		exportSubaccountEntitlements(subaccount, tempConfigDir, values)
	case tfutils.CmdEnvironmentInstanceParameter:
		exportSubaccountEnvironmentInstances(subaccount, tempConfigDir, values)
	case tfutils.CmdSubscriptionParameter:
		exportSubaccountSubscriptions(subaccount, tempConfigDir, values)
	case tfutils.CmdTrustConfigurationParameter:
		exportSubaccountTrustConfigurations(subaccount, tempConfigDir, values)
	}

	tfutils.ExecPostExportSteps(tempConfigDir, configDir, resourceFileName, techResourceNameLong)
}

func isSubset(superSet []string, subset []string) (string, bool) {
	for _, value := range subset {
		if !slices.Contains(superSet, value) {
			return value, false
		}
	}
	return "", true
}
