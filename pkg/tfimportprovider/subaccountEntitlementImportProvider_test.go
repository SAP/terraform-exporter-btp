package tfimportprovider

import (
	"testing"

	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
	"github.com/stretchr/testify/assert"
)

func TestCreateEntitlementImportBlock(t *testing.T) {
	resourceDoc := tfutils.EntityDocs{
		Import: "import {\n\t\t\t\tto = btp_subaccount_entitlement.<resource_name>\n\t\t\t\tid = \"<subaccount_id>,<service_name>,<plan_name>\"\n\t\t\t  }\n",
	}

	jsonString := "{\"application-logs:lite\": {\"category\":\"ELASTIC_SERVICE\",\"plan_description\":\"Free offering for development purposes\",\"plan_display_name\":\"lite\",\"plan_name\":\"lite\",\"quota_assigned\": 1,\"quota_remaining\": 1,\"service_display_name\":\"Application Logging Service\",\"service_name\":\"application-logs\"}}"
	dataEntitlement, _ := GetDataFromJsonString(jsonString)

	jsonStrinMultipleEntitlements := "{\"application-logs:lite\": {\"category\":\"ELASTIC_SERVICE\",\"plan_description\":\"Free offering for development purposes\",\"plan_display_name\":\"lite\",\"plan_name\":\"lite\",\"quota_assigned\": 1,\"quota_remaining\": 1,\"service_display_name\":\"Application Logging Service\",\"service_name\":\"application-logs\"},\"auditlog-api:default\": {\"category\":\"ELASTIC_SERVICE\",\"plan_description\":\"Default plan for Auditlog API\",\"plan_display_name\":\"Default\",\"plan_name\":\"default\",\"quota_assigned\": 1,\"quota_remaining\": 1,\"service_display_name\":\"[DEPRECATED] Audit Log Retrieval\",\"service_name\":\"auditlog-api\"}}"
	dataMultipleEntitlements, _ := GetDataFromJsonString(jsonStrinMultipleEntitlements)

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
			data:          dataEntitlement,
			subaccountId:  "12345",
			expectedBlock: "import {\n\t\t\t\tto = btp_subaccount_entitlement.application-logs_lite\n\t\t\t\tid = \"12345,application-logs,lite\"\n\t\t\t  }\n\n",
			expectError:   false,
		},
		{
			name:          "With filter values",
			data:          dataMultipleEntitlements,
			subaccountId:  "12345",
			filterValues:  []string{"application-logs_lite"},
			expectedBlock: "import {\n\t\t\t\tto = btp_subaccount_entitlement.application-logs_lite\n\t\t\t\tid = \"12345,application-logs,lite\"\n\t\t\t  }\n\n",
			expectError:   false,
		},
		{
			name:          "Missing entitlement in filter values",
			data:          dataMultipleEntitlements,
			subaccountId:  "12345",
			filterValues:  []string{"non-existent-entitlement"},
			expectedBlock: "",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			importBlock, err := CreateEntitlementImportBlock(tt.data, tt.subaccountId, tt.filterValues, resourceDoc)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBlock, importBlock)
			}
		})
	}
}
