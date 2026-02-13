package tfimportprovider

import (
	"fmt"
	"log"
	"slices"
	"strings"

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

func (tf *subaccountServiceBindingImportProvider) GetImportBlock(data map[string]any, levelId string, filterValues []string) (string, int, error) {
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

func createServiceBindingImportBlock(data map[string]any, subaccountId string, filterValues []string, resourceDoc tfutils.EntityDocs) (importBlock string, count int, err error) {
	count = 0
	serviceBindings := data["values"].([]any)

	if len(filterValues) != 0 {
		var subaccountAllServiceBindings []string

		for x, value := range serviceBindings {
			binding := value.(map[string]any)
			resourceName := fmt.Sprintf("%v", binding["name"])
			subaccountAllServiceBindings = append(subaccountAllServiceBindings, resourceName)
			if slices.Contains(filterValues, resourceName) {
				importBlock += templateServiceBindingImport(x, binding, subaccountId, resourceDoc)
				count++
			}
		}

		missingBinding, subset := isSubset(subaccountAllServiceBindings, filterValues)

		if !subset {
			return "", 0, fmt.Errorf("service binding %s not found in the subaccount. Please adjust it in the provided file", missingBinding)
		}

	} else {
		for x, value := range serviceBindings {
			binding := value.(map[string]any)
			importBlock += templateServiceBindingImport(x, binding, subaccountId, resourceDoc)
			count++
		}
	}
	return importBlock, count, nil
}

func templateServiceBindingImport(x int, binding map[string]any, subaccountId string, resourceDoc tfutils.EntityDocs) string {
	template := strings.ReplaceAll(resourceDoc.Import, "<resource_name>", "servicebinding_"+fmt.Sprint(x))
	template = strings.ReplaceAll(template, "<subaccount_id>", subaccountId)
	template = strings.ReplaceAll(template, "<service_binding_id>", fmt.Sprintf("%v", binding["id"]))
	return template + "\n"
}
