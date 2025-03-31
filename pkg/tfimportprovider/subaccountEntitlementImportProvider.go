package tfimportprovider

import (
	"fmt"
	"slices"
	"strings"

	"github.com/SAP/terraform-exporter-btp/pkg/defaultfilter"
	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
)

type subaccountEntitlementImportProvider struct {
	TfImportProvider
}

func newSubaccountEntitlementImportProvider() ITfImportProvider {
	return &subaccountEntitlementImportProvider{
		TfImportProvider: TfImportProvider{
			resourceType: tfutils.SubaccountEntitlementType,
		},
	}
}

func (tf *subaccountEntitlementImportProvider) GetImportBlock(data map[string]interface{}, levelId string, filterValues []string) (string, int, error) {
	count := 0
	subaccountId := levelId
	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.SubaccountEntitlementType, tfutils.SubaccountLevel)
	if err != nil {
		return "", count, err
	}

	importBlock, count, err := CreateEntitlementImportBlock(data, subaccountId, filterValues, resourceDoc)
	if err != nil {
		return "", count, err
	}

	return importBlock, count, nil
}

func CreateEntitlementImportBlock(data map[string]interface{}, subaccountId string, filterValues []string, resourceDoc tfutils.EntityDocs) (importBlock string, count int, err error) {
	count = 0

	if len(filterValues) != 0 {
		var subaccountAllEntitlements []string
		for key, value := range data {
			subaccountAllEntitlements = append(subaccountAllEntitlements, strings.ReplaceAll(key, ":", "_"))
			if slices.Contains(filterValues, strings.ReplaceAll(key, ":", "_")) {
				importBlock += templateEntitlementImport(count, value, subaccountId, resourceDoc)
				count++
			}
		}

		missingEntitlement, subset := isSubset(subaccountAllEntitlements, filterValues)

		if !subset {
			return "", 0, fmt.Errorf("entitlement %s not found in the subaccount. Please adjust it in the provided file", missingEntitlement)
		}

	} else {
		for _, value := range data {

			serviceName, planName := getServicePlanNameFromData(value)
			if defaultfilter.IsDefaultEntitlement(serviceName, planName) {
				continue
			}

			importBlock += templateEntitlementImport(count, value, subaccountId, resourceDoc)
			count++
		}
	}
	return importBlock, count, nil
}

func templateEntitlementImport(x int, value interface{}, subaccountId string, resourceDoc tfutils.EntityDocs) string {
	template := strings.ReplaceAll(resourceDoc.Import, "<resource_name>", "entitlement_"+fmt.Sprint(x))
	template = strings.ReplaceAll(template, "<subaccount_id>", subaccountId)
	if subMap, ok := value.(map[string]interface{}); ok {
		for subKey, subValue := range subMap {
			template = strings.ReplaceAll(template, "<"+subKey+">", fmt.Sprintf("%v", subValue))
		}
	}
	return template + "\n"
}

func getServicePlanNameFromData(value interface{}) (serviceName string, planName string) {
	if subMap, ok := value.(map[string]interface{}); ok {
		for subKey, subValue := range subMap {
			return subKey, fmt.Sprintf("%v", subValue)
		}
	}
	return "", ""
}
