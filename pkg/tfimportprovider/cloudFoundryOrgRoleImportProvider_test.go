package tfimportprovider

import (
	"testing"

	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
	"github.com/stretchr/testify/assert"
)

func TestCreateCfOrgRoleImportBlock(t *testing.T) {
	resourceDoc := tfutils.EntityDocs{
		Import: "import {\n\t\t\t\tto =  cloudfoundry_org_role.<resource_name>\n\t\t\t\tid = \"<role_guid>\"\"\n\t\t\t  }\n",
	}

	jsonString := "{\"org\":\"12345\",\"roles\": [{\"created_at\":\"2024-10-22T19:52:23Z\",\"id\":\"34567\",\"org\":\"12345\",\"type\":\"organization_user\",\"updated_at\":\"2024-10-22T19:52:23Z\",\"user\":\"56789\"}],\"type\": null,\"user\": null}"
	dataOrgRoles, _ := GetDataFromJsonString(jsonString)

	jsonStringMultipleOrgRoles := "{\"org\":\"12345\",\"roles\": [{\"created_at\":\"2024-10-22T19:52:23Z\",\"id\":\"34567\",\"org\":\"12345\",\"type\":\"organization_user\",\"updated_at\":\"2024-10-22T19:52:23Z\",\"user\":\"56789\"},    {\"created_at\":\"2024-10-22T19:52:23Z\",\"id\":\"45678\",\"org\":\"12345\",\"type\":\"organization_manager\",\"updated_at\":\"2024-10-22T19:52:23Z\",\"user\":\"56789\"}],\"type\": null,\"user\": null}"
	dataMultipleOrgRoles, _ := GetDataFromJsonString(jsonStringMultipleOrgRoles)

	tests := []struct {
		name          string
		data          map[string]any
		orgId         string
		filterValues  []string
		expectedBlock string
		expectedCount int
		expectError   bool
	}{

		{
			name:          "Valid data without filter",
			data:          dataOrgRoles,
			orgId:         "12345",
			filterValues:  []string{},
			expectedBlock: "import {\n\t\t\t\tto =  cloudfoundry_org_role.role_organization_user_0\n\t\t\t\tid = \"34567\"\"\n\t\t\t  }\n\n",
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "Valid data with matching filter",
			data:          dataMultipleOrgRoles,
			orgId:         "12345",
			filterValues:  []string{"organization_user_56789"},
			expectedBlock: "import {\n\t\t\t\tto =  cloudfoundry_org_role.role_organization_user_0\n\t\t\t\tid = \"34567\"\"\n\t\t\t  }\n\n",
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "Invalid filter value",
			data:          dataOrgRoles,
			orgId:         "12345",
			filterValues:  []string{"wrong-org-role"},
			expectedBlock: "",
			expectedCount: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			importBlock, count, err := createOrgRoleImportBlock(tt.data, tt.orgId, tt.filterValues, resourceDoc)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBlock, importBlock)
				assert.Equal(t, tt.expectedCount, count)
			}
		})
	}
}
