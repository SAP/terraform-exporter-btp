package tfimportprovider

import (
	"testing"

	tfutils "github.com/SAP/terraform-exporter-btp/tfutils"
	"github.com/stretchr/testify/assert"
)

func TestCreateSecuritySettingImportBlock(t *testing.T) {
	resourceDoc := tfutils.EntityDocs{
		Import: "import {\n\t\t\t\tto = btp_subaccount_security_settings.<resource_name>\n\t\t\t\tid = \"'<subaccount_id>'\"\"\n\t\t\t  }\n",
	}

	jsonString := "{\"access_token_validity\": -1,\"custom_email_domains\": [],\"default_identity_provider\":\"sap.default\",\"iframe_domains\":\"\",\"refresh_token_validity\": -1,\"subaccount_id\":\"12345\",\"treat_users_with_same_email_as_same_user\": false}"

	dataSecuritySettings, _ := GetDataFromJsonString(jsonString)

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
			data:          dataSecuritySettings,
			subaccountId:  "12345",
			filterValues:  []string{},
			expectedBlock: "import {\n\t\t\t\tto = btp_subaccount_security_settings.sec_setting\n\t\t\t\tid = \"12345\"\"\n\t\t\t  }\n\n",
			expectError:   false,
		},
		{
			name:          "With filter values",
			data:          dataSecuritySettings,
			subaccountId:  "12345",
			filterValues:  []string{"12345"},
			expectedBlock: "import {\n\t\t\t\tto = btp_subaccount_security_settings.sec_setting\n\t\t\t\tid = \"12345\"\"\n\t\t\t  }\n\n",
			expectError:   false,
		},
		{
			name:          "Invalid filter value",
			data:          dataSecuritySettings,
			subaccountId:  "12345",
			filterValues:  []string{"56789"},
			expectedBlock: "",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			importBlock, err := createSecuritySettingImportBlock(tt.data, tt.subaccountId, tt.filterValues, resourceDoc)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBlock, importBlock)
			}
		})
	}
}
