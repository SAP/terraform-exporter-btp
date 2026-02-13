package tfimportprovider

import (
	"fmt"
	"slices"
	"strings"

	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
)

type directoryEntitlementImportProvider struct {
	TfImportProvider
}

func newDirectoryEntitlementImportProvider() ITfImportProvider {
	return &directoryEntitlementImportProvider{
		TfImportProvider: TfImportProvider{
			resourceType: tfutils.DirectoryEntitlementType,
		},
	}
}

func (tf *directoryEntitlementImportProvider) GetImportBlock(data map[string]any, levelId string, filterValues []string) (string, int, error) {
	count := 0
	directoryId := levelId
	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.DirectoryEntitlementType, tfutils.DirectoryLevel)
	if err != nil {
		return "", count, err
	}

	importBlock, count, err := CreateDirEntitlementImportBlock(data, directoryId, filterValues, resourceDoc)
	if err != nil {
		return "", count, err
	}

	// we only import one directory at a time
	return importBlock, count, nil
}

func CreateDirEntitlementImportBlock(data map[string]any, directoryId string, filterValues []string, resourceDoc tfutils.EntityDocs) (importBlock string, count int, err error) {
	count = 0

	if len(filterValues) != 0 {
		var directoryAllEntitlements []string
		for key, value := range data {
			directoryAllEntitlements = append(directoryAllEntitlements, strings.ReplaceAll(key, ":", "_"))
			if slices.Contains(filterValues, strings.ReplaceAll(key, ":", "_")) {
				importBlock += templateDirEntitlementImport(count, value, directoryId, resourceDoc)
				count++
			}
		}

		missingEntitlement, subset := isSubset(directoryAllEntitlements, filterValues)

		if !subset {
			return "", 0, fmt.Errorf("entitlement %s not found in the directory. Please adjust it in the provided file", missingEntitlement)
		}

	} else {
		for _, value := range data {
			importBlock += templateDirEntitlementImport(count, value, directoryId, resourceDoc)
			count++
		}
	}
	return importBlock, count, nil
}

func templateDirEntitlementImport(x int, value any, directoryId string, resourceDoc tfutils.EntityDocs) string {
	template := strings.ReplaceAll(resourceDoc.Import, "<resource_name>", "entitlement_"+fmt.Sprint(x))
	template = strings.ReplaceAll(template, "<directory_id>", directoryId)
	if subMap, ok := value.(map[string]any); ok {
		for subKey, subValue := range subMap {
			template = strings.ReplaceAll(template, "<"+subKey+">", fmt.Sprintf("%v", subValue))
		}
	}
	return template + "\n"
}
