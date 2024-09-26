package tfimportprovider

import (
	"testing"

	tfutils "github.com/SAP/terraform-exporter-btp/tfutils"
	"github.com/stretchr/testify/assert"
)

func TestCreateEnvironmentInstanceImportBlock(t *testing.T) {
	resourceDoc := tfutils.EntityDocs{
		Import: "import {\n\t\t\t\tto = btp_subaccount_environment_instance.<resource_name>\n\t\t\t\tid = \"<subaccount_id>,<environment_instance_id>\"\n\t\t\t  }\n",
	}

	jsonString := "{\"id\":\"12345\",\"subaccount_id\":\"12345\",\"values\":[{\"broker_id\":\"E181A0CD-B424-4D07-87DA-0767EF912337\",\"created_date\":\"2024-09-16T08:06:16Z\",\"dashboard_url\":\"\",\"description\":\"\",\"environment_type\":\"cloudfoundry\",\"id\":\"8D363D57-86F1-4FF4-B3F8-C1FB7FFD34AC\",\"landscape_label\":\"cf-us10\",\"last_modified\":\"2024-09-16T08:06:20Z\",\"name\":\"btptfexporter-validate-xhyu0cme_btptfexporter-validate\",\"operation\":\"provision\",\"plan_id\":\"fc5abe63-2a7d-4848-babf-f63a5d316df1\",\"plan_name\":\"standard\",\"platform_id\":\"987654\",\"service_id\":\"fa31b750-375f-4268-bee1-604811a89fd9\",\"service_name\":\"cloudfoundry\",\"state\":\"OK\",\"tenant_id\":\"12345\",\"type\":\"Provision\"}]}"

	dataEnvInstance, _ := GetDataFromJsonString(jsonString)

	jsonStringMultiple := "{\"id\":\"12345\",\"subaccount_id\":\"12345\",\"values\":[{\"broker_id\":\"E181A0CD-B424-4D07-87DA-0767EF912337\",\"created_date\":\"2024-09-16T08:06:16Z\",\"dashboard_url\":\"\",\"description\":\"\",\"environment_type\":\"cloudfoundry\",\"id\":\"8D363D57-86F1-4FF4-B3F8-C1FB7FFD34AC\",\"landscape_label\":\"cf-us10\",\"last_modified\":\"2024-09-16T08:06:20Z\",\"name\":\"btptfexporter-validate-xhyu0cme_btptfexporter-validate\",\"operation\":\"provision\",\"plan_id\":\"fc5abe63-2a7d-4848-babf-f63a5d316df1\",\"plan_name\":\"standard\",\"platform_id\":\"987654\",\"service_id\":\"fa31b750-375f-4268-bee1-604811a89fd9\",\"service_name\":\"cloudfoundry\",\"state\":\"OK\",\"tenant_id\":\"12345\",\"type\":\"Provision\"},{\"broker_id\":\"E181A0CD-B424-4D07-87DA-0767EF912337\",\"created_date\":\"2024-09-16T08:06:16Z\",\"dashboard_url\":\"\",\"description\":\"\",\"environment_type\":\"cloudfoundry2\",\"id\":\"00000000-0000-0000-0000-000000000000\",\"landscape_label\":\"cf-us10\",\"last_modified\":\"2024-09-16T08:06:20Z\",\"name\":\"btptfexporter-validate-xhyu0cme_btptfexporter-validate\",\"operation\":\"provision\",\"plan_id\":\"fc5abe63-2a7d-4848-babf-f63a5d316df1\",\"plan_name\":\"standard\",\"platform_id\":\"987654\",\"service_id\":\"fa31b750-375f-4268-bee1-604811a89fd9\",\"service_name\":\"cloudfoundry\",\"state\":\"OK\",\"tenant_id\":\"12345\",\"type\":\"Provision\"}]}"
	dataEnvInstanceMultiple, _ := GetDataFromJsonString(jsonStringMultiple)

	tests := []struct {
		name          string
		data          map[string]interface{}
		subaccountId  string
		filterValues  []string
		expectedBlock string
		expectError   bool
	}{
		{
			name:          "No filter values",
			data:          dataEnvInstance,
			subaccountId:  "12345",
			expectedBlock: "import {\n\t\t\t\tto = btp_subaccount_environment_instance.cloudfoundry\n\t\t\t\tid = \"12345,8D363D57-86F1-4FF4-B3F8-C1FB7FFD34AC\"\n\t\t\t  }\n\n",
			expectError:   false,
		},
		{
			name:          "With filter values",
			data:          dataEnvInstanceMultiple,
			subaccountId:  "12345",
			filterValues:  []string{"cloudfoundry"},
			expectedBlock: "import {\n\t\t\t\tto = btp_subaccount_environment_instance.cloudfoundry\n\t\t\t\tid = \"12345,8D363D57-86F1-4FF4-B3F8-C1FB7FFD34AC\"\n\t\t\t  }\n\n",
			expectError:   false,
		},
		{
			name:          "Filter value not found",
			data:          dataEnvInstance,
			subaccountId:  "12345",
			filterValues:  []string{"kubernetes"},
			expectedBlock: "",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			importBlock, err := createEnvironmentInstanceImportBlock(tt.data, tt.subaccountId, tt.filterValues, resourceDoc)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBlock, importBlock)
			}
		})
	}
}
