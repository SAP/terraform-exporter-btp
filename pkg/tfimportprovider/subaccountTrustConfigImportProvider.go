package tfimportprovider

import (
	"fmt"
	"log"
	"slices"
	"strings"

	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
)

type subaccountTrustConfigImportProvider struct {
	TfImportProvider
}

func newSubaccountTrustConfigImportProvider() ITfImportProvider {
	return &subaccountTrustConfigImportProvider{
		TfImportProvider: TfImportProvider{
			resourceType: tfutils.SubaccountTrustConfigurationType,
		},
	}
}

func (tf *subaccountTrustConfigImportProvider) GetImportBlock(data map[string]interface{}, levelId string, filterValues []string) (string, int, error) {
	count := 0
	subaccountId := levelId

	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.SubaccountTrustConfigurationType, tfutils.SubaccountLevel)
	if err != nil {
		fmt.Print("\r\n")
		log.Fatalf("read doc failed!")
		return "", count, err
	}

	importBlock, count, err := createTrustConfigurationImportBlock(data, subaccountId, filterValues, resourceDoc)
	if err != nil {
		return "", count, err
	}

	return importBlock, count, nil
}

func createTrustConfigurationImportBlock(data map[string]interface{}, subaccountId string, filterValues []string, resourceDoc tfutils.EntityDocs) (importBlock string, count int, err error) {
	count = 0
	trusts := data["values"].([]interface{})

	if len(filterValues) != 0 {
		var subaccountAllTrusts []string

		for x, value := range trusts {
			trust := value.(map[string]interface{})
			subaccountAllTrusts = append(subaccountAllTrusts, fmt.Sprintf("%v", trust["origin"]))
			if slices.Contains(filterValues, fmt.Sprintf("%v", trust["origin"])) {
				importBlock += templateTrustImport(x, trust, subaccountId, resourceDoc)
				count++
			}
		}

		missingTrust, subset := isSubset(subaccountAllTrusts, filterValues)

		if !subset {
			return "", 0, fmt.Errorf("trust configuration %s not found in the subaccount. Please adjust it in the provided file", missingTrust)
		}
	} else {
		for x, value := range trusts {
			trust := value.(map[string]interface{})
			importBlock += templateTrustImport(x, trust, subaccountId, resourceDoc)
			count++
		}
	}

	return importBlock, count, nil
}

func templateTrustImport(x int, trust map[string]interface{}, subaccountId string, resourceDoc tfutils.EntityDocs) string {
	template := strings.Replace(resourceDoc.Import, "<resource_name>", "trust_"+fmt.Sprint(x), -1)
	template = strings.Replace(template, "<subaccount_id>", subaccountId, -1)
	template = strings.Replace(template, "<origin>", fmt.Sprintf("%v", trust["origin"]), -1)
	return template + "\n"
}
