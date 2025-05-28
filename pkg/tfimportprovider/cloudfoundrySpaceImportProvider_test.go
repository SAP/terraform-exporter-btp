package tfimportprovider

import (
	"testing"

	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
	"github.com/stretchr/testify/assert"
)

func TestCreateCfSpaceImportBlock(t *testing.T) {
	resourceDoc := tfutils.EntityDocs{
		Import: "import {\n\t\t\t\tto =  cloudfoundry_space.<resource_name>\n\t\t\t\tid = \"<space_guid>\"\"\n\t\t\t  }\n",
	}

	jsonString := "{\"name\": null,\"org\":\"12345\",\"spaces\": [{\"allow_ssh\": true,\"annotations\": null,\"created_at\":\"2024-11-04T11:57:47Z\",\"id\":\"9876543210\",\"isolation_segment\": null,\"labels\": null,\"name\":\"dev1\",\"org\":\"12345\",\"quota\": null,\"updated_at\":\"2024-11-04T11:57:47Z\"}]}"
	dataSpaceConfig, _ := GetDataFromJsonString(jsonString)

	jsonStringMultipleSpaces := "{\"name\": null,\"org\":\"12345\",\"spaces\": [{\"allow_ssh\": true,\"annotations\": null,\"created_at\":\"2024-11-04T11:57:47Z\",\"id\":\"9876543210\",\"isolation_segment\": null,\"labels\": null,\"name\":\"dev1\",\"org\":\"12345\",\"quota\": null,\"updated_at\":\"2024-11-04T11:57:47Z\"},    {\"allow_ssh\": true,\"annotations\": null,\"created_at\":\"2024-11-04T15:10:19Z\",\"id\":\"457812963\",\"isolation_segment\": null,\"labels\": null,\"name\":\"test\",\"org\":\"12345\",\"quota\": null,\"updated_at\":\"2024-11-04T15:10:19Z\"}]}"
	dataMultipleSpaces, _ := GetDataFromJsonString(jsonStringMultipleSpaces)

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
			data:          dataSpaceConfig,
			orgId:         "12345",
			filterValues:  []string{},
			expectedBlock: "import {\n\t\t\t\tto =  cloudfoundry_space.space_0\n\t\t\t\tid = \"9876543210\"\"\n\t\t\t  }\n\n",
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "Valid data with matching filter",
			data:          dataMultipleSpaces,
			orgId:         "12345",
			filterValues:  []string{"dev1"},
			expectedBlock: "import {\n\t\t\t\tto =  cloudfoundry_space.space_0\n\t\t\t\tid = \"9876543210\"\"\n\t\t\t  }\n\n",
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "Invalid filter value",
			data:          dataSpaceConfig,
			orgId:         "12345",
			filterValues:  []string{"wrong-space-name"},
			expectedBlock: "",
			expectedCount: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			importBlock, count, err := createSpaceImportBlock(tt.data, tt.orgId, tt.filterValues, resourceDoc)
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
