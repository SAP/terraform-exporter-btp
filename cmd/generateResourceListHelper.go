package cmd

import (
	"btptfexport/files"
	"btptfexport/output"
	"btptfexport/tfutils"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func getResourcesInfo(subaccount string, fileName string, resources string) {
	if len(resources) == 0 {
		log.Fatal("please provide the btp resources you want to get using --resources flag or provide 'all' to get all resources")
		return
	}

	tfutils.ConfigureProvider()

	var inputRes []string

	if resources == "all" {
		inputRes = AllowedResources
	} else {
		inputRes = strings.Split(resources, ",")

		for _, resource := range inputRes {
			if !(slices.Contains(AllowedResources, resource)) {
				log.Fatal("please check the resource provided. Currently supported resources are subaccount, entitlements, subscriptions, environment-instances and trust-configurations. Provide 'all' to check for all resources")
				return
			}
		}
	}

	spinner, err := output.StartSpinner("Collecting resources of subaccount")
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	result, err := tfutils.ReadDataSources(subaccount, inputRes)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("error: %s", err)
		return
	}

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("error getting current directory: %s", err)
		return
	}

	dataBlockFile := filepath.Join(currentDir, fileName)
	err = files.CreateFileWithContent(dataBlockFile, string(jsonBytes))
	if err != nil {
		log.Fatalf("create file %s failed!", dataBlockFile)
		return
	}

	err = output.StopSpinner(spinner)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

}
