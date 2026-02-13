package tfimportprovider

import (
	"fmt"
	"log"
	"slices"
	"strings"

	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
)

type cloudfoundrySpaceQuotaImportProvider struct {
	TfImportProvider
}

func newcloudfoundrySpaceQuotaImportProvider() ITfImportProvider {
	return &cloudfoundrySpaceQuotaImportProvider{
		TfImportProvider: TfImportProvider{
			resourceType: tfutils.CfSpaceQuotaType,
		},
	}
}

func (tf *cloudfoundrySpaceQuotaImportProvider) GetImportBlock(data map[string]any, levelId string, filterValues []string) (string, int, error) {
	count := 0
	orgId := levelId

	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.CfSpaceQuotaType, tfutils.OrganizationLevel)
	if err != nil {
		fmt.Print("\r\n")
		log.Fatalf("read doc failed!")
		return "", count, err
	}

	importBlock, count, err := createSpaceQuotaImportBlock(data, orgId, filterValues, resourceDoc)
	if err != nil {
		return "", count, err
	}

	return importBlock, count, nil
}

func createSpaceQuotaImportBlock(data map[string]any, orgId string, filterValues []string, resourceDoc tfutils.EntityDocs) (importBlock string, count int, err error) {
	count = 0
	quotas := data["space_quotas"].([]any)

	if len(filterValues) != 0 {
		var cfAllSpaceQuotas []string

		for x, value := range quotas {
			quota := value.(map[string]any)
			cfAllSpaceQuotas = append(cfAllSpaceQuotas, fmt.Sprintf("%v", quota["name"]))
			if slices.Contains(filterValues, fmt.Sprintf("%v", quota["name"])) {
				importBlock += templateSpaceQuotaImport(x, quota, resourceDoc)
				count++
			}
		}

		missingQuota, subset := isSubset(cfAllSpaceQuotas, filterValues)

		if !subset {
			return "", 0, fmt.Errorf("cloud foundry space quota %s not found in the organization with ID %s. Please adjust it in the provided file", missingQuota, orgId)
		}
	} else {
		for x, value := range quotas {
			quota := value.(map[string]any)
			importBlock += templateSpaceQuotaImport(x, quota, resourceDoc)
			count++
		}
	}

	return importBlock, count, nil
}

func templateSpaceQuotaImport(x int, quota map[string]any, resourceDoc tfutils.EntityDocs) string {
	template := strings.ReplaceAll(resourceDoc.Import, "<resource_name>", "space_quota_"+fmt.Sprintf("%v", x))
	template = strings.ReplaceAll(template, "<space_quota_guid>", fmt.Sprintf("%v", quota["id"]))
	return template + "\n"
}
