package tfimportprovider

import (
	"testing"

	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
	"github.com/stretchr/testify/assert"
)

func TestCreateCfServiceInstanceImportBlock(t *testing.T) {
	resourceDoc := tfutils.EntityDocs{
		Import: "import {\n\t\t\t\tto =  cloudfoundry_service_instance.<resource_name>\n\t\t\t\tid = \"<service_instance_guid>\"\"\n\t\t\t  }\n",
	}

	jsonString := "{\"name\": null,\"org\":\"12345\",\"service_instances\": [{\"id\":\"56789\",\"name\":\"test-instance-1\",\"route_service_url\": null,\"service_plan\":\"23456\",\"space\":\"67895\",\"syslog_drain_url\": null,\"tags\": null,\"type\":\"managed\",\"updated_at\":\"2024-01-17T13:20:17Z\",\"upgrade_available\": false}],\"space\": null}"
	dataServiceInstanceConfig, _ := GetDataFromJsonString(jsonString)

	jsonStringMultipleServiceInstances := "{\"name\": null,\"org\":\"12345\",\"service_instances\": [{\"id\":\"56789\",\"name\":\"test-instance-1\",\"route_service_url\": null,\"service_plan\":\"23456\",\"space\":\"67895\",\"syslog_drain_url\": null,\"tags\": null,\"type\":\"managed\",\"updated_at\":\"2024-01-17T13:20:17Z\",\"upgrade_available\": false},    {\"id\":\"45678\",\"name\":\"test-instance-2\",\"route_service_url\": null,\"service_plan\":\"34567\",\"space\":\"67895\",\"syslog_drain_url\": null,\"tags\": null,\"type\":\"managed\",\"updated_at\":\"2024-01-30T08:14:19Z\",\"upgrade_available\": false}],\"space\": null}"
	dataMultipleServiceinstances, _ := GetDataFromJsonString(jsonStringMultipleServiceInstances)

	tests := []struct {
		name          string
		data          map[string]interface{}
		orgId         string
		filterValues  []string
		expectedBlock string
		expectedCount int
		expectError   bool
	}{

		{
			name:          "Valid data without filter",
			data:          dataServiceInstanceConfig,
			orgId:         "12345",
			filterValues:  []string{},
			expectedBlock: "import {\n\t\t\t\tto =  cloudfoundry_service_instance.cf_serviceinstance_0\n\t\t\t\tid = \"56789\"\"\n\t\t\t  }\n\n",
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "Valid data with matching filter",
			data:          dataMultipleServiceinstances,
			orgId:         "12345",
			filterValues:  []string{"test-instance-2_34567"},
			expectedBlock: "import {\n\t\t\t\tto =  cloudfoundry_service_instance.cf_serviceinstance_1\n\t\t\t\tid = \"45678\"\"\n\t\t\t  }\n\n",
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "Invalid filter value",
			data:          dataServiceInstanceConfig,
			orgId:         "12345",
			filterValues:  []string{"wrong-instance-name"},
			expectedBlock: "",
			expectedCount: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			importBlock, count, err := createCfServiceInstanceImportBlock(tt.data, tt.orgId, tt.filterValues, resourceDoc)
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
