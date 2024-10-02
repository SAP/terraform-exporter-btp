package tfimportprovider

import (
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

func (tf *subaccountServiceInstanceImportProvider) GetImportBlock(data map[string]interface{}, subaccountId string, filterValues []string) (string, error) {
	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.SubaccountServiceInstanceType)
	if err != nil {
		log.Fatalf("read doc failed!")
		return "", err
	}

	importBlock, err := createServiceInstanceImportBlock(data, subaccountId, filterValues, resourceDoc)
	if err != nil {
		return "", err
	}

	return importBlock, nil
}

func createServiceInstanceImportBlock(data map[string]interface{}, subaccountId string, filterValues []string, resourceDoc tfutils.EntityDocs) (importBlock string, err error) {
	serviceInstances := data["values"].([]interface{})

	if len(filterValues) != 0 {
		var subaccountAllServiceInstances []string

		for _, value := range serviceInstances {
			instance := value.(map[string]interface{})
			resourceName := output.FormatServiceInstanceResourceName(fmt.Sprintf("%v", instance["name"]), fmt.Sprintf("%v", instance["serviceplan_id"]))
			subaccountAllServiceInstances = append(subaccountAllServiceInstances, resourceName)
			if slices.Contains(filterValues, resourceName) {
				importBlock += templateServiceInstanceImport(instance, subaccountId, resourceDoc)
			}
		}

		missingInstance, subset := isSubset(subaccountAllServiceInstances, filterValues)

		if !subset {
			return "", fmt.Errorf("service instance %s not found in the subaccount. Please adjust it in the provided file", missingInstance)
		}

	} else {
		for _, value := range serviceInstances {
			instance := value.(map[string]interface{})
			importBlock += templateServiceInstanceImport(instance, subaccountId, resourceDoc)
		}
	}
	return importBlock, nil
}

func templateServiceInstanceImport(instance map[string]interface{}, subaccountId string, resourceDoc tfutils.EntityDocs) string {
	resourceName := output.FormatServiceInstanceResourceName(fmt.Sprintf("%v", instance["name"]), fmt.Sprintf("%v", instance["serviceplan_id"]))
	template := strings.Replace(resourceDoc.Import, "<resource_name>", resourceName, -1)
	template = strings.Replace(template, "<subaccount_id>", subaccountId, -1)
	template = strings.Replace(template, "<service_instance_id>", fmt.Sprintf("%v", instance["id"]), -1)
	return template + "\n"
}
