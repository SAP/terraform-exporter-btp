package tfimportprovider

import (
	"testing"

	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
	"github.com/stretchr/testify/assert"
)

func TestCreateCfUserImportBlock(t *testing.T) {
	resourceDoc := tfutils.EntityDocs{
		Import: "import {\n\t\t\t\tto =  cloudfoundry_user.<resource_name>\n\t\t\t\tid = \"<user_guid>\"\"\n\t\t\t  }\n",
	}

	jsonString := "{\"org\":\"12345\",\"space\": null,\"users\": [{\"annotations\": null,\"created_at\":\"2017-06-08T07:19:21Z\",\"id\":\"34567\",\"labels\": null,\"origin\":\"sap.ids\",\"presentation_name\":\"user1@sap.com\",\"updated_at\":\"2019-06-06T09:01:48Z\",\"username\":\"user1@sap.com\"}]}"
	dataUser, _ := GetDataFromJsonString(jsonString)

	jsonStringMultipleUsers := "{\"org\":\"12345\",\"space\": null,\"users\": [{\"annotations\": null,\"created_at\":\"2017-06-08T07:19:21Z\",\"id\":\"34567\",\"labels\": null,\"origin\":\"sap.ids\",\"presentation_name\":\"user1@sap.com\",\"updated_at\":\"2019-06-06T09:01:48Z\",\"username\":\"user1@sap.com\"},    {\"annotations\": null,\"created_at\":\"2017-06-08T05:27:17Z\",\"id\":\"56789\",\"labels\": null,\"origin\":\"sap.ids\",\"presentation_name\":\"user2@sap.com\",\"updated_at\":\"2017-06-08T05:27:17Z\",\"username\":\"user2@sap.com\"}]}"
	dataMultipleUsers, _ := GetDataFromJsonString(jsonStringMultipleUsers)

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
			data:          dataUser,
			orgId:         "12345",
			filterValues:  []string{},
			expectedBlock: "import {\n\t\t\t\tto =  cloudfoundry_user.user_0\n\t\t\t\tid = \"34567\"\"\n\t\t\t  }\n\n",
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "Valid data with matching filter",
			data:          dataMultipleUsers,
			orgId:         "12345",
			filterValues:  []string{"user1@sap.com"},
			expectedBlock: "import {\n\t\t\t\tto =  cloudfoundry_user.user_0\n\t\t\t\tid = \"34567\"\"\n\t\t\t  }\n\n",
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "Invalid filter value",
			data:          dataUser,
			orgId:         "12345",
			filterValues:  []string{"wrong-username"},
			expectedBlock: "",
			expectedCount: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			importBlock, count, err := createUserImportBlock(tt.data, tt.orgId, tt.filterValues, resourceDoc)
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
