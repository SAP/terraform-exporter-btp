package tfimportprovider

import (
	"fmt"
	"log"
	"slices"
	"strings"

	output "github.com/SAP/terraform-exporter-btp/output"
	tfutils "github.com/SAP/terraform-exporter-btp/tfutils"
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

func (tf *subaccountRoleImportProvider) GetImportBlock(data map[string]interface{}, subaccountId string, filterValues []string) (string, error) {
	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.SubaccountRoleType)
	if err != nil {
		log.Fatalf("read doc failed!")
		return "", err
	}

	importBlock, err := createRoleImportBlock(data, subaccountId, filterValues, resourceDoc)
	if err != nil {
		return "", err
	}

	return importBlock, nil
}

func createRoleImportBlock(data map[string]interface{}, subaccountId string, filterValues []string, resourceDoc tfutils.EntityDocs) (importBlock string, err error) {

	roles := data["values"].([]interface{})

	if len(filterValues) != 0 {
		var subaccountAllRoles []string

		for _, value := range roles {
			role := value.(map[string]interface{})
			resourceName := output.FormatResourceNameGeneric(fmt.Sprintf("%v", role["name"]))
			subaccountAllRoles = append(subaccountAllRoles, resourceName)
			if slices.Contains(filterValues, resourceName) {
				importBlock += templateRoleImport(role, subaccountId, resourceDoc)
			}
		}

		missingRole, subset := isSubset(subaccountAllRoles, filterValues)

		if !subset {
			return "", fmt.Errorf("role %s not found in the subaccount. Please adjust it in the provided file", missingRole)
		}

	} else {
		for _, value := range roles {
			role := value.(map[string]interface{})
			importBlock += templateRoleImport(role, subaccountId, resourceDoc)
		}
	}
	return importBlock, nil
}

func templateRoleImport(role map[string]interface{}, subaccountId string, resourceDoc tfutils.EntityDocs) string {

	resourceDoc.Import = strings.Replace(resourceDoc.Import, "'", "", -1)
	resourceName := output.FormatResourceNameGeneric(fmt.Sprintf("%v", role["name"]))
	template := strings.Replace(resourceDoc.Import, "<resource_name>", resourceName, -1)
	template = strings.Replace(template, "<subaccount_id>", subaccountId, -1)
	template = strings.Replace(template, "<name>", fmt.Sprintf("%v", role["name"]), -1)
	template = strings.Replace(template, "<role_template_name>", fmt.Sprintf("%v", role["role_template_name"]), -1)
	template = strings.Replace(template, "<app_id>", fmt.Sprintf("%v", role["app_id"]), -1)
	return template + "\n"
}
