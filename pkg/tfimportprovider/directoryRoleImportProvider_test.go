package tfimportprovider

import (
	"testing"

	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
	"github.com/stretchr/testify/assert"
)

func TestCreateDirectoryRoleImportBlock(t *testing.T) {
	resourceDoc := tfutils.EntityDocs{
		Import: "import {\n\t\t\t\tto = btp_directory_role.<resource_name>\n\t\t\t\tid = \"'<directory_id>,<name>,<role_template_name>,<app_id>'\"\n\t\t\t  }\n",
	}

	jsonStringRoles := "{\"directory_id\":\"12345\",\"id\":\"12345\",\"values\": [{\"app_id\":\"uas!b36585\",\"app_name\":\"uas\",\"description\":\"Role for directory members with read-only authorizations for core commercialization operations, such as viewing directory usage information.\",\"name\":\"Directory Usage Reporting Viewer\",\"read_only\": true,\"role_template_name\":\"Directory_Usage_Reporting_Viewer\",\"scopes\": [{\"custom_grant_as_authority_to_apps\": null,\"custom_granted_apps\": null,\"description\":\"\",\"grant_as_authority_to_apps\": null,\"granted_apps\": null,\"name\":\"uas!b36585.UAS.reporting.directory.read\"},{\"custom_grant_as_authority_to_apps\": null,\"custom_granted_apps\": null,\"description\":\"Enable account navigation\",\"grant_as_authority_to_apps\": null,\"granted_apps\": [\"*\"],\"name\":\"xs_account.access\"}]}]}"
	dataRoles, _ := GetDataFromJsonString(jsonStringRoles)

	jsonStringMultipleRoles := "{\"directory_id\":\"12345\",\"id\":\"12345\",\"values\": [{\"app_id\":\"uas!b36585\",\"app_name\":\"uas\",\"description\":\"Role for directory members with read-only authorizations for core commercialization operations, such as viewing directory usage information.\",\"name\":\"Directory Usage Reporting Viewer\",\"read_only\": true,\"role_template_name\":\"Directory_Usage_Reporting_Viewer\",\"scopes\": [{\"custom_grant_as_authority_to_apps\": null,\"custom_granted_apps\": null,\"description\":\"\",\"grant_as_authority_to_apps\": null,\"granted_apps\": null,\"name\":\"uas!b36585.UAS.reporting.directory.read\"},{\"custom_grant_as_authority_to_apps\": null,\"custom_granted_apps\": null,\"description\":\"Enable account navigation\",\"grant_as_authority_to_apps\": null,\"granted_apps\": [\"*\"],\"name\":\"xs_account.access\"}]}, {\"app_id\":\"uas!b36585\",\"app_name\":\"uas\",\"description\":\"Role for directory members .\",\"name\":\"Directory Role\",\"read_only\": true,\"role_template_name\":\"Directory_Role\",\"scopes\": [{\"custom_grant_as_authority_to_apps\": null,\"custom_granted_apps\": null,\"description\":\"\",\"grant_as_authority_to_apps\": null,\"granted_apps\": null,\"name\":\"uas!b36585.UAS.reporting.directory.read\"},{\"custom_grant_as_authority_to_apps\": null,\"custom_granted_apps\": null,\"description\":\"Enable account navigation\",\"grant_as_authority_to_apps\": null,\"granted_apps\": [\"*\"],\"name\":\"xs_account.access\"}]}]}"
	dataMultipleRoles, _ := GetDataFromJsonString(jsonStringMultipleRoles)

	tests := []struct {
		name          string
		data          map[string]interface{}
		directoryId   string
		filterValues  []string
		expectedBlock string
		expectedCount int
		expectError   bool
	}{
		{
			name:          "No filter values",
			data:          dataRoles,
			directoryId:   "12345",
			expectedBlock: "import {\n\t\t\t\tto = btp_directory_role.directory_role_0\n\t\t\t\tid = \"12345,Directory Usage Reporting Viewer,Directory_Usage_Reporting_Viewer,uas!b36585\"\n\t\t\t  }\n\n",
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "With filter values",
			data:          dataMultipleRoles,
			directoryId:   "12345",
			filterValues:  []string{"directory_usage_reporting_viewer"},
			expectedBlock: "import {\n\t\t\t\tto = btp_directory_role.directory_role_0\n\t\t\t\tid = \"12345,Directory Usage Reporting Viewer,Directory_Usage_Reporting_Viewer,uas!b36585\"\n\t\t\t  }\n\n",
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "Filter value not found",
			data:          dataRoles,
			directoryId:   "12345",
			filterValues:  []string{"wrong-role"},
			expectedBlock: "",
			expectedCount: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			importBlock, count, err := createDirectoryRoleImportBlock(tt.data, tt.directoryId, tt.filterValues, resourceDoc)
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
