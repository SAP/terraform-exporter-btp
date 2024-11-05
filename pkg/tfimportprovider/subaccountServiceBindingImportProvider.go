package tfimportprovider

import (
	"fmt"
	"log"
	"slices"
	"strings"

	output "github.com/SAP/terraform-exporter-btp/pkg/output"
	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
)

type subaccountServiceBindingImportProvider struct {
	TfImportProvider
}

func newSubaccountServiceBindingImportProvider() ITfImportProvider {
	return &subaccountServiceBindingImportProvider{
		TfImportProvider: TfImportProvider{
			resourceType: tfutils.SubaccountServiceBindingType,
		},
	}
}

func (tf *subaccountServiceBindingImportProvider) GetImportBlock(data map[string]interface{}, levelId string, filterValues []string) (string, int, error) {
	count := 0
	subaccountId := levelId

	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.SubaccountServiceBindingType, tfutils.SubaccountLevel)
	if err != nil {
		fmt.Print("\r\n")
		log.Fatalf("read doc failed!")
		return "", count, err
	}

	importBlock, count, err := createServiceBindingImportBlock(data, subaccountId, filterValues, resourceDoc)
	if err != nil {
		return "", count, err
	}

	return importBlock, count, nil
}

func createServiceBindingImportBlock(data map[string]interface{}, subaccountId string, filterValues []string, resourceDoc tfutils.EntityDocs) (importBlock string, count int, err error) {
	count = 0
	serviceBindings := data["values"].([]interface{})

	if len(filterValues) != 0 {
		var subaccountAllServiceBindings []string

		for _, value := range serviceBindings {
			binding := value.(map[string]interface{})
			resourceName := fmt.Sprintf("%v", binding["name"])
			subaccountAllServiceBindings = append(subaccountAllServiceBindings, resourceName)
			if slices.Contains(filterValues, resourceName) {
				importBlock += templateServiceBindingImport(binding, subaccountId, resourceDoc)
				count++
			}
		}

		missingBinding, subset := isSubset(subaccountAllServiceBindings, filterValues)

		if !subset {
			return "", 0, fmt.Errorf("service binding %s not found in the subaccount. Please adjust it in the provided file", missingBinding)
		}

	} else {
		for _, value := range serviceBindings {
			binding := value.(map[string]interface{})
			importBlock += templateServiceBindingImport(binding, subaccountId, resourceDoc)
			count++
		}
	}
	return importBlock, count, nil
}

func templateServiceBindingImport(binding map[string]interface{}, subaccountId string, resourceDoc tfutils.EntityDocs) string {
	resourceName := output.FormatResourceNameGeneric(fmt.Sprintf("%v", binding["name"]))
	template := strings.Replace(resourceDoc.Import, "<resource_name>", resourceName, -1)
	template = strings.Replace(template, "<subaccount_id>", subaccountId, -1)
	template = strings.Replace(template, "<service_binding_id>", fmt.Sprintf("%v", binding["id"]), -1)
	return template + "\n"
}
