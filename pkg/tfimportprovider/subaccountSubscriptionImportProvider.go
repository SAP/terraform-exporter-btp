package tfimportprovider

import (
	"fmt"
	"log"
	"slices"
	"strings"

	output "github.com/SAP/terraform-exporter-btp/pkg/output"
	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
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

func (tf *subaccountSubscriptionImportProvider) GetImportBlock(data map[string]interface{}, levelId string, filterValues []string) (string, int, error) {
	count := 0
	subaccountId := levelId

	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.SubaccountSubscriptionType, tfutils.SubaccountLevel)
	if err != nil {
		return "", count, err
	}

	importBlock, count, err := createSubscriptionImportBlock(data, subaccountId, filterValues, resourceDoc)
	if err != nil {
		return "", count, err
	}

	return importBlock, count, nil
}

func createSubscriptionImportBlock(data map[string]interface{}, subaccountId string, filterValues []string, resourceDoc tfutils.EntityDocs) (importBlock string, count int, err error) {
	count = 0
	subscriptions := data["values"].([]interface{})

	var failedSubscriptions []string
	var inProgressSubscription []string
	if len(filterValues) != 0 {
		var subaccountAllSubscriptions []string

		for x, value := range subscriptions {
			subscription := value.(map[string]interface{})
			subaccountAllSubscriptions = append(subaccountAllSubscriptions, output.FormatSubscriptionResourceName(fmt.Sprintf("%v", subscription["app_name"]), fmt.Sprintf("%v", subscription["plan_name"])))
			if slices.Contains(filterValues, output.FormatSubscriptionResourceName(fmt.Sprintf("%v", subscription["app_name"]), fmt.Sprintf("%v", subscription["plan_name"]))) {
				if fmt.Sprintf("%v", subscription["state"]) == "SUBSCRIBED" {
					importBlock += templateSubscriptionImport(x, subscription, subaccountId, resourceDoc)
					count++
				} else if fmt.Sprintf("%v", subscription["state"]) == "SUBSCRIBE_FAILED" {
					failedSubscriptions = append(failedSubscriptions, output.FormatSubscriptionResourceName(fmt.Sprintf("%v", subscription["app_name"]), fmt.Sprintf("%v", subscription["plan_name"])))
				} else if fmt.Sprintf("%v", subscription["state"]) == "IN_PROCESS" {
					inProgressSubscription = append(inProgressSubscription, output.FormatSubscriptionResourceName(fmt.Sprintf("%v", subscription["app_name"]), fmt.Sprintf("%v", subscription["plan_name"])))
				}
			}
		}

		missingSubscription, subset := isSubset(subaccountAllSubscriptions, filterValues)

		if !subset {
			return "", 0, fmt.Errorf("subscription %s not found in the subaccount. Please adjust it in the provided file", missingSubscription)
		}

	} else {
		for x, value := range subscriptions {
			subscription := value.(map[string]interface{})
			if fmt.Sprintf("%v", subscription["state"]) == "SUBSCRIBED" {
				importBlock += templateSubscriptionImport(x, subscription, subaccountId, resourceDoc)
				count++
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
	return importBlock, count, nil
}

func templateSubscriptionImport(x int, subscription map[string]interface{}, subaccountId string, resourceDoc tfutils.EntityDocs) string {
	template := strings.ReplaceAll(resourceDoc.Import, "<resource_name>", "subscription_"+fmt.Sprint(x))
	template = strings.ReplaceAll(template, "<subaccount_id>", subaccountId)
	template = strings.ReplaceAll(template, "<app_name>", fmt.Sprintf("%v", subscription["app_name"]))
	template = strings.ReplaceAll(template, "<plan_name>", fmt.Sprintf("%v", subscription["plan_name"]))
	return template + "\n"
}
