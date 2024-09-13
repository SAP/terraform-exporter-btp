package cmd

import (
	"encoding/json"
	"testing"
)

func TestReadSubaccountEntilementsDataSource(t *testing.T) {

	dataBlock, err := readSubaccountEntilementsDataSource("5163621f-6a1e-4fbf-af3a-0f530a0dc4d5")

	if err != nil {
		t.Errorf("error creating dataBlock")
	}

	expectedValue := "data \"btp_subaccount_entitlements\" \"all\"{\n  subaccount_id = \"5163621f-6a1e-4fbf-af3a-0f530a0dc4d5\" \n}\n"

	if dataBlock != expectedValue {
		t.Errorf("got %q, wanted %q", dataBlock, expectedValue)
	}

}

func TestGetEntitlementsImportBlock(t *testing.T) {

	// jsonString := "{\"id\":\"5163621f-6a1e-4fbf-af3a-0f530a0dc4d4\",\"subaccount_id\":\"5163621f-6a1e-4fbf-af3a-0f530a0dc4d4\",\"values\":[{\"environment_type\":\"cloudfoundry\",\"id\":\"BCC1E35C-EC86-49FC-96C5-5A3BBBB248C9\"}]}"
	jsonString := "{\"application-logs:lite\":{\"category\":\"ELASTIC_SERVICE\",\"plan_description\":\"Free offering for development purposes\",\"plan_display_name\":\"lite\",\"plan_name\":\"lite\",\"quota_assigned\":1,\"quota_remaining\":1,\"service_display_name\":\"Application Logging Service\",\"service_name\":\"application-logs\"}}"
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		t.Errorf("error in unmarshalling")
	}

	importBlock, err := getEntitlementsImportBlock(data, "5163621f-6a1e-4fbf-af3a-0f530a0dc4d5", nil)
	if err != nil {
		t.Errorf("error creating importBlock")
	}

	expectedValue := "import {\n\t\t\t\tto = btp_subaccount_entitlement.application-logs_lite\n\t\t\t\tid = \"5163621f-6a1e-4fbf-af3a-0f530a0dc4d5,application-logs,lite\"\n\t\t\t  }\n"

	if importBlock != expectedValue {
		t.Errorf("got %q, wanted %q", importBlock, expectedValue)
	}

}
