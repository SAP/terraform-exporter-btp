package tfimportprovider

import (
	"fmt"
	"slices"
	"strings"

	tfutils "github.com/SAP/terraform-exporter-btp/tfutils"
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

func (tf *subaccountEnvInstanceImportProvider) GetImportBlock(data map[string]interface{}, subaccountId string, filterValues []string) (string, error) {

	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.SubaccountEnvironmentInstanceType)
	if err != nil {
		return "", err
	}

	importBlock, err := createEnvironmentInstanceImportBlock(data, subaccountId, filterValues, resourceDoc)
	if err != nil {
		return "", err
	}

	return importBlock, nil

}

func createEnvironmentInstanceImportBlock(data map[string]interface{}, subaccountId string, filterValues []string, resourceDoc tfutils.EntityDocs) (importBlock string, err error) {
	environmentInstances := data["values"].([]interface{})

	if len(filterValues) != 0 {
		var subaccountAllEnvInstances []string

		for _, value := range environmentInstances {
			environmentInstance := value.(map[string]interface{})
			subaccountAllEnvInstances = append(subaccountAllEnvInstances, fmt.Sprintf("%v", environmentInstance["environment_type"]))
			if slices.Contains(filterValues, fmt.Sprintf("%v", environmentInstance["environment_type"])) {
				importBlock += templateEnvironmentInstanceImport(environmentInstance, subaccountId, resourceDoc)
			}
		}

		missingEnvInstance, subset := isSubset(subaccountAllEnvInstances, filterValues)

		if !subset {
			return "", fmt.Errorf("environment instance %s not found in the subaccount. Please adjust it in the provided file", missingEnvInstance)

		}
	} else {

		for _, value := range environmentInstances {
			environmentInstance := value.(map[string]interface{})
			importBlock += templateEnvironmentInstanceImport(environmentInstance, subaccountId, resourceDoc)
		}
	}
	return importBlock, nil
}

func templateEnvironmentInstanceImport(environmentInstance map[string]interface{}, subaccountId string, resourceDoc tfutils.EntityDocs) string {
	template := strings.Replace(resourceDoc.Import, "<resource_name>", fmt.Sprintf("%v", environmentInstance["environment_type"]), -1)
	template = strings.Replace(template, "<subaccount_id>", subaccountId, -1)
	template = strings.Replace(template, "<environment_instance_id>", fmt.Sprintf("%v", environmentInstance["id"]), -1)
	return template + "\n"
}
