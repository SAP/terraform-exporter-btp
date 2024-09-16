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

func exportSubaccountSubscriptions(subaccountID string, configDir string, filterValues []string) {

	dataBlock, err := readDataSource(subaccountID, SubaccountSubscriptionType)
	if err != nil {
		return
	}

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("error getting current directory: %v", err)
		return
	}
	dataBlockFile := filepath.Join(TmpFolder, "main.tf")
	err = tfutils.CreateFileWithContent(dataBlockFile, dataBlock)
	if err != nil {
		log.Fatalf("create file %s failed!", dataBlockFile)
		return
	}

	jsonBytes, err := GetTfStateData(TmpFolder, SubaccountSubscriptionType)
	if err != nil {
		log.Fatalf("error json.Marshal: %v", err)
		return
	}

	var data map[string]interface{}
	jsonString := string(jsonBytes)

	err = json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	importBlock, err := getSubscriptionsImportBlock(data, subaccountID, filterValues)
	if err != nil {
		fmt.Println("error:", err)
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
		log.Fatalf("create file %v failed!", dataBlockFile)
		return
	}

	log.Println("subaccount subscriptions has been exported. Please check " + configDir + " folder")
}

func getSubscriptionsImportBlock(data map[string]interface{}, subaccountId string, filterValues []string) (string, error) {
	resourceDoc, err := getDocByResourceName(ResourcesKind, SubaccountSubscriptionType)
	if err != nil {
		return "", err
	}

	var importBlock string
	subscriptions := data["values"].([]interface{})

	if len(filterValues) != 0 {
		var subaccountAllSubscriptions []string

		for _, value := range subscriptions {
			subscription := value.(map[string]interface{})
			subaccountAllSubscriptions = append(subaccountAllSubscriptions, fmt.Sprintf("%v", subscription["app_name"])+"_"+fmt.Sprintf("%v", subscription["plan_name"]))
			if slices.Contains(filterValues, fmt.Sprintf("%v", subscription["app_name"])+"_"+fmt.Sprintf("%v", subscription["plan_name"])) {
				importBlock += templateSubscriptionImport(subscription, subaccountId, resourceDoc)
			}
		}

		missingSubscription, subset := isSubset(subaccountAllSubscriptions, filterValues)

		if !subset {
			return "", fmt.Errorf("subscription %s not found in the subaccount. Please adjust it in the provided file", missingSubscription)

		}

	} else {
		for _, value := range subscriptions {
			subscription := value.(map[string]interface{})
			if fmt.Sprintf("%v", subscription["state"]) != "NOT_SUBSCRIBED" {
				importBlock += templateSubscriptionImport(subscription, subaccountId, resourceDoc)
			}
		}
	}

	return importBlock, nil
}

func templateSubscriptionImport(subscription map[string]interface{}, subaccountId string, resourceDoc tfutils.EntityDocs) string {
	template := strings.Replace(resourceDoc.Import, "<resource_name>", strings.Replace(fmt.Sprintf("%v", subscription["app_name"]), "-", "_", -1), -1)
	template = strings.Replace(template, "<subaccount_id>", subaccountId, -1)
	template = strings.Replace(template, "<app_name>", fmt.Sprintf("%v", subscription["app_name"]), -1)
	template = strings.Replace(template, "<plan_name>", fmt.Sprintf("%v", subscription["plan_name"]), -1)
	return template + "\n"
}
