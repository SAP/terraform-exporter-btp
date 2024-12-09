package tfimportprovider

import (
	"testing"

	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
	"github.com/stretchr/testify/assert"
)

func TestCreateServiceBindingImportBlock(t *testing.T) {
	resourceDoc := tfutils.EntityDocs{
		Import: "import {\n\t\t\t\tto = btp_subaccount_service_binding.<resource_name>\n\t\t\t\tid = \"<subaccount_id>,<service_binding_id>\"\"\n\t\t\t  }\n",
	}

	jsonString := "{\"fields_filter\": null,\"id\":\"12345\",\"labels_filter\": null,\"subaccount_id\":\"12345\",\"values\": [{\"bind_resource\": null,\"context\":\"{\\\"crm_customer_id\\\":\\\"\\\",\\\"env_type\\\":\\\"sapcp\\\",\\\"global_account_id\\\":\\\"56789\\\",\\\"instance_name\\\":\\\"test\\\",\\\"license_type\\\":\\\"DEVELOPER\\\",\\\"origin\\\":\\\"sapcp\\\",\\\"platform\\\":\\\"sapcp\\\",\\\"region\\\":\\\"cf-eu10\\\",\\\"service_instance_id\\\":\\\"123678\\\",\\\"subaccount_id\\\":\\\"12345\\\",\\\"subdomain\\\":\\\"testcf-5vpfkrrj\\\",\\\"zone_id\\\":\\\"234567\\\"}\",\"created_date\":\"2024-09-25T05:51:08Z\",\"credentials\":\"{\\\"serviceInstanceId\\\":\\\"123678\\\",\\\"systemId\\\":\\\"45678\\\",\\\"uaa\\\":{\\\"apiurl\\\":\\\"https://api.authentication.eu10.hana.ondemand.com\\\",\\\"clientid\\\":\\\"sb-12345|one-mds-master!b28283\\\",\\\"clientsecret\\\":\\\"5dfa16a0-12345\\\",\\\"credential-type\\\":\\\"binding-secret\\\",\\\"identityzone\\\":\\\"testcf-5vpfkrrj\\\",\\\"identityzoneid\\\":\\\"12345\\\",\\\"sburl\\\":\\\"https://internal-xsuaa.authentication.eu10.hana.ondemand.com\\\",\\\"serviceInstanceId\\\":\\\"123678\\\",\\\"subaccountid\\\":\\\"12345\\\",\\\"tenantid\\\":\\\"12345\\\",\\\"tenantmode\\\":\\\"dedicated\\\",\\\"uaadomain\\\":\\\"authentication.eu10.hana.ondemand.com\\\",\\\"url\\\":\\\"https://testcf-5vpfkrrj.authentication.eu10.hana.ondemand.com\\\",\\\"verificationkey\\\":\\\"-----BEGIN PUBLIC KEY-----\\n12abcdef34\\n12abcd34ef\\n12abcd34ef\\n12abcd34ef\\nabcdef\\n-----END PUBLIC KEY-----\\\",\\\"xsappname\\\":\\\"d12345678|one-mds-master!b28283\\\",\\\"zoneid\\\":\\\"12345\\\"},\\\"uri\\\":\\\"https://one-mds.cfapps.eu10.hana.ondemand.com\\\"}\",\"id\":\"56789\",\"labels\": {\"subaccount_id\": [\"12345\\\"\"]},\"last_modified\":\"2024-09-25T05:51:09Z\",\"name\":\"my binding\",\"ready\": true,\"service_instance_id\":\"123678\"}]}"
	dataServiceBinding, _ := GetDataFromJsonString(jsonString)

	jsonStringMultipleBindings := "{\"fields_filter\": null,\"id\":\"12345\",\"labels_filter\": null,\"subaccount_id\":\"12345\",\"values\": [{\"bind_resource\": null,\"context\":\"{\\\"crm_customer_id\\\":\\\"\\\",\\\"env_type\\\":\\\"sapcp\\\",\\\"global_account_id\\\":\\\"56789\\\",\\\"instance_name\\\":\\\"test\\\",\\\"license_type\\\":\\\"DEVELOPER\\\",\\\"origin\\\":\\\"sapcp\\\",\\\"platform\\\":\\\"sapcp\\\",\\\"region\\\":\\\"cf-eu10\\\",\\\"service_instance_id\\\":\\\"123678\\\",\\\"subaccount_id\\\":\\\"12345\\\",\\\"subdomain\\\":\\\"testcf-5vpfkrrj\\\",\\\"zone_id\\\":\\\"234567\\\"}\",\"created_date\":\"2024-09-25T05:51:08Z\",\"credentials\":\"{\\\"serviceInstanceId\\\":\\\"123678\\\",\\\"systemId\\\":\\\"45678\\\",\\\"uaa\\\":{\\\"apiurl\\\":\\\"https://api.authentication.eu10.hana.ondemand.com\\\",\\\"clientid\\\":\\\"sb-12345|one-mds-master!b28283\\\",\\\"clientsecret\\\":\\\"5dfa16a0-12345\\\",\\\"credential-type\\\":\\\"binding-secret\\\",\\\"identityzone\\\":\\\"testcf-5vpfkrrj\\\",\\\"identityzoneid\\\":\\\"12345\\\",\\\"sburl\\\":\\\"https://internal-xsuaa.authentication.eu10.hana.ondemand.com\\\",\\\"serviceInstanceId\\\":\\\"123678\\\",\\\"subaccountid\\\":\\\"12345\\\",\\\"tenantid\\\":\\\"12345\\\",\\\"tenantmode\\\":\\\"dedicated\\\",\\\"uaadomain\\\":\\\"authentication.eu10.hana.ondemand.com\\\",\\\"url\\\":\\\"https://testcf-5vpfkrrj.authentication.eu10.hana.ondemand.com\\\",\\\"verificationkey\\\":\\\"-----BEGIN PUBLIC KEY-----\\n12abcdef34\\n12abcd34ef\\n12abcd34ef\\n12abcd34ef\\nabcdef\\n-----END PUBLIC KEY-----\\\",\\\"xsappname\\\":\\\"d12345678|one-mds-master!b28283\\\",\\\"zoneid\\\":\\\"12345\\\"},\\\"uri\\\":\\\"https://one-mds.cfapps.eu10.hana.ondemand.com\\\"}\",\"id\":\"56789\",\"labels\": {\"subaccount_id\": [\"12345\\\"\"]},\"last_modified\":\"2024-09-25T05:51:09Z\",\"name\":\"my binding\",\"ready\": true,\"service_instance_id\":\"123678\"},    {\"bind_resource\": null,\"context\":\"{\\\"crm_customer_id\\\":\\\"\\\",\\\"env_type\\\":\\\"sapcp\\\",\\\"global_account_id\\\":\\\"56789\\\",\\\"instance_name\\\":\\\"connectivity_instance\\\",\\\"license_type\\\":\\\"DEVELOPER\\\",\\\"origin\\\":\\\"sapcp\\\",\\\"platform\\\":\\\"sapcp\\\",\\\"region\\\":\\\"cf-eu10\\\",\\\"service_instance_id\\\":\\\"12345\\\",\\\"subaccount_id\\\":\\\"12345\\\",\\\"subdomain\\\":\\\"testcf-5vpfkrrj\\\",\\\"zone_id\\\":\\\"56789\\\"}\",\"created_date\":\"2024-09-26T08:17:00Z\",\"credentials\":\"{\\\"clientid\\\":\\\"sb-12345|connectivity!b17\\\",\\\"clientsecret\\\":\\\"A12345678\\\",\\\"connectivity_service\\\":{\\\"CAs_path\\\":\\\"/api/v1/CAs\\\",\\\"CAs_signing_path\\\":\\\"/api/v1/CAs/signing\\\",\\\"api_path\\\":\\\"/api/v1\\\",\\\"tunnel_path\\\":\\\"/api/v1/tunnel\\\",\\\"url\\\":\\\"https://connectivity.cf.eu10.hana.ondemand.com\\\"},\\\"credential-type\\\":\\\"binding-secret\\\",\\\"subaccount_id\\\":\\\"12345\\\",\\\"subaccount_subdomain\\\":\\\"testcf-5vpfkrrj\\\",\\\"token_service_domain\\\":\\\"authentication.eu10.hana.ondemand.com\\\",\\\"token_service_url\\\":\\\"https://testcf-5vpfkrrj.authentication.eu10.hana.ondemand.com/oauth/token\\\",\\\"token_service_url_pattern\\\":\\\"https://{tenant}.authentication.eu10.hana.ondemand.com/oauth/token\\\",\\\"token_service_url_pattern_tenant_key\\\":\\\"subaccount_subdomain\\\",\\\"xsappname\\\":\\\"clone12345!b507858|connectivity!b17\\\"}\",\"id\":\"34567\",\"labels\": {\"subaccount_id\": [\"12345\\\"\"]},\"last_modified\":\"2024-09-26T08:17:01Z\",\"name\":\"connection_binding\",\"ready\": true,\"service_instance_id\":\"56789\"}]}"
	dataMultipleServiceBindings, _ := GetDataFromJsonString(jsonStringMultipleBindings)

	tests := []struct {
		name          string
		data          map[string]interface{}
		subaccountId  string
		filterValues  []string
		expectedBlock string
		expectedCount int
		expectError   bool
	}{

		{
			name:          "Valid data without filter",
			data:          dataServiceBinding,
			subaccountId:  "12345",
			filterValues:  []string{},
			expectedBlock: "import {\n\t\t\t\tto = btp_subaccount_service_binding.servicebinding_0\n\t\t\t\tid = \"12345,56789\"\"\n\t\t\t  }\n\n",
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "Valid data with matching filter",
			data:          dataMultipleServiceBindings,
			subaccountId:  "12345",
			filterValues:  []string{"connection_binding"},
			expectedBlock: "import {\n\t\t\t\tto = btp_subaccount_service_binding.servicebinding_1\n\t\t\t\tid = \"12345,34567\"\"\n\t\t\t  }\n\n",
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "Invalid filter value",
			data:          dataServiceBinding,
			subaccountId:  "12345",
			filterValues:  []string{"wrong-binding-name"},
			expectedBlock: "",
			expectedCount: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			importBlock, count, err := createServiceBindingImportBlock(tt.data, tt.subaccountId, tt.filterValues, resourceDoc)
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
