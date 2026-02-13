package tfimportprovider

import (
	"testing"

	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
	"github.com/stretchr/testify/assert"
)

func TestCreateServiceInstanceImportBlock(t *testing.T) {
	resourceDoc := tfutils.EntityDocs{
		Import: "import {\n\t\t\t\tto = btp_subaccount_service_instance.<resource_name>\n\t\t\t\tid = \"<subaccount_id>,<service_instance_id>\"\"\n\t\t\t  }\n",
	}

	jsonString := "{\"fields_filter\": null,\"id\":\"12345\",\"labels_filter\": null,\"subaccount_id\":\"12345\",\"values\": [{\"context\":\"{\\\"platform\\\":\\\"sapcp\\\",\\\"origin\\\":\\\"sapcp\\\"}\",\"created_date\":\"2024-09-27T11:55:09Z\",\"dashboard_url\":\"\",\"id\":\"0b031c77-a995-44cf-b8ce-c7369ad2cd0f\",\"last_modified\":\"2024-09-27T11:55:10Z\",\"name\":\"audit-log-exporter\",\"platform_id\":\"service-manager\",\"ready\": true,\"serviceplan_id\":\"a50128a9-35fc-4624-9953-c79668ef3e5b\",\"usable\": true}]}"
	dataServiceInstance, _ := GetDataFromJsonString(jsonString)

	jsonStringWithCfInstance := "{\"fields_filter\": null,\"id\":\"12345\",\"labels_filter\": null,\"subaccount_id\":\"12345\",\"values\": [{\"context\":\"{\\\"platform\\\":\\\"cloudfoundry\\\",\\\"origin\\\":\\\"cloudfoundry\\\"}\",\"created_date\":\"2024-09-27T11:55:09Z\",\"dashboard_url\":\"\",\"id\":\"0b031c77-a995-44cf-b8ce-c7369ad2cd0f\",\"last_modified\":\"2024-09-27T11:55:10Z\",\"name\":\"audit-log-exporter\",\"platform_id\":\"service-manager\",\"ready\": true,\"serviceplan_id\":\"a50128a9-35fc-4624-9953-c79668ef3e5b\",\"usable\": true}]}"
	dataCfServiceInstance, _ := GetDataFromJsonString(jsonStringWithCfInstance)

	jsonStringMultipleInstances := "{\"fields_filter\": null,\"id\":\"12345\",\"labels_filter\": null,\"subaccount_id\":\"12345\",\"values\": [{\"context\":\"{\\\"platform\\\":\\\"sapcp\\\",\\\"origin\\\":\\\"sapcp\\\"}\",\"created_date\":\"2024-09-25T12:15:03Z\",\"dashboard_url\":\"\",\"id\":\"8c91e9ec-75cb-40d3-aec6-c36dc4582853\",\"last_modified\":\"2024-09-25T12:15:10Z\",\"name\":\"content-agent-test\",\"platform_id\":\"service-manager\",\"ready\": true,\"serviceplan_id\":\"e39c1e0a-838b-4b62-ab2a-3b23b58e5baf\",\"usable\": true},{\"context\":\"{\\\"platform\\\":\\\"sapcp\\\",\\\"origin\\\":\\\"sapcp\\\"}]\",\"created_date\":\"2024-09-27T11:55:09Z\",\"dashboard_url\":\"\",\"id\":\"0b031c77-a995-44cf-b8ce-c7369ad2cd0f\",\"last_modified\":\"2024-09-27T11:55:10Z\",\"name\":\"audit-log-exporter\",\"platform_id\":\"service-manager\",\"ready\": true,\"serviceplan_id\":\"a50128a9-35fc-4624-9953-c79668ef3e5b\",\"usable\": true}]}"
	dataMultipleServiceInstances, _ := GetDataFromJsonString(jsonStringMultipleInstances)

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
			name:          "Valid data without filter",
			data:          dataServiceInstance,
			subaccountId:  "12345",
			filterValues:  []string{},
			expectedBlock: "import {\n\t\t\t\tto = btp_subaccount_service_instance.serviceinstance_0\n\t\t\t\tid = \"12345,0b031c77-a995-44cf-b8ce-c7369ad2cd0f\"\"\n\t\t\t  }\n\n",
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "Instance with origin as cloudfoundry not sapcp",
			data:          dataCfServiceInstance,
			subaccountId:  "12345",
			filterValues:  []string{},
			expectedBlock: "",
			expectedCount: 0,
			expectError:   false,
		},
		{
			name:          "Valid data with matching filter",
			data:          dataMultipleServiceInstances,
			subaccountId:  "12345",
			filterValues:  []string{"audit-log-exporter_a50128a9-35fc-4624-9953-c79668ef3e5b"},
			expectedBlock: "import {\n\t\t\t\tto = btp_subaccount_service_instance.serviceinstance_1\n\t\t\t\tid = \"12345,0b031c77-a995-44cf-b8ce-c7369ad2cd0f\"\"\n\t\t\t  }\n\n",
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "Invalid filter value",
			data:          dataServiceInstance,
			subaccountId:  "12345",
			filterValues:  []string{"wrong-instance-name"},
			expectedBlock: "",
			expectedCount: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			importBlock, count, err := createServiceInstanceImportBlock(tt.data, tt.subaccountId, tt.filterValues, resourceDoc)
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
