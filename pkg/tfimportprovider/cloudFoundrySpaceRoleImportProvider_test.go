package tfimportprovider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/SAP/terraform-exporter-btp/pkg/output"
	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
	"github.com/stretchr/testify/assert"
)

func TestCreateCfSpaceRoleImportBlock(t *testing.T) {
	resourceDoc := tfutils.EntityDocs{
		Import: "import {\n\t\t\t\tto =  cloudfoundry_space_role.<resource_name>\n\t\t\t\tid = \"<role_guid>\"\"\n\t\t\t  }\n",
	}

	jsonString := "{\"roles\": [{\"created_at\":\"2025-01-07T08:31:08Z\",\"id\":\"23456\",\"space\":\"12345\",\"type\":\"space_manager\",\"updated_at\":\"2025-01-07T08:31:08Z\",\"user\":\"34567\"}],\"space\":\"12345\"}"
	dataSpaceRoles, _ := GetDataFromJsonString(jsonString)

	jsonStringMultipleSpaceRoles := "{\"roles\": [{\"created_at\":\"2025-01-07T08:31:08Z\",\"id\":\"23456\",\"space\":\"12345\",\"type\":\"space_manager\",\"updated_at\":\"2025-01-07T08:31:08Z\",\"user\":\"34567\"},    {\"created_at\":\"2025-01-07T08:31:08Z\",\"id\":\"56789\",\"space\":\"12345\",\"type\":\"space_developer\",\"updated_at\":\"2025-01-07T08:31:08Z\",\"user\":\"34567\"}],\"space\":\"12345\"}"
	dataMultipleSpaceRoles, _ := GetDataFromJsonString(jsonStringMultipleSpaceRoles)

	tests := []struct {
		name          string
		data          map[string]any
		spaceId       string
		filterValues  []string
		expectedBlock string
		expectedCount int
		expectError   bool
	}{

		{
			name:          "Valid data without filter",
			data:          dataSpaceRoles,
			spaceId:       "12345",
			filterValues:  []string{},
			expectedBlock: "import {\n\t\t\t\tto =  cloudfoundry_space_role.role_23456_space_manager_0\n\t\t\t\tid = \"23456\"\"\n\t\t\t  }\n\n",
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "Valid data with matching filter",
			data:          dataMultipleSpaceRoles,
			spaceId:       "12345",
			filterValues:  []string{"test-space_space-manager_test-user_12345"},
			expectedBlock: "import {\n\t\t\t\tto =  cloudfoundry_space_role.role_23456_space_manager_0\n\t\t\t\tid = \"23456\"\"\n\t\t\t  }\n\n",
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "Invalid filter value",
			data:          dataSpaceRoles,
			spaceId:       "12345",
			filterValues:  []string{"wrong_org_role"},
			expectedBlock: "",
			expectedCount: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.filterValues != nil {
				output.FormatRoles = func(spaceRoleType string, spaceId string, userId string) string {
					return fmt.Sprintf("%s_%s_%s_%s", "test-space", strings.ReplaceAll(spaceRoleType, "_", "-"), "test-user", spaceId)

				}
			}

			importBlock, count, err := createSpaceRoleImportBlock(tt.data, tt.spaceId, tt.filterValues, resourceDoc)
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
