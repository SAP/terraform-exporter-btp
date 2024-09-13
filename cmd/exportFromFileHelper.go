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
			generateConfigForResource(resName, value, subaccount, configDir)
		}
	}

	generateConfig(resourceFile, configDir)
}

func generateConfigForResource(resource string, values []string, subaccount string, configDir string) {
	if resource == "environment-instances" {
		exportEnvironmentInstances(subaccount, configDir, values)
	}
	if resource == "subaccount" {
		exportSubaccount(subaccount, configDir, values)
	}
	if resource == "entitlements" {
		exportEnvironmentInstances(subaccount, configDir, values)
	}
	if resource == "subscriptions" {
		exportSubaccountSubscriptions(subaccount, configDir, values)
	}
	if resource == "trust-configurations" {
		exportTrustConfigurations(subaccount, configDir, values)
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
