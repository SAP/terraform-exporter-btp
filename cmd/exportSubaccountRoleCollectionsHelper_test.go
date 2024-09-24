package cmd

import (
	"encoding/json"
	"testing"
)

func TestGetRoleCollectionsImportBlock(t *testing.T) {

	jsonString := "{\"id\":\"5163621f-6a1e-4fbf-af3a-0f530a0dc4d4\",\"subaccount_id\":\"5163621f-6a1e-4fbf-af3a-0f530a0dc4d4\",\"values\":[{\"description\": \"Operate the data transmission tunnels used by the Cloud Connector.\", \"name\": \"Cloud Connector Administrator\", \"read_only\": true, \"roles\": [{ \"description\": \"Operate the data transmission tunnels and client certificates used by the Cloud connector\", \"name\": \"Cloud Connector Administrator\", \"role_template_app_id\": \"connectivity!b7\", \"role_template_name\": \"Cloud_Connector_Administrator\"}]}]}"
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		t.Errorf("error in unmarshalling")
	}

	importBlock, err := getSubaccountRoleCollectionsImportBlock(data, "5163621f-6a1e-4fbf-af3a-0f530a0dc4d5", nil)
	if err != nil {
		t.Errorf("error creating importBlock")
	}

	expectedValue := "import {\n\t\t\t\tto = btp_subaccount_role_collection.cloud_connector_administrator\n\t\t\t\tid = \"5163621f-6a1e-4fbf-af3a-0f530a0dc4d5,Cloud Connector Administrator\"\n\t\t\t  }\n"

	if importBlock != expectedValue {
		t.Errorf("got %q, wanted %q", importBlock, expectedValue)
	}

}
