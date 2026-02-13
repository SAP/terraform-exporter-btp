package tfimportprovider

import (
	"testing"

	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
	"github.com/stretchr/testify/assert"
)

func TestCreateRoleCollectionImportBlock(t *testing.T) {
	resourceDoc := tfutils.EntityDocs{
		Import: "import {\n\t\t\t\tto = btp_subaccount_role_collection.<resource_name>\n\t\t\t\tid = \"'<subaccount_id>,<name>'\"\n\t\t\t  }\n",
	}

	jsonString := "{\"id\":\"12345\",\"subaccount_id\":\"12345\",\"values\": [{\"description\":\"Administrative access to the subaccount\",\"name\":\"Subaccount Administrator\",\"read_only\": true,\"roles\": [{\"description\":\"Administrative access to service brokers and environments on a subaccount level.\",\"name\":\"Subaccount Service Administrator\",\"role_template_app_id\":\"service-manager!b1476\",\"role_template_name\":\"Subaccount_Service_Administrator\"},        {\"description\":\"Manage authorizations, trusted identity providers, and users.\",\"name\":\"User and Role Administrator\",\"role_template_app_id\":\"xsuaa!t8\",\"role_template_name\":\"xsuaa_admin\"},        {\"description\":\"Manage destination configurations, certificates and signing keys for SAML assertions issued by the Destination service\",\"name\":\"Destination Administrator\",\"role_template_app_id\":\"destination-xsappname!b62\",\"role_template_name\":\"Destination_Administrator\"},        {\"description\":\"Operate the data transmission tunnels and client certificates used by the Cloud connector\",\"name\":\"Cloud Connector Administrator\",\"role_template_app_id\":\"connectivity!b7\",\"role_template_name\":\"Cloud_Connector_Administrator\"},        {\"description\":\"Role for subaccount members with read-write authorizations for core commercialization operations, such as viewing subaccount entitlements, and creating and deleting environment instances.\",\"name\":\"Subaccount Admin\",\"role_template_app_id\":\"cis-local!b4\",\"role_template_name\":\"Subaccount_Admin\"}]}]}"

	data, _ := GetDataFromJsonString(jsonString)

	jsonStringMultiple := "{\"id\":\"12345\",\"subaccount_id\":\"12345\",\"values\": [{\"description\":\"Administrative access to the subaccount\",\"name\":\"Subaccount Administrator\",\"read_only\": true,\"roles\": [{\"description\":\"Administrative access to service brokers and environments on a subaccount level.\",\"name\":\"Subaccount Service Administrator\",\"role_template_app_id\":\"service-manager!b1476\",\"role_template_name\":\"Subaccount_Service_Administrator\"},        {\"description\":\"Manage authorizations, trusted identity providers, and users.\",\"name\":\"User and Role Administrator\",\"role_template_app_id\":\"xsuaa!t8\",\"role_template_name\":\"xsuaa_admin\"},        {\"description\":\"Manage destination configurations, certificates and signing keys for SAML assertions issued by the Destination service\",\"name\":\"Destination Administrator\",\"role_template_app_id\":\"destination-xsappname!b62\",\"role_template_name\":\"Destination_Administrator\"},        {\"description\":\"Operate the data transmission tunnels and client certificates used by the Cloud connector\",\"name\":\"Cloud Connector Administrator\",\"role_template_app_id\":\"connectivity!b7\",\"role_template_name\":\"Cloud_Connector_Administrator\"},        {\"description\":\"Role for subaccount members with read-write authorizations for core commercialization operations, such as viewing subaccount entitlements, and creating and deleting environment instances.\",\"name\":\"Subaccount Admin\",\"role_template_app_id\":\"cis-local!b4\",\"role_template_name\":\"Subaccount_Admin\"}]},    {\"description\":\"Administrative access to service brokers and environments on a subaccount level.\",\"name\":\"Subaccount Service Administrator\",\"read_only\": true,\"roles\": [{\"description\":\"Administrative access to service brokers and environments on a subaccount level.\",\"name\":\"Subaccount_Service_Admin\",\"role_template_app_id\":\"service-manager!b1476\",\"role_template_name\":\"Subaccount_Service_Admin\"}]},    {\"description\":\"Read-only access to the subaccount\",\"name\":\"Subaccount Viewer\",\"read_only\": true,\"roles\": [{\"description\":\"Read-only access for authorizations, trusted identity providers, and users.\",\"name\":\"User and Role Auditor\",\"role_template_app_id\":\"xsuaa!t8\",\"role_template_name\":\"xsuaa_auditor\"},        {\"description\":\"Read-only access to service brokers and environments on a subaccount level.\",\"name\":\"Subaccount Service Auditor\",\"role_template_app_id\":\"service-manager!b1476\",\"role_template_name\":\"Subaccount_Service_Auditor\"},        {\"description\":\"Role for subaccount members with read-only authorizations for core commercialization operations, such as viewing subaccount entitlements, details of environment instances, and job results.\",\"name\":\"Subaccount Viewer\",\"role_template_app_id\":\"cis-local!b4\",\"role_template_name\":\"Subaccount_Viewer\"},        {\"description\":\"View destination configurations, certificates and signing keys for SAML assertions issued by the Destination service\",\"name\":\"Destination Viewer\",\"role_template_app_id\":\"destination-xsappname!b62\",\"role_template_name\":\"Destination_Viewer\"},        {\"description\":\"View the data transmission tunnels and client certificates used by the Cloud connector to communicate with back-end systems\",\"name\":\"Cloud Connector Auditor\",\"role_template_app_id\":\"connectivity!b7\",\"role_template_name\":\"Cloud_Connector_Auditor\"}]}]}"
	dataMultiple, _ := GetDataFromJsonString(jsonStringMultiple)

	tests := []struct {
		name          string
		data          map[string]any
		subaccountId  string
		filterValues  []string
		expectedBlock string
		expectedCount int
		expectError   bool
	}{
		{
			name:          "No filter values",
			data:          data,
			subaccountId:  "12345",
			expectedBlock: "import {\n\t\t\t\tto = btp_subaccount_role_collection.rolecollection_0\n\t\t\t\tid = \"12345,Subaccount Administrator\"\n\t\t\t  }\n\n",
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "With filter values",
			data:          dataMultiple,
			subaccountId:  "12345",
			filterValues:  []string{"subaccount_administrator"},
			expectedBlock: "import {\n\t\t\t\tto = btp_subaccount_role_collection.rolecollection_0\n\t\t\t\tid = \"12345,Subaccount Administrator\"\n\t\t\t  }\n\n",
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "Filter value not found",
			data:          data,
			subaccountId:  "12345",
			filterValues:  []string{"subacccount_viewer"},
			expectedBlock: "",
			expectedCount: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			importBlock, count, err := createRoleCollectionImportBlock(tt.data, tt.subaccountId, tt.filterValues, resourceDoc)
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
