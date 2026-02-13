package tfimportprovider

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/SAP/terraform-exporter-btp/pkg/defaultfilter"
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

func (tf *subaccountRoleCollectionImportProvider) GetImportBlock(data map[string]any, levelId string, filterValues []string) (string, int, error) {
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

func createRoleCollectionImportBlock(data map[string]any, subaccountId string, filterValues []string, resourceDoc tfutils.EntityDocs) (importBlock string, count int, err error) {
	count = 0
	roleCollections := data["values"].([]any)

	if len(filterValues) != 0 {
		var subaccountAllRoleCollections []string

		for x, value := range roleCollections {
			roleCollection := value.(map[string]any)

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
		defaultRoleCollections := defaultfilter.FetchDefaultRoleCollectionsBySubaccount(subaccountId)

		for x, value := range roleCollections {
			roleCollection := value.(map[string]any)

			// exclude default role collections from export
			if defaultfilter.IsRoleCollectionInDefaultList(fmt.Sprintf("%v", roleCollection["name"]), defaultRoleCollections) {
				continue
			}

			importBlock += templateRoleCollectionImport(x, roleCollection, subaccountId, resourceDoc)
			count++
		}
	}

	return importBlock, count, nil

}

func templateRoleCollectionImport(x int, roleCollection map[string]any, subaccountId string, resourceDoc tfutils.EntityDocs) string {
	resourceDoc.Import = strings.ReplaceAll(resourceDoc.Import, "'", "")
	template := strings.ReplaceAll(resourceDoc.Import, "<resource_name>", "rolecollection_"+fmt.Sprint(x))
	template = strings.ReplaceAll(template, "<subaccount_id>", subaccountId)
	template = strings.ReplaceAll(template, "<name>", fmt.Sprintf("%v", roleCollection["name"]))
	return template + "\n"
}
