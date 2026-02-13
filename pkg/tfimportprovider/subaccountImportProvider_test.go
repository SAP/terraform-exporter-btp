package tfimportprovider

import (
	"testing"

	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
	"github.com/stretchr/testify/assert"
)

func TestCreateSubaccountImportBlock(t *testing.T) {
	resourceDoc := tfutils.EntityDocs{
		Import: "import {\n\t\t\t\tto = btp_subaccount.<resource_name>\n\t\t\t\tid = \"<subaccount_id>\"\"\n\t\t\t  }\n",
	}

	jsonString := "{\"beta_enabled\": false,\"created_by\":\"someone@sap.com\",\"created_date\":\"2024-09-11T10:32:19Z\",\"description\":\"\",\"id\":\"12345\",\"labels\": null,\"last_modified\":\"2024-09-11T10:32:42Z\",\"name\":\"BTP tfexporter\",\"parent_features\": null,\"parent_id\":\"543121\",\"region\":\"us10\",\"state\":\"OK\",\"subdomain\":\"btptfexporter\",\"usage\":\"NOT_USED_FOR_PRODUCTION\"}"

	dataSubaccount, _ := GetDataFromJsonString(jsonString)

	tests := []struct {
		name          string
		data          map[string]any
		subaccountId  string
		filterValues  []string
		expectedBlock string
		expectError   bool
	}{

		{
			name:          "Valid data without filter",
			data:          dataSubaccount,
			subaccountId:  "12345",
			filterValues:  []string{},
			expectedBlock: "import {\n\t\t\t\tto = btp_subaccount.subaccount_0\n\t\t\t\tid = \"12345\"\"\n\t\t\t  }\n\n",
			expectError:   false,
		},
		{
			name:          "Valid data with matching filter",
			data:          dataSubaccount,
			subaccountId:  "12345",
			filterValues:  []string{"BTP tfexporter"},
			expectedBlock: "import {\n\t\t\t\tto = btp_subaccount.subaccount_0\n\t\t\t\tid = \"12345\"\"\n\t\t\t  }\n\n",
			expectError:   false,
		},
		{
			name:          "Invalid filter value",
			data:          dataSubaccount,
			subaccountId:  "12345",
			filterValues:  []string{"wrong-subaccount"},
			expectedBlock: "",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			importBlock, err := createSubaccountImportBlock(tt.data, tt.subaccountId, tt.filterValues, resourceDoc)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBlock, importBlock)
			}
		})
	}
}
