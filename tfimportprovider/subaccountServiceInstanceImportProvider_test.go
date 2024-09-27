package tfimportprovider

import (
	"testing"

	tfutils "github.com/SAP/terraform-exporter-btp/tfutils"
	"github.com/stretchr/testify/assert"
)

func TestCreateServiceInstanceImportBlock(t *testing.T) {
	resourceDoc := tfutils.EntityDocs{
		Import: "import {\n\t\t\t\tto = btp_subaccount_service_instance.<resource_name>\n\t\t\t\tid = \"<subaccount_id>,<service_instance_id>\"\"\n\t\t\t  }\n",
	}

	jsonString := ""
	dataServiceInstance, _ := GetDataFromJsonString(jsonString)

	jsonStringMultipleInstances := ""
	dataMultipleServiceInstances, _ := GetDataFromJsonString(jsonStringMultipleInstances)

	tests := []struct {
		name          string
		data          map[string]interface{}
		subaccountId  string
		filterValues  []string
		expectedBlock string
		expectError   bool
	}{

		{
			name:          "Valid data without filter",
			data:          dataServiceInstance,
			subaccountId:  "12345",
			filterValues:  []string{},
			expectedBlock: "import {\n\t\t\t\tto = btp_subaccount_service_instance.\n\t\t\t\tid = \"12345,\"\"\n\t\t\t  }\n\n",
			expectError:   false,
		},
		{
			name:          "Valid data with matching filter",
			data:          dataMultipleServiceInstances,
			subaccountId:  "12345",
			filterValues:  []string{"connection_instance_12345"},
			expectedBlock: "import {\n\t\t\t\tto = btp_subaccount_service_instance.\n\t\t\t\tid = \"12345,\"\"\n\t\t\t  }\n\n",
			expectError:   false,
		},
		{
			name:          "Invalid filter value",
			data:          dataServiceInstance,
			subaccountId:  "12345",
			filterValues:  []string{"wrong-instance-name"},
			expectedBlock: "",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			importBlock, err := createServiceInstanceImportBlock(tt.data, tt.subaccountId, tt.filterValues, resourceDoc)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBlock, importBlock)
			}
		})
	}
}
