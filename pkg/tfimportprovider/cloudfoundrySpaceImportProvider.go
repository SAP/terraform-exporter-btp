package tfimportprovider

import (
	"fmt"
	"log"
	"slices"
	"strings"

	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
)

type cloudfoundrySpaceImportProvider struct {
	TfImportProvider
}

func newcloudfoundrySpaceImportProvider() ITfImportProvider {
	return &cloudfoundrySpaceImportProvider{
		TfImportProvider: TfImportProvider{
			resourceType: tfutils.CfSpaceType,
		},
	}
}

func (tf *cloudfoundrySpaceImportProvider) GetImportBlock(data map[string]interface{}, levelId string, filterValues []string) (string, int, error) {
	count := 0
	orgId := levelId

	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.CfSpaceType, tfutils.OrganizationLevel)
	if err != nil {
		fmt.Print("\r\n")
		log.Fatalf("read doc failed!")
		return "", count, err
	}

	importBlock, count, err := createSpaceImportBlock(data, orgId, filterValues, resourceDoc)
	if err != nil {
		return "", count, err
	}

	return importBlock, count, nil
}

func createSpaceImportBlock(data map[string]interface{}, orgId string, filterValues []string, resourceDoc tfutils.EntityDocs) (importBlock string, count int, err error) {
	count = 0
	spaces := data["spaces"].([]interface{})

	if len(filterValues) != 0 {
		var cfAllSpaces []string

		for x, value := range spaces {
			space := value.(map[string]interface{})
			cfAllSpaces = append(cfAllSpaces, fmt.Sprintf("%v", space["name"]))
			if slices.Contains(filterValues, fmt.Sprintf("%v", space["name"])) {
				importBlock += templateSpaceImport(x, space, resourceDoc)
				count++
			}
		}

		missingSpace, subset := isSubset(cfAllSpaces, filterValues)

		if !subset {
			return "", 0, fmt.Errorf("cloud foudndry space %s not found in the organization with ID %s. Please adjust it in the provided file", missingSpace, orgId)
		}
	} else {
		for x, value := range spaces {
			space := value.(map[string]interface{})
			importBlock += templateSpaceImport(x, space, resourceDoc)
			count++
		}
	}

	return importBlock, count, nil
}

func templateSpaceImport(x int, space map[string]interface{}, resourceDoc tfutils.EntityDocs) string {
	template := strings.Replace(resourceDoc.Import, "<resource_name>", "space_"+fmt.Sprintf("%v", space["name"]), -1)
	template = strings.Replace(template, "<space_guid>", fmt.Sprintf("%v", space["id"]), -1)
	return template + "\n"
}
