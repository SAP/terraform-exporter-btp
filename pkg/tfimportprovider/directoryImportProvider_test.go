package tfimportprovider

import (
	"testing"

	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
	"github.com/stretchr/testify/assert"
)

func TestCreateDirectoryImportBlock(t *testing.T) {
	resourceDoc := tfutils.EntityDocs{
		Import: "import {\n\t\t\t\tto = btp_directory.<resource_name>\n\t\t\t\tid = \"<directory_id>\"\"\n\t\t\t  }\n",
	}

	jsonString := "{\"created_by\":\"some.body@sap.com\",\"created_date\":\"2024-10-09T12:05:16Z\",\"description\":\"\",\"features\": [\"AUTHORIZATIONS\",\"DEFAULT\",\"ENTITLEMENTS\"],\"id\":\"12345\",\"labels\": {\"test\": [\"exporter\"]},\"last_modified\":\"2024-10-09T12:05:32Z\",\"name\":\"btptfexporter-validate-dir\",\"parent_id\":\"67890\",\"state\":\"OK\",\"subdomain\":\"12345\"}"

	dataDirectory, _ := GetDataFromJsonString(jsonString)

	tests := []struct {
		name          string
		data          map[string]interface{}
		directoryId   string
		filterValues  []string
		expectedBlock string
		expectError   bool
	}{

		{
			name:          "Valid data without filter",
			data:          dataDirectory,
			directoryId:   "12345",
			filterValues:  []string{},
			expectedBlock: "import {\n\t\t\t\tto = btp_directory.directory_0\n\t\t\t\tid = \"12345\"\"\n\t\t\t  }\n\n",
			expectError:   false,
		},
		{
			name:          "Valid data with matching filter",
			data:          dataDirectory,
			directoryId:   "12345",
			filterValues:  []string{"btptfexporter-validate-dir"},
			expectedBlock: "import {\n\t\t\t\tto = btp_directory.directory_0\n\t\t\t\tid = \"12345\"\"\n\t\t\t  }\n\n",
			expectError:   false,
		},
		{
			name:          "Invalid filter value",
			data:          dataDirectory,
			directoryId:   "12345",
			filterValues:  []string{"wrong-directory"},
			expectedBlock: "",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			importBlock, err := createDirectoryImportBlock(tt.data, tt.directoryId, tt.filterValues, resourceDoc)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBlock, importBlock)
			}
		})
	}
}
