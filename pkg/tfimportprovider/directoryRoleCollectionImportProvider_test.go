package tfimportprovider

import (
	"testing"

	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
	"github.com/stretchr/testify/assert"
)

func TestCreateDirectoryRoleCollectionImportBlock(t *testing.T) {
	resourceDoc := tfutils.EntityDocs{
		Import: "import {\n\t\t\t\tto = btp_directory_role_collection.<resource_name>\n\t\t\t\tid = \"'<directory_id>,<name>'\"\n\t\t\t  }\n",
	}

	jsonString := "{\"directory_id\":\"12345\",\"id\":\"12345\",\"values\": [{\"description\":\"Administrative access to the directory\",\"name\":\"Directory Administrator\",\"read_only\": true,\"roles\": [{\"description\":\"Manage authorizations, trusted identity providers, and users.\",\"name\":\"User and Role Administrator\",\"role_template_app_id\":\"xsuaa!t1\",\"role_template_name\":\"xsuaa_admin\"},        {\"description\":\"Role for directory members with read-only authorizations for core commercialization operations, such as viewing directory usage information.\",\"name\":\"Directory Usage Reporting Viewer\",\"role_template_app_id\":\"uas!b36585\",\"role_template_name\":\"Directory_Usage_Reporting_Viewer\"},        {\"description\":\"Role for directory members with read-write authorizations for core commercialization operations, such as updating directories, setting entitlements, and creating, updating, and deleting subaccounts.\",\"name\":\"Directory Admin\",\"role_template_app_id\":\"cis-central!b14\",\"role_template_name\":\"Directory_Admin\"}]}]}"
	data, _ := GetDataFromJsonString(jsonString)

	jsonStringMultiple := "{\"directory_id\":\"12345\",\"id\":\"12345\",\"values\": [{\"description\":\"Administrative access to the directory\",\"name\":\"Directory Administrator\",\"read_only\": true,\"roles\": [{\"description\":\"Manage authorizations, trusted identity providers, and users.\",\"name\":\"User and Role Administrator\",\"role_template_app_id\":\"xsuaa!t1\",\"role_template_name\":\"xsuaa_admin\"},        {\"description\":\"Role for directory members with read-only authorizations for core commercialization operations, such as viewing directory usage information.\",\"name\":\"Directory Usage Reporting Viewer\",\"role_template_app_id\":\"uas!b36585\",\"role_template_name\":\"Directory_Usage_Reporting_Viewer\"},        {\"description\":\"Role for directory members with read-write authorizations for core commercialization operations, such as updating directories, setting entitlements, and creating, updating, and deleting subaccounts.\",\"name\":\"Directory Admin\",\"role_template_app_id\":\"cis-central!b14\",\"role_template_name\":\"Directory_Admin\"}]},    {\"description\":\"Read-only access to the directory\",\"name\":\"Directory Viewer\",\"read_only\": true,\"roles\": [{\"description\":\"Read-only access for authorizations, trusted identity providers, and users.\",\"name\":\"User and Role Auditor\",\"role_template_app_id\":\"xsuaa!t1\",\"role_template_name\":\"xsuaa_auditor\"},        {\"description\":\"Role for directory members with read-only authorizations for core commercialization operations, such as viewing directories, subaccounts, entitlements, and regions.\",\"name\":\"Directory Viewer\",\"role_template_app_id\":\"cis-central!b14\",\"role_template_name\":\"Directory_Viewer\"},        {\"description\":\"Role for directory members with read-only authorizations for core commercialization operations, such as viewing directory usage information.\",\"name\":\"Directory Usage Reporting Viewer\",\"role_template_app_id\":\"uas!b36585\",\"role_template_name\":\"Directory_Usage_Reporting_Viewer\"}]}]}"
	dataMultiple, _ := GetDataFromJsonString(jsonStringMultiple)

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
			data:          data,
			directoryId:   "12345",
			expectedBlock: "import {\n\t\t\t\tto = btp_directory_role_collection.directory_administrator\n\t\t\t\tid = \"12345,Directory Administrator\"\n\t\t\t  }\n\n",
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "With filter values",
			data:          dataMultiple,
			directoryId:   "12345",
			filterValues:  []string{"directory_administrator"},
			expectedBlock: "import {\n\t\t\t\tto = btp_directory_role_collection.directory_administrator\n\t\t\t\tid = \"12345,Directory Administrator\"\n\t\t\t  }\n\n",
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "Filter value not found",
			data:          data,
			directoryId:   "12345",
			filterValues:  []string{"invalid_rolecollection"},
			expectedBlock: "",
			expectedCount: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			importBlock, count, err := createDirectoryRoleCollectionImportBlock(tt.data, tt.directoryId, tt.filterValues, resourceDoc)
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
