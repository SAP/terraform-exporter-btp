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

func (tf *directoryEntitlementImportProvider) GetImportBlock(data map[string]interface{}, levelId string, filterValues []string) (string, error) {
	directoryId := levelId
	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.DirectoryEntitlementType)
	if err != nil {
		return "", err
	}

	importBlock, err := CreateDirEntitlementImportBlock(data, directoryId, filterValues, resourceDoc)
	if err != nil {
		return "", err
	}

	return importBlock, nil
}

func CreateDirEntitlementImportBlock(data map[string]interface{}, directoryId string, filterValues []string, resourceDoc tfutils.EntityDocs) (importBlock string, err error) {

	if len(filterValues) != 0 {
		var directoryAllEntitlements []string
		for key, value := range data {
			directoryAllEntitlements = append(directoryAllEntitlements, strings.Replace(key, ":", "_", -1))
			if slices.Contains(filterValues, strings.Replace(key, ":", "_", -1)) {
				importBlock += templateDirEntitlementImport(key, value, directoryId, resourceDoc)
			}
		}

		missingEntitlement, subset := isSubset(directoryAllEntitlements, filterValues)

		if !subset {
			return "", fmt.Errorf("entitlement %s not found in the directory. Please adjust it in the provided file", missingEntitlement)
		}

	} else {
		for key, value := range data {
			importBlock += templateDirEntitlementImport(key, value, directoryId, resourceDoc)
		}
	}
	return importBlock, nil
}

func templateDirEntitlementImport(key string, value interface{}, directoryId string, resourceDoc tfutils.EntityDocs) string {
	template := strings.Replace(resourceDoc.Import, "<resource_name>", strings.Replace(key, ":", "_", -1), -1)
	template = strings.Replace(template, "<directory_id>", directoryId, -1)
	if subMap, ok := value.(map[string]interface{}); ok {
		for subKey, subValue := range subMap {
			template = strings.Replace(template, "<"+subKey+">", fmt.Sprintf("%v", subValue), -1)
		}
	}
	return template + "\n"
}
