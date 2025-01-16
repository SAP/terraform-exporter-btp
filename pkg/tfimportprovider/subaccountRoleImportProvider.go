package tfimportprovider

import (
	"fmt"
	"log"
	"slices"
	"strings"

	output "github.com/SAP/terraform-exporter-btp/pkg/output"
	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
)

type subaccountRoleImportProvider struct {
	TfImportProvider
}

func newSubaccountRoleImportProvider() ITfImportProvider {
	return &subaccountRoleImportProvider{
		TfImportProvider: TfImportProvider{
			resourceType: tfutils.SubaccountRoleType,
		},
	}
}

func (tf *subaccountRoleImportProvider) GetImportBlock(data map[string]interface{}, levelId string, filterValues []string) (string, int, error) {
	count := 0
	subaccountId := levelId

	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.SubaccountRoleType, tfutils.SubaccountLevel)
	if err != nil {
		fmt.Print("\r\n")
		log.Fatalf("read doc failed!")
		return "", count, err
	}

	importBlock, count, err := createRoleImportBlock(data, subaccountId, filterValues, resourceDoc)
	if err != nil {
		return "", count, err
	}

	return importBlock, count, nil
}

func createRoleImportBlock(data map[string]interface{}, subaccountId string, filterValues []string, resourceDoc tfutils.EntityDocs) (importBlock string, count int, err error) {
	count = 0
	roles := data["values"].([]interface{})

	if len(filterValues) != 0 {
		var subaccountAllRoles []string

		for x, value := range roles {
			role := value.(map[string]interface{})
			resourceName := output.FormatResourceNameGeneric(fmt.Sprintf("%v", role["name"]))
			subaccountAllRoles = append(subaccountAllRoles, resourceName)
			if slices.Contains(filterValues, resourceName) {
				importBlock += templateRoleImport(x, role, subaccountId, resourceDoc)
				count++
			}
		}

		missingRole, subset := isSubset(subaccountAllRoles, filterValues)

		if !subset {
			return "", 0, fmt.Errorf("role %s not found in the subaccount. Please adjust it in the provided file", missingRole)
		}

	} else {
		for x, value := range roles {
			role := value.(map[string]interface{})
			importBlock += templateRoleImport(x, role, subaccountId, resourceDoc)
			count++
		}
	}
	return importBlock, count, nil
}

func templateRoleImport(x int, role map[string]interface{}, subaccountId string, resourceDoc tfutils.EntityDocs) string {

	resourceDoc.Import = strings.Replace(resourceDoc.Import, "'", "", -1)
	template := strings.Replace(resourceDoc.Import, "<resource_name>", "role_"+fmt.Sprint(x), -1)
	template = strings.Replace(template, "<subaccount_id>", subaccountId, -1)
	template = strings.Replace(template, "<name>", fmt.Sprintf("%v", role["name"]), -1)
	template = strings.Replace(template, "<role_template_name>", fmt.Sprintf("%v", role["role_template_name"]), -1)
	template = strings.Replace(template, "<app_id>", fmt.Sprintf("%v", role["app_id"]), -1)
	return template + "\n"
}
