package tfimportprovider

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/SAP/terraform-exporter-btp/pkg/output"
	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
)

type cloudfoundryOrgRolesImportProvider struct {
	TfImportProvider
}

func newCloudfoundryOrgRolesImportProvider() ITfImportProvider {
	return &cloudfoundryOrgRolesImportProvider{
		TfImportProvider: TfImportProvider{
			resourceType: tfutils.CfOrgRoleType,
		},
	}
}

func (tf *cloudfoundryOrgRolesImportProvider) GetImportBlock(data map[string]interface{}, levelId string, filterValues []string) (string, int, error) {
	count := 0
	orgId := levelId
	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.CfOrgRoleType, tfutils.OrganizationLevel)
	if err != nil {
		fmt.Print("\r\n")
		log.Fatalf("read doc failed!")
		return "", count, err
	}
	importBlock, count, err := createOrgRoleImportBlock(data, orgId, filterValues, resourceDoc)
	if err != nil {
		return "", count, err
	}
	return importBlock, count, nil
}
func createOrgRoleImportBlock(data map[string]interface{}, orgId string, filterValues []string, resourceDoc tfutils.EntityDocs) (importBlock string, count int, err error) {
	count = 0
	roles := data["roles"].([]interface{})
	if len(filterValues) != 0 {
		var cfAllOrgRoles []string
		for x, value := range roles {
			role := value.(map[string]interface{})
			cfAllOrgRoles = append(cfAllOrgRoles, output.FormatOrgRoleResourceName(fmt.Sprintf("%v", role["type"]), fmt.Sprintf("%v", role["user"])))
			if slices.Contains(filterValues, output.FormatOrgRoleResourceName(fmt.Sprintf("%v", role["type"]), fmt.Sprintf("%v", role["user"]))) {
				importBlock += templateOrgRoleImport(x, role, resourceDoc)
				count++
			}
		}
		missingRole, subset := isSubset(cfAllOrgRoles, filterValues)
		if !subset {
			return "", 0, fmt.Errorf("cloud foudndry org role %s not found in the organization with ID %s. Please adjust it in the provided file", missingRole, orgId)
		}
	} else {
		for x, value := range roles {
			role := value.(map[string]interface{})
			importBlock += templateOrgRoleImport(x, role, resourceDoc)
			count++
		}
	}
	return importBlock, count, nil
}
func templateOrgRoleImport(x int, role map[string]interface{}, resourceDoc tfutils.EntityDocs) string {
	template := strings.Replace(resourceDoc.Import, "<resource_name>", "role_"+fmt.Sprintf("%v", role["type"])+"_"+fmt.Sprintf("%v", x), -1)
	template = strings.Replace(template, "<role_guid>", fmt.Sprintf("%v", role["id"]), -1)
	return template + "\n"
}
