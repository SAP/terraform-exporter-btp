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

	jsonString := "{\"auditlog-management:default\": {\"category\":\"ELASTIC_SERVICE\",\"plan_description\":\"Free offering for development purposes\",\"plan_display_name\":\"lite\",\"plan_name\":\"default\",\"quota_assigned\": 1,\"quota_remaining\": 1,\"service_display_name\":\"auditlog-management\",\"service_name\":\"auditlog-management\"}}"
	dataEntitlement, _ := GetDataFromJsonString(jsonString)

	jsonStrinMultipleEntitlements := "{\"auditlog-management:default\": {\"category\":\"ELASTIC_SERVICE\",\"plan_description\":\"Free offering for development purposes\",\"plan_display_name\":\"default\",\"plan_name\":\"default\",\"quota_assigned\": 1,\"quota_remaining\": 1,\"service_display_name\":\"auditlog-management\",\"service_name\":\"auditlog-management\"},\"auditlog-api:default\": {\"category\":\"ELASTIC_SERVICE\",\"plan_description\":\"Default plan for Auditlog API\",\"plan_display_name\":\"Default\",\"plan_name\":\"default\",\"quota_assigned\": 1,\"quota_remaining\": 1,\"service_display_name\":\"[DEPRECATED] Audit Log Retrieval\",\"service_name\":\"auditlog-api\"}}"
	dataMultipleEntitlements, _ := GetDataFromJsonString(jsonStrinMultipleEntitlements)

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
			data:          dataEntitlement,
			subaccountId:  "12345",
			expectedBlock: "import {\n\t\t\t\tto = btp_subaccount_entitlement.entitlement_0\n\t\t\t\tid = \"12345,auditlog-management,default\"\n\t\t\t  }\n\n",
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "With filter values",
			data:          dataMultipleEntitlements,
			subaccountId:  "12345",
			filterValues:  []string{"auditlog-management_default"},
			expectedBlock: "import {\n\t\t\t\tto = btp_subaccount_entitlement.entitlement_0\n\t\t\t\tid = \"12345,auditlog-management,default\"\n\t\t\t  }\n\n",
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "Missing entitlement in filter values",
			data:          dataMultipleEntitlements,
			subaccountId:  "12345",
			filterValues:  []string{"non-existent-entitlement"},
			expectedBlock: "",
			expectedCount: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			importBlock, count, err := CreateEntitlementImportBlock(tt.data, tt.subaccountId, tt.filterValues, resourceDoc)
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
