package tfimportprovider

import (
	"fmt"
	"log"
	"slices"
	"strings"

	output "github.com/SAP/terraform-exporter-btp/pkg/output"
	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
)

type directoryRoleImportProvider struct {
	TfImportProvider
}

func newDirectoryRoleImportProvider() ITfImportProvider {
	return &directoryRoleImportProvider{
		TfImportProvider: TfImportProvider{
			resourceType: tfutils.DirectoryRoleType,
		},
	}
}

func (tf *directoryRoleImportProvider) GetImportBlock(data map[string]interface{}, levelId string, filterValues []string) (string, error) {

	directoryId := levelId

	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.DirectoryRoleType)
	if err != nil {
		log.Fatalf("read doc failed!")
		return "", err
	}

	importBlock, err := createDirectoryRoleImportBlock(data, directoryId, filterValues, resourceDoc)
	if err != nil {
		return "", err
	}

	return importBlock, nil
}

func createDirectoryRoleImportBlock(data map[string]interface{}, directoryId string, filterValues []string, resourceDoc tfutils.EntityDocs) (importBlock string, err error) {

	roles := data["values"].([]interface{})

	if len(filterValues) != 0 {
		var directoryAllRoles []string

		for _, value := range roles {
			role := value.(map[string]interface{})
			resourceName := output.FormatResourceNameGeneric(fmt.Sprintf("%v", role["name"]))
			directoryAllRoles = append(directoryAllRoles, resourceName)
			if slices.Contains(filterValues, resourceName) {
				importBlock += templateDirectoryRoleImport(role, directoryId, resourceDoc)
			}
		}

		missingRole, subset := isSubset(directoryAllRoles, filterValues)

		if !subset {
			return "", fmt.Errorf("role %s not found in the directory. Please adjust it in the provided file", missingRole)
		}

	} else {
		for _, value := range roles {
			role := value.(map[string]interface{})
			importBlock += templateDirectoryRoleImport(role, directoryId, resourceDoc)
		}
	}
	return importBlock, nil
}

func templateDirectoryRoleImport(role map[string]interface{}, directoryId string, resourceDoc tfutils.EntityDocs) string {

	resourceDoc.Import = strings.Replace(resourceDoc.Import, "'", "", -1)
	resourceName := output.FormatResourceNameGeneric(fmt.Sprintf("%v", role["name"]))
	template := strings.Replace(resourceDoc.Import, "<resource_name>", resourceName, -1)
	template = strings.Replace(template, "<directory_id>", directoryId, -1)
	template = strings.Replace(template, "<name>", fmt.Sprintf("%v", role["name"]), -1)
	template = strings.Replace(template, "<role_template_name>", fmt.Sprintf("%v", role["role_template_name"]), -1)
	template = strings.Replace(template, "<app_id>", fmt.Sprintf("%v", role["app_id"]), -1)
	return template + "\n"
}