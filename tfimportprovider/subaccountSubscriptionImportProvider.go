package tfimportprovider

import (
	"fmt"
	"log"
	"slices"
	"strings"

	output "github.com/SAP/terraform-exporter-btp/output"
	tfutils "github.com/SAP/terraform-exporter-btp/tfutils"
)

type subaccountSubscriptionImportProvider struct {
	TfImportProvider
}

func newSubaccountSubscriptionImportProvider() ITfImportProvider {
	return &subaccountSubscriptionImportProvider{
		TfImportProvider: TfImportProvider{
			resourceType: tfutils.SubaccountSubscriptionType,
		},
	}
}

func (tf *subaccountSubscriptionImportProvider) GetImportBlock(data map[string]interface{}, subaccountId string, filterValues []string) (string, error) {

	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.SubaccountSubscriptionType)
	if err != nil {
		return "", err
	}

	importBlock, err := createSubscriptionImportBlock(data, subaccountId, filterValues, resourceDoc)
	if err != nil {
		return "", err
	}

	return importBlock, nil
}

func createSubscriptionImportBlock(data map[string]interface{}, subaccountId string, filterValues []string, resourceDoc tfutils.EntityDocs) (importBlock string, err error) {
	subscriptions := data["values"].([]interface{})

	var failedSubscriptions []string
	var inProgressSubscription []string
	if len(filterValues) != 0 {
		var subaccountAllSubscriptions []string

		for _, value := range subscriptions {
			subscription := value.(map[string]interface{})
			subaccountAllSubscriptions = append(subaccountAllSubscriptions, output.FormatSubscriptionResourceName(fmt.Sprintf("%v", subscription["app_name"]), fmt.Sprintf("%v", subscription["plan_name"])))
			if slices.Contains(filterValues, output.FormatSubscriptionResourceName(fmt.Sprintf("%v", subscription["app_name"]), fmt.Sprintf("%v", subscription["plan_name"]))) {
				if fmt.Sprintf("%v", subscription["state"]) == "SUBSCRIBED" {
					importBlock += templateSubscriptionImport(subscription, subaccountId, resourceDoc)
				} else if fmt.Sprintf("%v", subscription["state"]) == "SUBSCRIBE_FAILED" {
					failedSubscriptions = append(failedSubscriptions, output.FormatSubscriptionResourceName(fmt.Sprintf("%v", subscription["app_name"]), fmt.Sprintf("%v", subscription["plan_name"])))
				} else if fmt.Sprintf("%v", subscription["state"]) == "IN_PROCESS" {
					inProgressSubscription = append(inProgressSubscription, output.FormatSubscriptionResourceName(fmt.Sprintf("%v", subscription["app_name"]), fmt.Sprintf("%v", subscription["plan_name"])))
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
				failedSubscriptions = append(failedSubscriptions, output.FormatSubscriptionResourceName(fmt.Sprintf("%v", subscription["app_name"]), fmt.Sprintf("%v", subscription["plan_name"])))
			} else if fmt.Sprintf("%v", subscription["state"]) == "IN_PROCESS" {
				inProgressSubscription = append(inProgressSubscription, output.FormatSubscriptionResourceName(fmt.Sprintf("%v", subscription["app_name"]), fmt.Sprintf("%v", subscription["plan_name"])))
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
