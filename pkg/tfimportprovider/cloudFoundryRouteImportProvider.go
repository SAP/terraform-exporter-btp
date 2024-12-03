package tfimportprovider

import (
	"fmt"
	"log"
	"slices"
	"strings"

	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
)

type cloudfoundryRouteImportProvider struct {
	TfImportProvider
}

func newCloudfoundryRouteImportProvider() ITfImportProvider {
	return &cloudfoundryRouteImportProvider{
		TfImportProvider: TfImportProvider{
			resourceType: tfutils.CfRouteType,
		},
	}
}
func (tf *cloudfoundryRouteImportProvider) GetImportBlock(data map[string]interface{}, levelId string, filterValues []string) (string, int, error) {
	count := 0
	orgId := levelId
	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.CfRouteType, tfutils.OrganizationLevel)
	if err != nil {
		fmt.Print("\r\n")
		log.Fatalf("read doc failed!")
		return "", count, err
	}
	importBlock, count, err := createRouteImportBlock(data, orgId, filterValues, resourceDoc)
	if err != nil {
		return "", count, err
	}
	return importBlock, count, nil
}
func createRouteImportBlock(data map[string]interface{}, orgId string, filterValues []string, resourceDoc tfutils.EntityDocs) (importBlock string, count int, err error) {
	count = 0
	routes := data["routes"].([]interface{})
	if len(filterValues) != 0 {
		var cfAllRoutes []string
		for x, value := range routes {
			route := value.(map[string]interface{})
			cfAllRoutes = append(cfAllRoutes, fmt.Sprintf("%v", route["url"]))
			if slices.Contains(filterValues, fmt.Sprintf("%v", route["url"])) {
				importBlock += templateRouteImport(x, route, resourceDoc)
				count++
			}
		}
		missingRoute, subset := isSubset(cfAllRoutes, filterValues)
		if !subset {
			return "", 0, fmt.Errorf("cloud foundry route %s not found in the organization with ID %s. Please adjust it in the provided file", missingRoute, orgId)
		}
	} else {
		for x, value := range routes {
			route := value.(map[string]interface{})
			importBlock += templateRouteImport(x, route, resourceDoc)
			count++
		}
	}
	return importBlock, count, nil
}
func templateRouteImport(x int, route map[string]interface{}, resourceDoc tfutils.EntityDocs) string {
	template := strings.Replace(resourceDoc.Import, "<resource_name>", "route_"+fmt.Sprintf("%v", x), -1)
	template = strings.Replace(template, "<route_guid>", fmt.Sprintf("%v", route["id"]), -1)
	return template + "\n"
}
