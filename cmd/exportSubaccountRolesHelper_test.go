package cmd

import (
	"encoding/json"
	"testing"
)

func TestGetRolesImportBlock(t *testing.T) {

	jsonString := "{\"id\":\"5163621f-6a1e-4fbf-af3a-0f530a0dc4d4\",\"subaccount_id\":\"5163621f-6a1e-4fbf-af3a-0f530a0dc4d4\",\"values\":[{\"app_id\": \"destination-xsappname!b62\", \"app_name\": \"destination-xsappname\", \"description\": \"Manage destination configurations, certificates and signing keys for SAML assertions issued by the Destination service\", \"name\": \"Application Destination Administrator\", \"read_only\": \"true\", \"role_template_name\": \"Application_Destination_Administrator\"}]}"
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		t.Errorf("error in unmarshalling")
	}

	importBlock, err := getSubaccountRolesImportBlock(data, "5163621f-6a1e-4fbf-af3a-0f530a0dc4d5", nil)
	if err != nil {
		t.Errorf("error creating importBlock")
	}

	expectedValue := "import {\n\t\t\t\tto = btp_subaccount_role.application_destination_administrator\n\t\t\t\tid = \"5163621f-6a1e-4fbf-af3a-0f530a0dc4d5,Application Destination Administrator,Application_Destination_Administrator,destination-xsappname!b62\"\n\t\t\t  }\n"

	if importBlock != expectedValue {
		t.Errorf("got %q, wanted %q", importBlock, expectedValue)
	}

}
