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

type directoryRoleCollectionImportProvider struct {
	TfImportProvider
}

func newDirectoryRoleCollectionImportProvider() ITfImportProvider {
	return &directoryRoleCollectionImportProvider{
		TfImportProvider: TfImportProvider{
			resourceType: tfutils.DirectoryRoleCollectionType,
		},
	}
}

func (tf *directoryRoleCollectionImportProvider) GetImportBlock(data map[string]interface{}, levelId string, filterValues []string) (string, int, error) {
	count := 0
	directoryId := levelId

	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.DirectoryRoleCollectionType, tfutils.DirectoryLevel)
	if err != nil {
		fmt.Print("\r\n")
		log.Fatalf("read doc failed!")
		return "", count, err
	}

	importBlock, count, err := createDirectoryRoleCollectionImportBlock(data, directoryId, filterValues, resourceDoc)
	if err != nil {
		return "", count, err
	}

	return importBlock, count, nil
}

func createDirectoryRoleCollectionImportBlock(data map[string]interface{}, directoryId string, filterValues []string, resourceDoc tfutils.EntityDocs) (importBlock string, count int, err error) {
	count = 0
	roleCollections := data["values"].([]interface{})

	if len(filterValues) != 0 {
		var directoryAllRoleCollections []string

		for x, value := range roleCollections {
			roleCollection := value.(map[string]interface{})
			resourceName := output.FormatResourceNameGeneric(fmt.Sprintf("%v", roleCollection["name"]))
			directoryAllRoleCollections = append(directoryAllRoleCollections, resourceName)
			if slices.Contains(filterValues, resourceName) {
				importBlock += templateDirectoryRoleCollectionImport(x, roleCollection, directoryId, resourceDoc)
				count++
			}
		}

		missingRoleCollection, subset := isSubset(directoryAllRoleCollections, filterValues)

		if !subset {
			return "", 0, fmt.Errorf("role collection %s not found in the directory. Please adjust it in the provided file", missingRoleCollection)
		}

	} else {
		defaultRoleCollections := defaultfilter.FetchDefaultRoleCollectionsByDirectory(directoryId)

		for x, value := range roleCollections {
			roleCollection := value.(map[string]interface{})

			// exclude default role collections from export
			if defaultfilter.IsRoleCollectionInDefaultList(fmt.Sprintf("%v", roleCollection["name"]), defaultRoleCollections) {
				continue
			}

			importBlock += templateDirectoryRoleCollectionImport(x, roleCollection, directoryId, resourceDoc)
			count++
		}
	}

	return importBlock, count, nil

}

func templateDirectoryRoleCollectionImport(x int, roleCollection map[string]interface{}, directoryId string, resourceDoc tfutils.EntityDocs) string {
	resourceDoc.Import = strings.Replace(resourceDoc.Import, "'", "", -1)
	template := strings.Replace(resourceDoc.Import, "<resource_name>", "rolecollection_"+fmt.Sprint(x), -1)
	template = strings.Replace(template, "<directory_id>", directoryId, -1)
	template = strings.Replace(template, "<name>", fmt.Sprintf("%v", roleCollection["name"]), -1)
	return template + "\n"
}
