package cmd

import (
	"btptfexport/tfutils"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func getResourcesInfo(subaccount string, fileName string, resources string) {
	type Resource struct {
		Name   string
		Values []string
	}

	type resourcesArr struct {
		Btp_resources []Resource
	}

	var btpSubaccountResources []Resource

	if len(resources) == 0 {
		log.Fatal("please provide the btp resources you want to get using --resources flag or provide 'all' to get all resources")
		return
	}

	configureProvider()

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

	for _, resource := range inputRes {

		values, err := readDataSources(subaccount, resource)
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		data := Resource{Name: resource, Values: values}
		btpSubaccountResources = append(btpSubaccountResources, data)

	}

	res2 := resourcesArr{Btp_resources: btpSubaccountResources}
	jsonBytes, err := json.MarshalIndent(res2, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("error getting current directory:", err)
		return
	}
	dataBlockFile := filepath.Join(currentDir, fileName)
	err = tfutils.CreateFileWithContent(dataBlockFile, string(jsonBytes))
	if err != nil {
		log.Fatalf("create file %s failed!", dataBlockFile)
		return
	}

}

func readDataSources(subaccountID string, btpResource string) ([]string, error) {
	dataBlockFile := filepath.Join(TmpFolder, "main.tf")
	var jsonBytes []byte

	btpResourceType := translateResourceParamToTechnicalName(btpResource)

	dataBlock, err := readDataSource(subaccountID, btpResourceType)
	if err != nil {
		fmt.Println("error getting data source:", err)
		return nil, err
	}

	err = tfutils.CreateFileWithContent(dataBlockFile, dataBlock)
	if err != nil {
		log.Fatalf("create file %s failed!", dataBlockFile)
		return nil, err
	}

	jsonBytes, err = GetTfStateData(TmpFolder, btpResourceType)
	if err != nil {
		log.Fatalf("error json.Marshal: %s", err)
		return nil, err
	}

	var data map[string]interface{}

	err = json.Unmarshal(jsonBytes, &data)
	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	}

	return transformDataToStringArray(btpResource, data), nil
}

func transformDataToStringArray(btpResource string, data map[string]interface{}) []string {
	var stringArr []string

	switch btpResource {
	case CmdSubaccountParameter:
		stringArr = []string{fmt.Sprintf("%v", data["name"])}
	case CmdEntitlementParameter:
		for key := range data {
			key := strings.Replace(key, ":", "_", -1)
			stringArr = append(stringArr, key)
		}
	case CmdSubscriptionParameter:
		subscriptions := data["values"].([]interface{})
		for _, value := range subscriptions {
			subscription := value.(map[string]interface{})
			if fmt.Sprintf("%v", subscription["state"]) != "NOT_SUBSCRIBED" {
				stringArr = append(stringArr, fmt.Sprintf("%v", subscription["app_name"])+"_"+fmt.Sprintf("%v", subscription["plan_name"]))
			}
		}
	case CmdEnvironmentInstanceParameter:
		environmentInstances := data["values"].([]interface{})
		for _, value := range environmentInstances {
			environmentInstance := value.(map[string]interface{})
			stringArr = append(stringArr, fmt.Sprintf("%v", environmentInstance["environment_type"]))
		}
	case CmdTrustConfigurationParameter:
		trusts := data["values"].([]interface{})
		for _, value := range trusts {
			trust := value.(map[string]interface{})
			stringArr = append(stringArr, fmt.Sprintf("%v", trust["origin"]))
		}
	}
	return stringArr
}
