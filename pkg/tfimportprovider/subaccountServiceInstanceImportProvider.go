package tfimportprovider

import (
	"encoding/json"
	"fmt"
	"log"
	"slices"
	"strings"

	output "github.com/SAP/terraform-exporter-btp/pkg/output"
	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
)

type subaccountServiceInstanceImportProvider struct {
	TfImportProvider
}

func newSubaccountServiceInstanceImportProvider() ITfImportProvider {
	return &subaccountServiceInstanceImportProvider{
		TfImportProvider: TfImportProvider{
			resourceType: tfutils.SubaccountServiceInstanceType,
		},
	}
}

func (tf *subaccountServiceInstanceImportProvider) GetImportBlock(data map[string]interface{}, levelId string, filterValues []string) (string, int, error) {
	count := 0
	subaccountId := levelId

	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.SubaccountServiceInstanceType, tfutils.SubaccountLevel)
	if err != nil {
		fmt.Print("\r\n")
		log.Fatalf("read doc failed!")
		return "", count, err
	}

	importBlock, count, err := createServiceInstanceImportBlock(data, subaccountId, filterValues, resourceDoc)
	if err != nil {
		return "", count, err
	}

	return importBlock, count, nil
}

func createServiceInstanceImportBlock(data map[string]interface{}, subaccountId string, filterValues []string, resourceDoc tfutils.EntityDocs) (importBlock string, count int, err error) {
	count = 0
	serviceInstances := data["values"].([]interface{})

	if len(filterValues) != 0 {
		var subaccountAllServiceInstances []string

		for x, value := range serviceInstances {
			instance := value.(map[string]interface{})
			resourceName := output.FormatServiceInstanceResourceName(fmt.Sprintf("%v", instance["name"]), fmt.Sprintf("%v", instance["serviceplan_id"]))
			subaccountAllServiceInstances = append(subaccountAllServiceInstances, resourceName)
			if slices.Contains(filterValues, resourceName) {
				importBlock += templateServiceInstanceImport(x, instance, subaccountId, resourceDoc)
				count++
			}
		}

		missingInstance, subset := isSubset(subaccountAllServiceInstances, filterValues)

		if !subset {
			return "", 0, fmt.Errorf("service instance %s not found in the subaccount. Please adjust it in the provided file", missingInstance)
		}

	} else {
		for x, value := range serviceInstances {
			instance := value.(map[string]interface{})
			context := fmt.Sprintf("%v", instance["context"])
			var contextData map[string]interface{}
			if err := json.Unmarshal([]byte(context), &contextData); err != nil {
				return "", 0, fmt.Errorf("error unmarshalling context data: %v", err)
			}
			if contextData["origin"] != "sapcp" {
				continue
			}
			importBlock += templateServiceInstanceImport(x, instance, subaccountId, resourceDoc)
			count++

		}
	}
	return importBlock, count, nil
}

func templateServiceInstanceImport(x int, instance map[string]interface{}, subaccountId string, resourceDoc tfutils.EntityDocs) string {
	template := strings.ReplaceAll(resourceDoc.Import, "<resource_name>", "serviceinstance_"+fmt.Sprint(x))
	template = strings.ReplaceAll(template, "<subaccount_id>", subaccountId)
	template = strings.ReplaceAll(template, "<service_instance_id>", fmt.Sprintf("%v", instance["id"]))
	return template + "\n"
}
