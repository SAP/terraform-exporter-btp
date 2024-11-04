package tfimportprovider

import (
	"fmt"
	"slices"
	"strings"

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
	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.SubaccountEntitlementType, levelId)
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
			subaccountAllEntitlements = append(subaccountAllEntitlements, strings.Replace(key, ":", "_", -1))
			if slices.Contains(filterValues, strings.Replace(key, ":", "_", -1)) {
				importBlock += templateEntitlementImport(key, value, subaccountId, resourceDoc)
				count++
			}
		}

		missingEntitlement, subset := isSubset(subaccountAllEntitlements, filterValues)

		if !subset {
			return "", 0, fmt.Errorf("entitlement %s not found in the subaccount. Please adjust it in the provided file", missingEntitlement)
		}

	} else {
		for key, value := range data {
			importBlock += templateEntitlementImport(key, value, subaccountId, resourceDoc)
			count++
		}
	}
	return importBlock, count, nil
}

func templateEntitlementImport(key string, value interface{}, subaccountId string, resourceDoc tfutils.EntityDocs) string {
	template := strings.Replace(resourceDoc.Import, "<resource_name>", strings.Replace(key, ":", "_", -1), -1)
	template = strings.Replace(template, "<subaccount_id>", subaccountId, -1)
	if subMap, ok := value.(map[string]interface{}); ok {
		for subKey, subValue := range subMap {
			template = strings.Replace(template, "<"+subKey+">", fmt.Sprintf("%v", subValue), -1)
		}
	}
	return template + "\n"
}
