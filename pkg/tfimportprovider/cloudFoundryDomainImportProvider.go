package tfimportprovider

import (
	"fmt"
	"log"
	"slices"
	"strings"

	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
)

type cloudfoundryDomainImportProvider struct {
	TfImportProvider
}

func newCloudfoundryDomainImportProvider() ITfImportProvider {
	return &cloudfoundryDomainImportProvider{
		TfImportProvider: TfImportProvider{
			resourceType: tfutils.CfDomainType,
		},
	}
}
func (tf *cloudfoundryDomainImportProvider) GetImportBlock(data map[string]interface{}, levelId string, filterValues []string) (string, int, error) {
	count := 0
	orgId := levelId
	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.CfDomainType, tfutils.OrganizationLevel)
	if err != nil {
		fmt.Print("\r\n")
		log.Fatalf("read doc failed!")
		return "", count, err
	}
	importBlock, count, err := createDomainImportBlock(data, orgId, filterValues, resourceDoc)
	if err != nil {
		return "", count, err
	}
	return importBlock, count, nil
}
func createDomainImportBlock(data map[string]interface{}, orgId string, filterValues []string, resourceDoc tfutils.EntityDocs) (importBlock string, count int, err error) {
	count = 0
	domains := data["domains"].([]interface{})
	if len(filterValues) != 0 {
		var cfAllDomains []string
		for x, value := range domains {
			domain := value.(map[string]interface{})
			cfAllDomains = append(cfAllDomains, fmt.Sprintf("%v", domain["name"]))
			if slices.Contains(filterValues, fmt.Sprintf("%v", domain["name"])) {
				importBlock += templateDomainImport(x, domain, resourceDoc)
				count++
			}
		}
		missingSpace, subset := isSubset(cfAllDomains, filterValues)
		if !subset {
			return "", 0, fmt.Errorf("cloud foundry domain %s not found in the organization with ID %s. Please adjust it in the provided file", missingSpace, orgId)
		}
	} else {
		for x, value := range domains {
			domain := value.(map[string]interface{})
			importBlock += templateDomainImport(x, domain, resourceDoc)
			count++
		}
	}
	return importBlock, count, nil
}
func templateDomainImport(x int, domain map[string]interface{}, resourceDoc tfutils.EntityDocs) string {
	template := strings.Replace(resourceDoc.Import, "<resource_name>", "domain_"+fmt.Sprintf("%v", x), -1)
	template = strings.Replace(template, "<domain_guid>", fmt.Sprintf("%v", domain["id"]), -1)
	return template + "\n"
}
