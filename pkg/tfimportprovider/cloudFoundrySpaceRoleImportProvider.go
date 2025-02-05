package tfimportprovider

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/SAP/terraform-exporter-btp/pkg/output"
	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
)

type cloudfoundrySpaceRolesImportProvider struct {
	TfImportProvider
}

func newCloudfoundrySpaceRolesImportProvider() ITfImportProvider {
	return &cloudfoundrySpaceRolesImportProvider{
		TfImportProvider: TfImportProvider{
			resourceType: tfutils.CfSpaceRoleType,
		},
	}
}

func (tf *cloudfoundrySpaceRolesImportProvider) GetImportBlock(data map[string]interface{}, levelId string, filterValues []string) (string, int, error) {
	count := 0
	orgId := levelId
	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.CfSpaceRoleType, tfutils.OrganizationLevel)
	if err != nil {
		fmt.Print("\r\n")
		log.Fatalf("read doc failed!")
		return "", count, err
	}
	importBlock, count, err := createSpaceRoleImportBlock(data, orgId, filterValues, resourceDoc)
	if err != nil {
		return "", count, err
	}
	return importBlock, count, nil
}
func createSpaceRoleImportBlock(data map[string]interface{}, orgId string, filterValues []string, resourceDoc tfutils.EntityDocs) (importBlock string, count int, err error) {
	count = 0
	roles := data["roles"].([]interface{})
	if len(filterValues) != 0 {
		var cfAllSpaceRoles []string
		for x, value := range roles {
			role := value.(map[string]interface{})
			cfAllSpaceRoles = append(cfAllSpaceRoles, output.FormatSpaceRoleResourceName(fmt.Sprintf("%v", role["type"]), fmt.Sprintf("%v", role["space"]), fmt.Sprintf("%v", role["user"])))
			if slices.Contains(filterValues, output.FormatSpaceRoleResourceName(fmt.Sprintf("%v", role["type"]), fmt.Sprintf("%v", role["space"]), fmt.Sprintf("%v", role["user"]))) {
				importBlock += templateSpaceRoleImport(x, role, resourceDoc)
				count++
			}
		}
		missingRole, subset := isSubset(cfAllSpaceRoles, filterValues)
		if !subset {
			return "", 0, fmt.Errorf("cloud foudndry org role %s not found in the organization with ID %s. Please adjust it in the provided file", missingRole, orgId)
		}
	} else {
		for x, value := range roles {
			role := value.(map[string]interface{})
			importBlock += templateSpaceRoleImport(x, role, resourceDoc)
			count++
		}
	}
	return importBlock, count, nil
}
func templateSpaceRoleImport(x int, role map[string]interface{}, resourceDoc tfutils.EntityDocs) string {
	//Needs to be removed once import statement in document is fixed
	template := strings.Replace(resourceDoc.Import, "cloudfoundry_role", "cloudfoundry_space_role", -1)

	template = strings.Replace(template, "<resource_name>", "role_"+fmt.Sprintf("%v", role["id"])+"_"+fmt.Sprintf("%v", role["type"])+"_"+fmt.Sprintf("%v", x), -1)
	template = strings.Replace(template, "<role_guid>", fmt.Sprintf("%v", role["id"]), -1)
	return template + "\n"
}
