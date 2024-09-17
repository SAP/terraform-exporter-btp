package cmd

import (
	"btptfexport/tfutils"
	"fmt"
	"log"
	"slices"
	"strings"
)

func exportSubaccountSubscriptions(subaccountID string, configDir string, filterValues []string) {

	fmt.Println("")
	spinner, err := startSpinner("crafting import block for " + strings.ToUpper(string(SubaccountSubscriptionType)))
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	data, err := fetchImportConfiguration(subaccountID, SubaccountSubscriptionType, TmpFolder)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	importBlock, err := getSubaccountSubscriptionsImportBlock(data, subaccountID, filterValues)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	if len(importBlock) == 0 {
		log.Println("No subscription found for the given subaccount")
		return
	}

	err = writeImportConfiguration(configDir, SubaccountSubscriptionType, importBlock)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	err = stopSpinner(spinner)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}
}

func getSubaccountSubscriptionsImportBlock(data map[string]interface{}, subaccountId string, filterValues []string) (string, error) {

	resourceDoc, err := getDocByResourceName(ResourcesKind, SubaccountSubscriptionType)
	if err != nil {
		return "", err
	}

	var importBlock string
	subscriptions := data["values"].([]interface{})

	var failedSubscriptions []string
	var inProgressSubscription []string
	if len(filterValues) != 0 {
		var subaccountAllSubscriptions []string

		for _, value := range subscriptions {
			subscription := value.(map[string]interface{})
			subaccountAllSubscriptions = append(subaccountAllSubscriptions, fmt.Sprintf("%v", subscription["app_name"])+"_"+fmt.Sprintf("%v", subscription["plan_name"]))
			if slices.Contains(filterValues, fmt.Sprintf("%v", subscription["app_name"])+"_"+fmt.Sprintf("%v", subscription["plan_name"])) {
				if fmt.Sprintf("%v", subscription["state"]) == "SUBSCRIBED" {
					importBlock += templateSubscriptionImport(subscription, subaccountId, resourceDoc)
				} else if fmt.Sprintf("%v", subscription["state"]) == "SUBSCRIBE_FAILED" {
					failedSubscriptions = append(failedSubscriptions, fmt.Sprintf("%v", subscription["app_name"])+"_"+fmt.Sprintf("%v", subscription["plan_name"]))
				} else if fmt.Sprintf("%v", subscription["state"]) == "IN_PROCESS" {
					inProgressSubscription = append(inProgressSubscription, fmt.Sprintf("%v", subscription["app_name"])+"_"+fmt.Sprintf("%v", subscription["plan_name"]))
				}
			}
		}

		missingSubscription, subset := isSubset(subaccountAllSubscriptions, filterValues)

		if !subset {
			return "", fmt.Errorf("subscription %s not found in the subaccount. Please adjust it in the provided file", missingSubscription)
		}

	} else {
		for _, value := range subscriptions {
			subscription := value.(map[string]interface{})
			if fmt.Sprintf("%v", subscription["state"]) == "SUBSCRIBED" {
				importBlock += templateSubscriptionImport(subscription, subaccountId, resourceDoc)
			} else if fmt.Sprintf("%v", subscription["state"]) == "SUBSCRIBE_FAILED" {
				failedSubscriptions = append(failedSubscriptions, fmt.Sprintf("%v", subscription["app_name"])+"_"+fmt.Sprintf("%v", subscription["plan_name"]))

			} else if fmt.Sprintf("%v", subscription["state"]) == "IN_PROCESS" {
				inProgressSubscription = append(inProgressSubscription, fmt.Sprintf("%v", subscription["app_name"])+"_"+fmt.Sprintf("%v", subscription["plan_name"]))

			}
		}
	}

	if len(failedSubscriptions) != 0 {
		failedSubscriptionsStr := strings.Join(failedSubscriptions, ", ")
		log.Println("Skipping failed subscriptions: " + failedSubscriptionsStr)
	}
	if len(inProgressSubscription) != 0 {
		inProgressSubscriptionStr := strings.Join(inProgressSubscription, ", ")
		log.Println("Skipping in progress subscriptions: " + inProgressSubscriptionStr)
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
