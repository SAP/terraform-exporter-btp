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

func (tf *subaccountRoleCollectionImportProvider) GetImportBlock(data map[string]interface{}, levelId string, filterValues []string) (string, int, error) {
	count := 0
	subaccountId := levelId

	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.SubaccountRoleCollectionType, tfutils.SubaccountLevel)
	if err != nil {
		fmt.Print("\r\n")
		log.Fatalf("read doc failed!")
		return "", count, err
	}

	importBlock, count, err := createRoleCollectionImportBlock(data, subaccountId, filterValues, resourceDoc)
	if err != nil {
		return "", count, err
	}

	return importBlock, count, nil
}

func createRoleCollectionImportBlock(data map[string]interface{}, subaccountId string, filterValues []string, resourceDoc tfutils.EntityDocs) (importBlock string, count int, err error) {
	count = 0
	roleCollections := data["values"].([]interface{})

	if len(filterValues) != 0 {
		var subaccountAllRoleCollections []string

		for x, value := range roleCollections {
			roleCollection := value.(map[string]interface{})
			resourceName := output.FormatResourceNameGeneric(fmt.Sprintf("%v", roleCollection["name"]))
			subaccountAllRoleCollections = append(subaccountAllRoleCollections, resourceName)
			if slices.Contains(filterValues, resourceName) {
				importBlock += templateRoleCollectionImport(x, roleCollection, subaccountId, resourceDoc)
				count++
			}
		}

		missingRoleCollection, subset := isSubset(subaccountAllRoleCollections, filterValues)

		if !subset {
			return "", 0, fmt.Errorf("role collection %s not found in the subaccount. Please adjust it in the provided file", missingRoleCollection)
		}

	} else {
		for x, value := range roleCollections {
			roleCollection := value.(map[string]interface{})
			importBlock += templateRoleCollectionImport(x, roleCollection, subaccountId, resourceDoc)
			count++
		}
	}

	return importBlock, count, nil

}

func templateRoleCollectionImport(x int, roleCollection map[string]interface{}, subaccountId string, resourceDoc tfutils.EntityDocs) string {
	resourceDoc.Import = strings.Replace(resourceDoc.Import, "'", "", -1)
	template := strings.Replace(resourceDoc.Import, "<resource_name>", "rolecollection_"+fmt.Sprint(x), -1)
	template = strings.Replace(template, "<subaccount_id>", subaccountId, -1)
	template = strings.Replace(template, "<name>", fmt.Sprintf("%v", roleCollection["name"]), -1)
	return template + "\n"
}
