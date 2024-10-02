package tfimportprovider

import (
	"fmt"
	"log"
	"slices"
	"strings"

	output "github.com/SAP/terraform-exporter-btp/pkg/output"
	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
)

type subaccountRoleCollectionImportProvider struct {
	TfImportProvider
}

func newSubaccountRoleCollectionImportProvider() ITfImportProvider {
	return &subaccountRoleCollectionImportProvider{
		TfImportProvider: TfImportProvider{
			resourceType: tfutils.SubaccountRoleCollectionType,
		},
	}
}

func (tf *subaccountRoleCollectionImportProvider) GetImportBlock(data map[string]interface{}, subaccountId string, filterValues []string) (string, error) {
	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.SubaccountRoleCollectionType)
	if err != nil {
		log.Fatalf("read doc failed!")
		return "", err
	}

	importBlock, err := createRoleCollectionImportBlock(data, subaccountId, filterValues, resourceDoc)
	if err != nil {
		return "", err
	}

	return importBlock, nil
}

func createRoleCollectionImportBlock(data map[string]interface{}, subaccountId string, filterValues []string, resourceDoc tfutils.EntityDocs) (importBlock string, err error) {

	roleCollections := data["values"].([]interface{})

	if len(filterValues) != 0 {
		var subaccountAllRoleCollections []string

		for _, value := range roleCollections {
			roleCollection := value.(map[string]interface{})
			resourceName := output.FormatResourceNameGeneric(fmt.Sprintf("%v", roleCollection["name"]))
			subaccountAllRoleCollections = append(subaccountAllRoleCollections, resourceName)
			if slices.Contains(filterValues, resourceName) {
				importBlock += templateRoleCollectionImport(roleCollection, subaccountId, resourceDoc)
			}
		}

		missingRoleCollection, subset := isSubset(subaccountAllRoleCollections, filterValues)

		if !subset {
			return "", fmt.Errorf("role collection %s not found in the subaccount. Please adjust it in the provided file", missingRoleCollection)
		}

	} else {
		for _, value := range roleCollections {
			roleCollection := value.(map[string]interface{})
			importBlock += templateRoleCollectionImport(roleCollection, subaccountId, resourceDoc)
		}
	}

	return importBlock, nil

}

func templateRoleCollectionImport(roleCollection map[string]interface{}, subaccountId string, resourceDoc tfutils.EntityDocs) string {
	resourceDoc.Import = strings.Replace(resourceDoc.Import, "'", "", -1)
	resourceName := output.FormatResourceNameGeneric(fmt.Sprintf("%v", roleCollection["name"]))
	template := strings.Replace(resourceDoc.Import, "<resource_name>", resourceName, -1)
	template = strings.Replace(template, "<subaccount_id>", subaccountId, -1)
	template = strings.Replace(template, "<name>", fmt.Sprintf("%v", roleCollection["name"]), -1)
	return template + "\n"
}
