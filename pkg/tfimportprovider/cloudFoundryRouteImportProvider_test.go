package tfimportprovider

import (
	"testing"

	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
	"github.com/stretchr/testify/assert"
)

func TestCreateCfRouteImportBlock(t *testing.T) {
	resourceDoc := tfutils.EntityDocs{
		Import: "import {\n\t\t\t\tto =  cloudfoundry_route.<resource_name>\n\t\t\t\tid = \"<route_guid>\"\"\n\t\t\t  }\n",
	}

	jsonString := "{\"org\": 12345,\"routes\": [{\"annotations\": null,\"created_at\":\"2024-01-17T13:29:02Z\",\"destinations\": [{\"app_id\":\"34567\",\"app_process_type\":\"web\",\"id\":\"23456\",\"port\": 8080,\"protocol\":\"http1\",\"weight\": null}],\"domain\":\"45678\",\"host\":\"test-host1\",\"id\":\"23456\",\"labels\": null,\"path\": null,\"port\": null,\"protocol\":\"http\",\"space\":\"56789\",\"updated_at\":\"2024-01-17T13:29:02Z\",\"url\":\"test-host1.eu12.hana.ondemand.com\"}],\"space\": null}"
	dataRoute, _ := GetDataFromJsonString(jsonString)

	jsonStringMultipleRoutes := "{\"org\": 12345,\"routes\": [{\"annotations\": null,\"created_at\":\"2024-01-17T13:29:02Z\",\"destinations\": [{\"app_id\":\"34567\",\"app_process_type\":\"web\",\"id\":\"23456\",\"port\": 8080,\"protocol\":\"http1\",\"weight\": null}],\"domain\":\"45678\",\"host\":\"test-host\",\"id\":\"23456\",\"labels\": null,\"path\": null,\"port\": null,\"protocol\":\"http\",\"space\":\"56789\",\"updated_at\":\"2024-01-17T13:29:02Z\",\"url\":\"test-host.eu12.hana.ondemand.com\"},    {\"annotations\": null,\"created_at\":\"2024-01-17T13:29:02Z\",\"destinations\": [{\"app_id\":\"35672\",\"app_process_type\":\"web\",\"id\":\"67895\",\"port\": 8080,\"protocol\":\"http1\",\"weight\": null}],\"domain\":\"45678\",\"host\":\"test-host2\",\"id\":\"67895\",\"labels\": null,\"path\": null,\"port\": null,\"protocol\":\"http\",\"space\":\"56789\",\"updated_at\":\"2024-01-17T13:29:02Z\",\"url\":\"test-host2.eu12.hana.ondemand.com\"}],\"space\": null}"
	dataMultipleRoutes, _ := GetDataFromJsonString(jsonStringMultipleRoutes)

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
			data:          dataRoute,
			orgId:         "12345",
			filterValues:  []string{},
			expectedBlock: "import {\n\t\t\t\tto =  cloudfoundry_route.route_0\n\t\t\t\tid = \"23456\"\"\n\t\t\t  }\n\n",
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "Valid data with matching filter",
			data:          dataMultipleRoutes,
			orgId:         "12345",
			filterValues:  []string{"test-host2.eu12.hana.ondemand.com"},
			expectedBlock: "import {\n\t\t\t\tto =  cloudfoundry_route.route_1\n\t\t\t\tid = \"67895\"\"\n\t\t\t  }\n\n",
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "Invalid filter value",
			data:          dataRoute,
			orgId:         "12345",
			filterValues:  []string{"wrong-route-url"},
			expectedBlock: "",
			expectedCount: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			importBlock, count, err := createRouteImportBlock(tt.data, tt.orgId, tt.filterValues, resourceDoc)
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
