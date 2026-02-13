package tfimportprovider

import (
	"testing"

	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
	"github.com/stretchr/testify/assert"
)

func TestCreateSpaceQuotaImportBlock(t *testing.T) {
	resourceDoc := tfutils.EntityDocs{
		Import: "import {\n\t\t\t\tto =  cloudfoundry_space_quota.<resource_name>\n\t\t\t\tid = \"<space_quota_guid>\"\"\n\t\t\t  }\n",
	}

	jsonString := "{\"name\": null,\"org\":\"12345\",\"space_quotas\": [{\"allow_paid_service_plans\": false,\"created_at\":\"2024-12-03T07:03:20Z\",\"id\":\"12345678\",\"instance_memory\": 15000,\"name\":\"test-quota1\",\"org\":\"12345\",\"spaces\": null,\"total_app_instances\": 10,\"total_app_log_rate_limit\": null,\"total_app_tasks\": null,\"total_memory\": 2048,\"total_route_ports\": null,\"total_routes\": 10,\"total_service_keys\": null,\"total_services\": 100,\"updated_at\":\"2024-12-03T07:03:20Z\"}]}"
	dataSpaceQuotaConfig, _ := GetDataFromJsonString(jsonString)

	jsonStringMultipleSpacesQuotas := "{\"name\": null,\"org\":\"12345\",\"space_quotas\": [{\"allow_paid_service_plans\": false,\"created_at\":\"2024-12-03T07:03:20Z\",\"id\":\"12345678\",\"instance_memory\": 15000,\"name\":\"test-quota1\",\"org\":\"12345\",\"spaces\": null,\"total_app_instances\": 10,\"total_app_log_rate_limit\": null,\"total_app_tasks\": null,\"total_memory\": 2048,\"total_route_ports\": null,\"total_routes\": 10,\"total_service_keys\": null,\"total_services\": 100,\"updated_at\":\"2024-12-03T07:03:20Z\"},    {\"allow_paid_service_plans\": false,\"created_at\":\"2024-12-03T07:03:20Z\",\"id\":\"23456789\",\"instance_memory\": 15000,\"name\":\"test-quota2\",\"org\":\"12345\",\"spaces\": null,\"total_app_instances\": 10,\"total_app_log_rate_limit\": null,\"total_app_tasks\": null,\"total_memory\": 2048,\"total_route_ports\": null,\"total_routes\": 10,\"total_service_keys\": null,\"total_services\": 100,\"updated_at\":\"2024-12-03T07:03:20Z\"}]}"
	dataMultipleSpacesQuotas, _ := GetDataFromJsonString(jsonStringMultipleSpacesQuotas)

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
			data:          dataSpaceQuotaConfig,
			orgId:         "12345",
			filterValues:  []string{},
			expectedBlock: "import {\n\t\t\t\tto =  cloudfoundry_space_quota.space_quota_0\n\t\t\t\tid = \"12345678\"\"\n\t\t\t  }\n\n",
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "Valid data with matching filter",
			data:          dataMultipleSpacesQuotas,
			orgId:         "12345",
			filterValues:  []string{"test-quota2"},
			expectedBlock: "import {\n\t\t\t\tto =  cloudfoundry_space_quota.space_quota_1\n\t\t\t\tid = \"23456789\"\"\n\t\t\t  }\n\n",
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "Invalid filter value",
			data:          dataSpaceQuotaConfig,
			orgId:         "12345",
			filterValues:  []string{"wrong-space-quota-name"},
			expectedBlock: "",
			expectedCount: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			importBlock, count, err := createSpaceQuotaImportBlock(tt.data, tt.orgId, tt.filterValues, resourceDoc)
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
