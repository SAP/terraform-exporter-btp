package tfimportprovider

import (
	"fmt"
	"log"
	"slices"
	"strings"

	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
)

type cloudfoundryUserImportProvider struct {
	TfImportProvider
}

func newcloudfoundryUserImportProvider() ITfImportProvider {
	return &cloudfoundryUserImportProvider{
		TfImportProvider: TfImportProvider{
			resourceType: tfutils.CfUserType,
		},
	}
}

func (tf *cloudfoundryUserImportProvider) GetImportBlock(data map[string]interface{}, levelId string, filterValues []string) (string, int, error) {
	count := 0
	orgId := levelId

	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.CfUserType, tfutils.OrganizationLevel)
	if err != nil {
		fmt.Print("\r\n")
		log.Fatalf("read doc failed!")
		return "", count, err
	}

	importBlock, count, err := createUserImportBlock(data, orgId, filterValues, resourceDoc)
	if err != nil {
		return "", count, err
	}

	return importBlock, count, nil
}

func createUserImportBlock(data map[string]interface{}, orgId string, filterValues []string, resourceDoc tfutils.EntityDocs) (importBlock string, count int, err error) {
	count = 0
	users := data["users"].([]interface{})

	if len(filterValues) != 0 {
		var cfAllUsers []string

		for x, value := range users {
			user := value.(map[string]interface{})
			cfAllUsers = append(cfAllUsers, fmt.Sprintf("%v", user["username"]))
			if slices.Contains(filterValues, fmt.Sprintf("%v", user["username"])) {
				importBlock += templateUserImport(x, user, resourceDoc)
				count++
			}
		}

		missingSpace, subset := isSubset(cfAllUsers, filterValues)

		if !subset {
			return "", 0, fmt.Errorf("cloud foudndry user %s not found in the organization with ID %s. Please adjust it in the provided file", missingSpace, orgId)
		}
	} else {
		for x, value := range users {
			user := value.(map[string]interface{})
			importBlock += templateUserImport(x, user, resourceDoc)
			count++
		}
	}

	return importBlock, count, nil
}

func templateUserImport(x int, user map[string]interface{}, resourceDoc tfutils.EntityDocs) string {
	template := strings.Replace(resourceDoc.Import, "<resource_name>", "user_"+fmt.Sprintf("%v", x), -1)
	template = strings.Replace(template, "<user_guid>", fmt.Sprintf("%v", user["id"]), -1)
	return template + "\n"
}
