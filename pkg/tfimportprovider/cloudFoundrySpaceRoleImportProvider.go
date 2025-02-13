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
	spaceId := levelId
	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.CfSpaceRoleType, tfutils.OrganizationLevel)
	if err != nil {
		fmt.Print("\r\n")
		log.Fatalf("read doc failed!")
		return "", count, err
	}
	importBlock, count, err := createSpaceRoleImportBlock(data, spaceId, filterValues, resourceDoc)
	if err != nil {
		return "", count, err
	}
	return importBlock, count, nil
}
func createSpaceRoleImportBlock(data map[string]interface{}, spaceId string, filterValues []string, resourceDoc tfutils.EntityDocs) (importBlock string, count int, err error) {
	count = 0
	roles := data["roles"].([]interface{})
	if len(filterValues) != 0 {
		var cfAllSpaceRoles []string
		for x, value := range roles {
			role := value.(map[string]interface{})
			var formattedRoleName = output.FormatRoles
			cfAllSpaceRoles = append(cfAllSpaceRoles, formattedRoleName(fmt.Sprintf("%v", role["type"]), fmt.Sprintf("%v", role["space"]), fmt.Sprintf("%v", role["user"])))
			if slices.Contains(filterValues, formattedRoleName(fmt.Sprintf("%v", role["type"]), fmt.Sprintf("%v", role["space"]), fmt.Sprintf("%v", role["user"]))) {
				importBlock += templateSpaceRoleImport(x, role, resourceDoc)
				count++
			}
		}
		missingRole, subset := isSubset(cfAllSpaceRoles, filterValues)
		if !subset {
			return "", 0, fmt.Errorf("cloud foudndry space role %s not found in the space with ID %s. Please adjust it in the provided file", missingRole, spaceId)
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
	template := strings.Replace(resourceDoc.Import, "<resource_name>", "role_"+fmt.Sprintf("%v", role["id"])+"_"+fmt.Sprintf("%v", role["type"])+"_"+fmt.Sprintf("%v", x), -1)
	template = strings.Replace(template, "<role_guid>", fmt.Sprintf("%v", role["id"]), -1)
	return template + "\n"
}
