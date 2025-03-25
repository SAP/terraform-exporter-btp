package tfimportprovider

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/SAP/terraform-exporter-btp/pkg/output"
	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
)

type cloudfoundryServiceInstanceImportProvider struct {
	TfImportProvider
}

func newCloudfoundryServiceInstanceImportProvider() ITfImportProvider {
	return &cloudfoundryServiceInstanceImportProvider{
		TfImportProvider: TfImportProvider{
			resourceType: tfutils.CfServiceInstanceType,
		},
	}
}

func (tf *cloudfoundryServiceInstanceImportProvider) GetImportBlock(data map[string]interface{}, levelId string, filterValues []string) (string, int, error) {
	count := 0
	orgId := levelId

	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.CfServiceInstanceType, tfutils.OrganizationLevel)
	if err != nil {
		fmt.Print("\r\n")
		log.Fatalf("read doc failed!")
		return "", count, err
	}

	importBlock, count, err := createCfServiceInstanceImportBlock(data, orgId, filterValues, resourceDoc)
	if err != nil {
		return "", count, err
	}

	return importBlock, count, nil
}

func createCfServiceInstanceImportBlock(data map[string]interface{}, orgId string, filterValues []string, resourceDoc tfutils.EntityDocs) (importBlock string, count int, err error) {
	count = 0
	serviceInstances := data["service_instances"].([]interface{})

	if len(filterValues) != 0 {
		var cfAllServiceInstances []string

		for x, value := range serviceInstances {
			serviceInstance := value.(map[string]interface{})
			resourceName := output.FormatServiceInstanceResourceName(fmt.Sprintf("%v", serviceInstance["name"]), fmt.Sprintf("%v", serviceInstance["service_plan"]))
			cfAllServiceInstances = append(cfAllServiceInstances, resourceName)
			if slices.Contains(filterValues, resourceName) {
				importBlock += templateCfServiceInstanceImport(x, serviceInstance, resourceDoc)
				count++
			}
		}

		missingServiceInstance, subset := isSubset(cfAllServiceInstances, filterValues)

		if !subset {
			return "", 0, fmt.Errorf("cloud foundry service instance %s not found in the organization with ID %s. Please adjust it in the provided file", missingServiceInstance, orgId)
		}
	} else {
		for x, value := range serviceInstances {
			serviceInstance := value.(map[string]interface{})
			importBlock += templateCfServiceInstanceImport(x, serviceInstance, resourceDoc)
			count++
		}
	}

	return importBlock, count, nil
}

func templateCfServiceInstanceImport(x int, serviceInstance map[string]interface{}, resourceDoc tfutils.EntityDocs) string {
	template := strings.ReplaceAll(resourceDoc.Import, "<resource_name>", "cf_serviceinstance_"+fmt.Sprintf("%v", x))
	template = strings.ReplaceAll(template, "<service_instance_guid>", fmt.Sprintf("%v", serviceInstance["id"]))
	return template + "\n"
}
