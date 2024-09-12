package cmd

import (
	"btptfexporter/tfutils"
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

	configureProvider()
	var btpSubaccountResources []Resource
	btpResources := []string{"subaccount", "entitlements", "subscriptions", "environment-instances", "trust-configurations"}
	if resources == "all" {
		for _, resource := range btpResources {
			values, err := readDataSources(subaccount, resource)
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			data := Resource{Name: resource, Values: values}
			btpSubaccountResources = append(btpSubaccountResources, data)

		}
	} else {
		if len(resources) == 0 {
			log.Fatal("please provide btp resources you want to get using --resources flag")
			return
		}

		inputRes := strings.Split(resources, ",")

		for _, resource := range inputRes {
			if !(slices.Contains(btpResources, resource)) {
				log.Fatal("please check the resource provided. Currently supporte resources are subaccount, entitlements, subscriptions, environment-instances and trust-configurations. Provide 'all' to check for all resources")
				return
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
	if btpResource == "subaccount" {
		dataBlock, err := readSubaccountDataSource(subaccountID)
		if err != nil {
			fmt.Println("error getting data source:", err)
			return nil, err
		}
		err = tfutils.CreateFileWithContent(dataBlockFile, dataBlock)
		if err != nil {
			log.Fatalf("create file %s failed!", dataBlockFile)
			return nil, err
		}

		jsonBytes, err = GetTfStateData(TmpFolder, SubaccountType)
		if err != nil {
			log.Fatalf("error json.Marshal: %s", err)
			return nil, err
		}
	} else if btpResource == "entitlements" {
		dataBlock, err := readSubaccountEntilementsDataSource(subaccountID)
		if err != nil {
			fmt.Println("error getting data source:", err)
			return nil, err
		}

		err = tfutils.CreateFileWithContent(dataBlockFile, dataBlock)
		if err != nil {
			log.Fatalf("create file %s failed!", dataBlockFile)
			return nil, err
		}

		jsonBytes, err = GetTfStateData(TmpFolder, SubaccountEntitlementType)
		if err != nil {
			log.Fatalf("error json.Marshal: %s", err)
			return nil, err
		}
	} else if btpResource == "subscriptions" {
		dataBlock, err := readSubaccountSubscriptionDataSource(subaccountID)
		if err != nil {
			fmt.Println("error getting data source:", err)
			return nil, err
		}

		err = tfutils.CreateFileWithContent(dataBlockFile, dataBlock)
		if err != nil {
			log.Fatalf("create file %s failed!", dataBlockFile)
			return nil, err
		}

		jsonBytes, err = GetTfStateData(TmpFolder, SubaccountSubscriptionType)
		if err != nil {
			log.Fatalf("error json.Marshal: %s", err)
			return nil, err
		}
	} else if btpResource == "environment-instances" {
		dataBlock, err := readDataSource(subaccountID)
		if err != nil {
			fmt.Println("error getting data source:", err)
			return nil, err
		}

		err = tfutils.CreateFileWithContent(dataBlockFile, dataBlock)
		if err != nil {
			log.Fatalf("create file %s failed!", dataBlockFile)
			return nil, err
		}

		jsonBytes, err = GetTfStateData(TmpFolder, EnvironmentInstanceType)
		if err != nil {
			log.Fatalf("error json.Marshal: %s", err)
			return nil, err
		}
	} else if btpResource == "trust-configurations" {
		dataBlock, err := readSubaccountTrustConfigurationsDataSource(subaccountID)
		if err != nil {
			fmt.Println("error getting data source:", err)
			return nil, err
		}

		err = tfutils.CreateFileWithContent(dataBlockFile, dataBlock)
		if err != nil {
			log.Fatalf("create file %s failed!", dataBlockFile)
			return nil, err
		}

		jsonBytes, err = GetTfStateData(TmpFolder, SubaccountTrustConfigurationType)
		if err != nil {
			log.Fatalf("error json.Marshal: %s", err)
			return nil, err
		}
	}
	var data map[string]interface{}
	err := json.Unmarshal(jsonBytes, &data)
	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	}
	var stringArr []string
	if btpResource == "subaccount" {
		stringArr = []string{fmt.Sprintf("%v", data["name"])}
	} else if btpResource == "entitlements" {
		for key, _ := range data {

			key := strings.Replace(key, ":", "_", -1)
			stringArr = append(stringArr, key)
		}
	} else if btpResource == "subscriptions" {
		subscriptions := data["values"].([]interface{})
		for _, value := range subscriptions {
			subscription := value.(map[string]interface{})
			if fmt.Sprintf("%v", subscription["state"]) != "NOT_SUBSCRIBED" {
				stringArr = append(stringArr, fmt.Sprintf("%v", subscription["app_name"])+"_"+fmt.Sprintf("%v", subscription["plan_name"]))
			}
		}
	} else if btpResource == "environment-instances" {
		environmentInstances := data["values"].([]interface{})
		for _, value := range environmentInstances {
			environmentInstance := value.(map[string]interface{})
			stringArr = append(stringArr, fmt.Sprintf("%v", environmentInstance["environment_type"]))
		}
	} else if btpResource == "trust-configurations" {
		trusts := data["values"].([]interface{})
		for _, value := range trusts {
			trust := value.(map[string]interface{})
			stringArr = append(stringArr, fmt.Sprintf("%v", trust["origin"]))
		}
	}
	return stringArr, nil
}
