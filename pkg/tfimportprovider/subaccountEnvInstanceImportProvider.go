package tfimportprovider

import (
	"fmt"
	"slices"
	"strings"

	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
)

type subaccountEnvInstanceImportProvider struct {
	TfImportProvider
}

func newSubaccountEnvInstanceImportProvider() ITfImportProvider {
	return &subaccountEnvInstanceImportProvider{
		TfImportProvider: TfImportProvider{
			resourceType: tfutils.SubaccountEnvironmentInstanceType,
		},
	}
}

func (tf *subaccountEnvInstanceImportProvider) GetImportBlock(data map[string]interface{}, levelId string, filterValues []string) (string, int, error) {
	count := 0
	subaccountId := levelId

	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.SubaccountEnvironmentInstanceType, levelId)
	if err != nil {
		return "", count, err
	}

	importBlock, count, err := createEnvironmentInstanceImportBlock(data, subaccountId, filterValues, resourceDoc)
	if err != nil {
		return "", count, err
	}

	return importBlock, count, nil

}

func createEnvironmentInstanceImportBlock(data map[string]interface{}, subaccountId string, filterValues []string, resourceDoc tfutils.EntityDocs) (importBlock string, count int, err error) {
	count = 0
	environmentInstances := data["values"].([]interface{})

	if len(filterValues) != 0 {
		var subaccountAllEnvInstances []string

		for _, value := range environmentInstances {
			environmentInstance := value.(map[string]interface{})
			subaccountAllEnvInstances = append(subaccountAllEnvInstances, fmt.Sprintf("%v", environmentInstance["environment_type"]))
			if slices.Contains(filterValues, fmt.Sprintf("%v", environmentInstance["environment_type"])) {
				importBlock += templateEnvironmentInstanceImport(environmentInstance, subaccountId, resourceDoc)
				count++
			}
		}

		missingEnvInstance, subset := isSubset(subaccountAllEnvInstances, filterValues)

		if !subset {
			return "", 0, fmt.Errorf("environment instance %s not found in the subaccount. Please adjust it in the provided file", missingEnvInstance)

		}
	} else {

		for _, value := range environmentInstances {
			environmentInstance := value.(map[string]interface{})
			importBlock += templateEnvironmentInstanceImport(environmentInstance, subaccountId, resourceDoc)
			count++
		}
	}
	return importBlock, count, nil
}

func templateEnvironmentInstanceImport(environmentInstance map[string]interface{}, subaccountId string, resourceDoc tfutils.EntityDocs) string {
	template := strings.Replace(resourceDoc.Import, "<resource_name>", fmt.Sprintf("%v", environmentInstance["environment_type"]), -1)
	template = strings.Replace(template, "<subaccount_id>", subaccountId, -1)
	template = strings.Replace(template, "<environment_instance_id>", fmt.Sprintf("%v", environmentInstance["id"]), -1)
	return template + "\n"
}
