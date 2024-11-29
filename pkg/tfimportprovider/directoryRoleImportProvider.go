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

func (tf *directoryRoleImportProvider) GetImportBlock(data map[string]interface{}, levelId string, filterValues []string) (string, int, error) {
	count := 0
	directoryId := levelId

	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.DirectoryRoleType, tfutils.DirectoryLevel)
	if err != nil {
		fmt.Print("\r\n")
		log.Fatalf("read doc failed!")
		return "", count, err
	}

	importBlock, count, err := createDirectoryRoleImportBlock(data, directoryId, filterValues, resourceDoc)
	if err != nil {
		return "", count, err
	}

	return importBlock, count, nil
}

func createDirectoryRoleImportBlock(data map[string]interface{}, directoryId string, filterValues []string, resourceDoc tfutils.EntityDocs) (importBlock string, count int, err error) {
	count = 0
	roles := data["values"].([]interface{})

	if len(filterValues) != 0 {
		var directoryAllRoles []string

		for x, value := range roles {
			role := value.(map[string]interface{})
			resourceName := output.FormatRoleResourceName(fmt.Sprintf("%v", role["name"]))
			directoryAllRoles = append(directoryAllRoles, resourceName)
			if slices.Contains(filterValues, resourceName) {
				importBlock += templateDirectoryRoleImport(x, role, directoryId, resourceDoc)
				count++
			}
		}

		missingRole, subset := isSubset(directoryAllRoles, filterValues)

		if !subset {
			return "", 0, fmt.Errorf("role %s not found in the directory. Please adjust it in the provided file", missingRole)
		}

	} else {
		for x, value := range roles {
			role := value.(map[string]interface{})
			importBlock += templateDirectoryRoleImport(x, role, directoryId, resourceDoc)
			count++
		}
	}
	return importBlock, count, nil
}

func templateDirectoryRoleImport(x int, role map[string]interface{}, directoryId string, resourceDoc tfutils.EntityDocs) string {

	resourceDoc.Import = strings.Replace(resourceDoc.Import, "'", "", -1)
	template := strings.Replace(resourceDoc.Import, "<resource_name>", "directory_role_"+fmt.Sprint(x), -1)
	template = strings.Replace(template, "<directory_id>", directoryId, -1)
	template = strings.Replace(template, "<name>", fmt.Sprintf("%v", role["name"]), -1)
	template = strings.Replace(template, "<role_template_name>", fmt.Sprintf("%v", role["role_template_name"]), -1)
	template = strings.Replace(template, "<app_id>", fmt.Sprintf("%v", role["app_id"]), -1)
	return template + "\n"
}
