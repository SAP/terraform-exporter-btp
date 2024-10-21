package tfimportprovider

import (
	"testing"

	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
	"github.com/stretchr/testify/assert"
)

func TestCreateTrustConfigurationImportBlock(t *testing.T) {
	resourceDoc := tfutils.EntityDocs{
		Import: "import {\n\t\t\t\tto = btp_subaccount_trust_configuration.<resource_name>\n\t\t\t\tid = \"<subaccount_id>,<origin>\"\"\n\t\t\t  }\n",
	}

	jsonString := "{\"id\":\"12345\",\"subaccount_id\":\"12345\",\"values\": [{\"auto_create_shadow_users\": false,\"available_for_user_logon\": false,\"description\":\"Custom Platform Identity Provider\",\"domain\": null,\"id\":\"terraformtest-platform\",\"identity_provider\":\"terraformtest.accounts400.ondemand.com\",\"link_text\":\"\",\"name\":\"terraformtest-platform\",\"origin\":\"terraformtest-platform\",\"protocol\":\"OpenID Connect\",\"read_only\": true,\"status\":\"active\",\"type\":\"Platform\"}]}"
	dataTrustConfig, _ := GetDataFromJsonString(jsonString)

	jsonStringMultipleTrustConfig := "{\"id\":\"12345\",\"subaccount_id\":\"12345\",\"values\": [{\"auto_create_shadow_users\": false,\"available_for_user_logon\": true,\"description\":\"\",\"domain\": null,\"id\":\"sap.default\",\"identity_provider\":\"\",\"link_text\":\"Default Identity Provider\",\"name\":\"sap.default\",\"origin\":\"sap.default\",\"protocol\":\"OpenID Connect\",\"read_only\": false,\"status\":\"active\",\"type\":\"Application\"},    {\"auto_create_shadow_users\": false,\"available_for_user_logon\": false,\"description\":\"Identity Authentication tenant terraformeds.accounts.ondemand.com used for platform users\",\"domain\":\"terraformeds.accounts.ondemand.com\",\"id\":\"terraformeds-platform\",\"identity_provider\":\"terraformeds.accounts.ondemand.com\",\"link_text\":\"\",\"name\":\"terraformeds.accounts.ondemand.com (platform users)\",\"origin\":\"terraformeds-platform\",\"protocol\":\"OpenID Connect\",\"read_only\": true,\"status\":\"active\",\"type\":\"Platform\"},    {\"auto_create_shadow_users\": false,\"available_for_user_logon\": false,\"description\":\"Custom Platform Identity Provider\",\"domain\": null,\"id\":\"terraformtest-platform\",\"identity_provider\":\"terraformtest.accounts400.ondemand.com\",\"link_text\":\"\",\"name\":\"terraformtest-platform\",\"origin\":\"terraformtest-platform\",\"protocol\":\"OpenID Connect\",\"read_only\": true,\"status\":\"active\",\"type\":\"Platform\"}]}"
	dataMultipleTrust, _ := GetDataFromJsonString(jsonStringMultipleTrustConfig)

	tests := []struct {
		name          string
		data          map[string]interface{}
		subaccountId  string
		filterValues  []string
		expectedBlock string
		expectedCount int
		expectError   bool
	}{

		{
			name:          "Valid data without filter",
			data:          dataTrustConfig,
			subaccountId:  "12345",
			filterValues:  []string{},
			expectedBlock: "import {\n\t\t\t\tto = btp_subaccount_trust_configuration.trust0\n\t\t\t\tid = \"12345,terraformtest-platform\"\"\n\t\t\t  }\n\n",
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "Valid data with matching filter",
			data:          dataMultipleTrust,
			subaccountId:  "12345",
			filterValues:  []string{"terraformtest-platform"},
			expectedBlock: "import {\n\t\t\t\tto = btp_subaccount_trust_configuration.trust2\n\t\t\t\tid = \"12345,terraformtest-platform\"\"\n\t\t\t  }\n\n",
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "Invalid filter value",
			data:          dataTrustConfig,
			subaccountId:  "12345",
			filterValues:  []string{"wrong-trust-name"},
			expectedBlock: "",
			expectedCount: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			importBlock, count, err := createTrustConfigurationImportBlock(tt.data, tt.subaccountId, tt.filterValues, resourceDoc)
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
