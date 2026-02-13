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

func (tf *directoryRoleImportProvider) GetImportBlock(data map[string]any, levelId string, filterValues []string) (string, int, error) {
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

func createDirectoryRoleImportBlock(data map[string]any, directoryId string, filterValues []string, resourceDoc tfutils.EntityDocs) (importBlock string, count int, err error) {
	count = 0
	roles := data["values"].([]any)

	if len(filterValues) != 0 {
		var directoryAllRoles []string

		for x, value := range roles {
			role := value.(map[string]any)
			resourceName := output.FormatResourceNameGeneric(fmt.Sprintf("%v", role["name"]))
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
			role := value.(map[string]any)

			// exclude default roles from export
			if defaultfilter.IsRoleInDefaultList(fmt.Sprintf("%v", role["name"]), defaultfilter.FetchDefaultRolesByDirectory(directoryId)) {
				continue
			}

			importBlock += templateDirectoryRoleImport(x, role, directoryId, resourceDoc)
			count++
		}
	}
	return importBlock, count, nil
}

func templateDirectoryRoleImport(x int, role map[string]any, directoryId string, resourceDoc tfutils.EntityDocs) string {

	resourceDoc.Import = strings.ReplaceAll(resourceDoc.Import, "'", "")
	template := strings.ReplaceAll(resourceDoc.Import, "<resource_name>", "role_"+fmt.Sprint(x))
	template = strings.ReplaceAll(template, "<directory_id>", directoryId)
	template = strings.ReplaceAll(template, "<name>", fmt.Sprintf("%v", role["name"]))
	template = strings.ReplaceAll(template, "<role_template_name>", fmt.Sprintf("%v", role["role_template_name"]))
	template = strings.ReplaceAll(template, "<app_id>", fmt.Sprintf("%v", role["app_id"]))
	return template + "\n"
}
