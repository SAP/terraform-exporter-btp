package cmd

import (
	"btptfexporter/tfutils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
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
		fmt.Println("err")
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

func generateConfigForResource(resource string, values []string, subaccout string, configDir string) {
	//fmt.Println(resource)
	//fmt.Println(values)
	if resource == "environment-instances" {
		getEnvInstanceConfig(values, subaccout, configDir)
	}
	if resource == "subaccount" {
		exportSubaccount(subaccout, configDir)
	}
	if resource == "entitlements" {
		getEntitlementConfig(values, subaccout, configDir)
	}
	if resource == "subscriptions" {
		getSubscriptionConfig(values, subaccout, configDir)
	}
	if resource == "trust-configurations" {
		getTrustConfig(values, subaccout, configDir)
	}
}

func getTrustConfig(values []string, subaccountID string, configDir string) {
	dataBlock, err := readSubaccountTrustConfigurationsDataSource(subaccountID)
	if err != nil {
		fmt.Println("error getting data source:", err)
		return
	}

	dataBlockFile := filepath.Join(TmpFolder, "main.tf")
	err = tfutils.CreateFileWithContent(dataBlockFile, dataBlock)
	if err != nil {
		log.Fatalf("create file %s failed!", dataBlockFile)
		return
	}

	jsonBytes, err := getTrustConfigurationsTfStateData(TmpFolder)
	if err != nil {
		log.Fatalf("error json.Marshal: %s", err)
		return
	}

	jsonString := string(jsonBytes)
	var data map[string]interface{}
	err = json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	importBlock, err := getImportBlock4(data, subaccountID, values)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("error getting current directory:", err)
		return
	}

	if len(importBlock) == 0 {
		log.Println("No trust configuration found for the given subaccount")
		return
	}

	importFileName := "subaccount_trust_configurations_import.tf"
	importFileName = filepath.Join(currentDir, configDir, importFileName)

	err = tfutils.CreateFileWithContent(importFileName, importBlock)
	if err != nil {
		log.Fatalf("create file %s failed!", dataBlockFile)
		return
	}

	log.Println("subaccount trust configuration has been exported. Please check " + configDir + " folder")
}

func getImportBlock4(data map[string]interface{}, subaccountId string, values []string) (string, error) {
	choice := "btp_subaccount_trust_configuration"
	resource_doc, err := tfutils.GetDocsForResource("SAP", "btp", "btp", "resources", choice, BtpProviderVersion, "github.com")
	if err != nil {
		log.Fatalf("read doc failed!")
		return "", err
	}

	var importBlock string
	trusts := data["values"].([]interface{})

	for x, value := range trusts {

		trust := value.(map[string]interface{})
		if slices.Contains(values, fmt.Sprintf("%v", trust["origin"])) {
			template := strings.Replace(resource_doc.Import, "<resource_name>", "trust"+fmt.Sprint(x), -1)
			template = strings.Replace(template, "<subaccount_id>", subaccountId, -1)
			template = strings.Replace(template, "<origin>", fmt.Sprintf("%v", trust["origin"]), -1)
			importBlock += template + "\n"
		}
	}

	return importBlock, nil
}

func getSubscriptionConfig(values []string, subaccountID string, configDir string) {
	dataBlock, err := readSubaccountSubscriptionDataSource(subaccountID)
	if err != nil {
		fmt.Println("error getting data source:", err)
		return
	}

	dataBlockFile := filepath.Join(TmpFolder, "main.tf")
	err = tfutils.CreateFileWithContent(dataBlockFile, dataBlock)
	if err != nil {
		log.Fatalf("create file %s failed!", dataBlockFile)
		return
	}

	jsonBytes, err := getSubscriptionsTfStateData(TmpFolder)
	if err != nil {
		log.Fatalf("error json.Marshal: %s", err)
		return
	}

	jsonString := string(jsonBytes)
	var data map[string]interface{}
	err = json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	importBlock, err := getImportBlock3(data, subaccountID, values)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("error getting current directory:", err)
		return
	}

	if len(importBlock) == 0 {
		log.Println("No subscription found for the given subaccount")
		return
	}

	importFileName := "subaccount_subscriptions_import.tf"
	importFileName = filepath.Join(currentDir, configDir, importFileName)

	err = tfutils.CreateFileWithContent(importFileName, importBlock)
	if err != nil {
		log.Fatalf("create file %s failed!", dataBlockFile)
		return
	}

	log.Println("subaccount subscriptions has been exported. Please check " + configDir + " folder")

}

func getImportBlock3(data map[string]interface{}, subaccountId string, values []string) (string, error) {
	choice := "btp_subaccount_subscription"
	resource_doc, err := tfutils.GetDocsForResource("SAP", "btp", "btp", "resources", choice, BtpProviderVersion, "github.com")
	if err != nil {
		log.Fatalf("read doc failed!")
		return "", err
	}

	var importBlock string
	subscriptions := data["values"].([]interface{})

	for _, value := range subscriptions {
		subscription := value.(map[string]interface{})
		if slices.Contains(values, fmt.Sprintf("%v", subscription["app_name"])+"_"+fmt.Sprintf("%v", subscription["plan_name"])) {
			template := strings.Replace(resource_doc.Import, "<resource_name>", strings.Replace(fmt.Sprintf("%v", subscription["app_name"]), "-", "_", -1), -1)
			template = strings.Replace(template, "<subaccount_id>", subaccountId, -1)
			template = strings.Replace(template, "<app_name>", fmt.Sprintf("%v", subscription["app_name"]), -1)
			template = strings.Replace(template, "<plan_name>", fmt.Sprintf("%v", subscription["plan_name"]), -1)
			importBlock += template + "\n"
		}
	}

	return importBlock, nil
}

func getEntitlementConfig(values []string, subaccountID string, configDir string) {
	dataBlock, err := readSubaccountEntilementsDataSource(subaccountID)
	if err != nil {
		fmt.Println("error getting data source:", err)
		return
	}

	dataBlockFile := filepath.Join(TmpFolder, "main.tf")
	err = tfutils.CreateFileWithContent(dataBlockFile, dataBlock)
	if err != nil {
		log.Fatalf("create file %s failed!", dataBlockFile)
		return
	}

	jsonBytes, err := getEntitlementTfStateData(TmpFolder)
	if err != nil {
		log.Fatalf("error json.Marshal: %s", err)
		return
	}

	jsonString := string(jsonBytes)
	var data map[string]interface{}
	err = json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	importBlock, err := getImportBlock2(data, subaccountID, values)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	if len(importBlock) == 0 {
		log.Println("No Entitlement found for the given subaccount")
		return
	}

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("error getting current directory:", err)
		return
	}

	importFileName := "subaccount_entitlements_import.tf"
	importFileName = filepath.Join(currentDir, configDir, importFileName)

	err = tfutils.CreateFileWithContent(importFileName, importBlock)
	if err != nil {
		log.Fatalf("create file %s failed!", dataBlockFile)
		return
	}

	log.Println("btp subaccount entitlements has been exported. Please check " + configDir + " folder")

}

func getImportBlock2(data map[string]interface{}, subaccountId string, values []string) (string, error) {

	choice := "btp_subaccount_entitlement"
	resource_doc, err := tfutils.GetDocsForResource("SAP", "btp", "btp", "resources", choice, BtpProviderVersion, "github.com")
	if err != nil {
		log.Fatalf("read doc failed!")
		return "", err
	}

	var importBlock string
	for key, value := range data {

		if slices.Contains(values, strings.Replace(key, ":", "_", -1)) {

			template := strings.Replace(resource_doc.Import, "<resource_name>", strings.Replace(key, ":", "_", -1), -1)
			template = strings.Replace(template, "<subaccount_id>", subaccountId, -1)
			if subMap, ok := value.(map[string]interface{}); ok {
				for subKey, subValue := range subMap {
					template = strings.Replace(template, "<"+subKey+">", fmt.Sprintf("%v", subValue), -1)
				}
			}
			importBlock += template + "\n"
		}
	}

	return importBlock, nil

}

func getEnvInstanceConfig(values []string, subaccountID string, configDir string) {

	dataBlock, err := readDataSource(subaccountID)
	if err != nil {
		fmt.Println("error getting data source:", err)
		return
	}

	dataBlockFile := filepath.Join(TmpFolder, "main.tf")
	err = tfutils.CreateFileWithContent(dataBlockFile, dataBlock)
	if err != nil {
		log.Fatalf("create file %s failed!", dataBlockFile)
		return
	}

	jsonBytes, err := getTfStateData(TmpFolder)
	if err != nil {
		log.Fatalf("error json.Marshal: %s", err)
		return
	}

	jsonString := string(jsonBytes)
	var data map[string]interface{}
	err = json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	importBlock, err := getImportBlock1(data, subaccountID, values)

	if err != nil {
		fmt.Println("error:", err)
		return
	}

	if len(importBlock) == 0 {
		log.Println("No environment instance found for the given subaccount")
		return
	}

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("error getting current directory:", err)
		return
	}
	importFileName := "btp_environment_instances_import.tf"
	importFileName = filepath.Join(currentDir, configDir, importFileName)

	err = tfutils.CreateFileWithContent(importFileName, importBlock)
	if err != nil {
		log.Fatalf("create file %s failed!", dataBlockFile)
		return
	}

	log.Println(" environment instances have been exported. Please check " + configDir + " folder")

}

func getImportBlock1(data map[string]interface{}, subaccountId string, values []string) (string, error) {
	choice := "btp_subaccount_environment_instance"
	resource_doc, err := tfutils.GetDocsForResource("SAP", "btp", "btp", "resources", choice, BtpProviderVersion, "github.com")
	if err != nil {
		log.Fatalf("read doc failed!")
		return "", err
	}

	var importBlock string
	environmentInstances := data["values"].([]interface{})

	for _, value := range environmentInstances {

		environmentInstance := value.(map[string]interface{})
		if slices.Contains(values, fmt.Sprintf("%v", environmentInstance["environment_type"])) {
			template := strings.Replace(resource_doc.Import, "<resource_name>", fmt.Sprintf("%v", environmentInstance["environment_type"]), -1)
			template = strings.Replace(template, "<subaccount_id>", subaccountId, -1)
			template = strings.Replace(template, "<environment_instance_id>", fmt.Sprintf("%v", environmentInstance["id"]), -1)
			importBlock += template + "\n"
		}
	}

	return importBlock, nil
}
