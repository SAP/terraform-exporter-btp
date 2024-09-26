package tfimportprovider

import (
	"testing"

	tfutils "github.com/SAP/terraform-exporter-btp/tfutils"
	"github.com/stretchr/testify/assert"
)

func TestCreateRoleImportBlock(t *testing.T) {
	resourceDoc := tfutils.EntityDocs{
		Import: "import {\n\t\t\t\tto = btp_subaccount_role.<resource_name>\n\t\t\t\tid = \"'<subaccount_id>,<name>,<role_template_name>,<app_id>'\"\n\t\t\t  }\n",
	}

	jsonStringRoles := "{\"id\":\"12345\",\"subaccount_id\":\"12345\",\"values\": [{\"app_id\":\"retention-manager-service!b1824\",\"app_name\":\"retention-manager-service\",\"description\":\"Administrator\",\"name\":\"Administrator\",\"read_only\": true,\"role_template_name\":\"Administrator\",\"scopes\": [{\"custom_grant_as_authority_to_apps\": null,\"custom_granted_apps\": null,\"description\":\"Scope for enabling use of archiving of legal grounds\",\"grant_as_authority_to_apps\": null,\"granted_apps\": null,\"name\":\"retention-manager-service!b1824.business.admin\"}]}]}"
	dataRoles, _ := GetDataFromJsonString(jsonStringRoles)

	jsonStringMultipleRoles := "{\"id\":\"12345\",\"subaccount_id\":\"12345\",\"values\": [{\"app_id\":\"retention-manager-service!b1824\",\"app_name\":\"retention-manager-service\",\"description\":\"Administrator\",\"name\":\"Administrator\",\"read_only\": true,\"role_template_name\":\"Administrator\",\"scopes\": [{\"custom_grant_as_authority_to_apps\": null,\"custom_granted_apps\": null,\"description\":\"Scope for enabling use of archiving of legal grounds\",\"grant_as_authority_to_apps\": null,\"granted_apps\": null,\"name\":\"retention-manager-service!b1824.business.admin\"}]},    {\"app_id\":\"business-context-manager-service!b28142\",\"app_name\":\"business-context-manager-service\",\"description\":\"Business Context Manager Administrator\",\"name\":\"BusinessContextManagerAdministrator\",\"read_only\": true,\"role_template_name\":\"BusinessContextManagerAdministrator\",\"scopes\": [{\"custom_grant_as_authority_to_apps\": null,\"custom_granted_apps\": null,\"description\":\"Business Context Manager Administration\",\"grant_as_authority_to_apps\": null,\"granted_apps\": null,\"name\":\"business-context-manager-service!b28142.BusinessContext.Administrator\"}]}]}"
	dataMultipleRoles, _ := GetDataFromJsonString(jsonStringMultipleRoles)

	tests := []struct {
		name          string
		data          map[string]interface{}
		subaccountId  string
		filterValues  []string
		expectedBlock string
		expectError   bool
	}{
		{
			name:          "No filter values",
			data:          dataRoles,
			subaccountId:  "12345",
			expectedBlock: "import {\n\t\t\t\tto = btp_subaccount_role.administrator\n\t\t\t\tid = \"12345,Administrator,Administrator,retention-manager-service!b1824\"\n\t\t\t  }\n\n",
			expectError:   false,
		},
		{
			name:          "With filter values",
			data:          dataMultipleRoles,
			subaccountId:  "12345",
			filterValues:  []string{"BusinessContextManagerAdministrator"},
			expectedBlock: "import {\n\t\t\t\tto = btp_subaccount_role.businesscontextmanageradministrator\n\t\t\t\tid = \"12345,BusinessContextManagerAdministrator,BusinessContextManagerAdministrator,business-context-manager-service!b28142\"\n\t\t\t  }\n\n",
			expectError:   false,
		},
		{
			name:          "Filter value not found",
			data:          dataRoles,
			subaccountId:  "12345",
			filterValues:  []string{"wrong-role"},
			expectedBlock: "",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			importBlock, err := createRoleImportBlock(tt.data, tt.subaccountId, tt.filterValues, resourceDoc)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBlock, importBlock)
			}
		})
	}
}
