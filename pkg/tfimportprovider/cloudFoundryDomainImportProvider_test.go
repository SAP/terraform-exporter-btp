package tfimportprovider

import (
	"testing"

	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
	"github.com/stretchr/testify/assert"
)

func TestCreateCfDomainImportBlock(t *testing.T) {
	resourceDoc := tfutils.EntityDocs{
		Import: "import {\n\t\t\t\tto =  cloudfoundry_domain.<resource_name>\n\t\t\t\tid = \"<domain_guid>\"\"\n\t\t\t  }\n",
	}

	jsonString := "{\"domains\": [{\"annotations\": null,\"created_at\":\"2022-01-02T08:56:47Z\",\"id\":\"23456\",\"internal\": false,\"labels\": null,\"name\":\"test-domain1.com\",\"org\":\"12345\",\"router_group\": null,\"shared_orgs\": null,\"supported_protocols\": [\"http\"],\"updated_at\":\"2023-01-02T08:56:47Z\"}],\"org\":\"12345\"}"
	dataDomain, _ := GetDataFromJsonString(jsonString)

	jsonStringMultipleDomains := "{\"domains\": [{\"annotations\": null,\"created_at\":\"2022-01-02T08:56:47Z\",\"id\":\"23456\",\"internal\": false,\"labels\": null,\"name\":\"test-domain1.com\",\"org\":\"12345\",\"router_group\": null,\"shared_orgs\": null,\"supported_protocols\": [\"http\"],\"updated_at\":\"2023-01-02T08:56:47Z\"},    {\"annotations\": null,\"created_at\":\"2022-01-02T09:29:59Z\",\"id\":\"34567\",\"internal\": false,\"labels\": null,\"name\":\"test-domain2.com\",\"org\":\"12345\",\"router_group\": null,\"shared_orgs\": null,\"supported_protocols\": [\"http\"],\"updated_at\":\"2023-01-02T09:29:59Z\"}],\"org\":\"12345\"}"
	dataMultipleDomains, _ := GetDataFromJsonString(jsonStringMultipleDomains)

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
			data:          dataDomain,
			orgId:         "12345",
			filterValues:  []string{},
			expectedBlock: "import {\n\t\t\t\tto =  cloudfoundry_domain.domain_0\n\t\t\t\tid = \"23456\"\"\n\t\t\t  }\n\n",
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "Valid data with matching filter",
			data:          dataMultipleDomains,
			orgId:         "12345",
			filterValues:  []string{"test-domain1.com"},
			expectedBlock: "import {\n\t\t\t\tto =  cloudfoundry_domain.domain_0\n\t\t\t\tid = \"23456\"\"\n\t\t\t  }\n\n",
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "Invalid filter value",
			data:          dataDomain,
			orgId:         "12345",
			filterValues:  []string{"wrong-domain-name"},
			expectedBlock: "",
			expectedCount: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			importBlock, count, err := createDomainImportBlock(tt.data, tt.orgId, tt.filterValues, resourceDoc)
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
