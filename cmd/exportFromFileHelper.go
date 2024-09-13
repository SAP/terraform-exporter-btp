package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
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
	if resource == "environment-instances" {
		execPreExportSteps("saenvinstanceconf")
		exportEnvironmentInstances(subaccount, "saenvinstanceconf", values)
		execPostExportSteps("saenvinstanceconf", configDir, resourceFileName, "SUBACCOUNT ENVIRONMENT INSTANCES")
	}
	if resource == "subaccount" {
		execPreExportSteps("saconf")
		exportSubaccount(subaccount, "saconf", values)
		execPostExportSteps("saconf", configDir, resourceFileName, "SUBACCOUNT")
	}
	if resource == "entitlements" {
		execPreExportSteps("saentitlementconf")
		exportSubaccountEntitlements(subaccount, "saentitlementconf", values)
		execPostExportSteps("saentitlementconf", configDir, resourceFileName, "SUBACCOUNT ENTITLEMENTS")
	}
	if resource == "subscriptions" {
		execPreExportSteps("sasubscriptionconf")
		exportSubaccountSubscriptions(subaccount, "sasubscriptionconf", values)
		execPostExportSteps("sasubscriptionconf", configDir, resourceFileName, "SUBACCOUNT  SUBSCRIPTIONS")
	}
	if resource == "trust-configurations" {
		execPreExportSteps("satrustconf")
		exportTrustConfigurations(subaccount, "satrustconf", values)
		execPostExportSteps("satrustconf", configDir, resourceFileName, "SUBACCOUNT TRUST CONFIGURATIONS")
	}
}

func isSubset(superSet []string, subset []string) (string, bool) {
	for _, value := range subset {
		if !slices.Contains(superSet, value) {
			return value, false
		}
	}
	return "", true
}
