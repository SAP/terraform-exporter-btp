package cmd

import (
	"encoding/json"
	"testing"
)

func TestGetServiceBindingsImportBlock(t *testing.T) {

	jsonString := "{\"id\":\"5163621f-6a1e-4fbf-af3a-0f530a0dc4d4\",\"subaccount_id\":\"5163621f-6a1e-4fbf-af3a-0f530a0dc4d4\",\"values\":[{\"id\": \"910e9a7d-0fb4-4428-a813-56550e683579\",\"name\": \"test binding\"}]}"
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		t.Errorf("error in unmarshalling")
	}

	importBlock, err := getSubaccountServiceBindingImportBlock(data, "5163621f-6a1e-4fbf-af3a-0f530a0dc4d5", nil)
	if err != nil {
		t.Errorf("error creating importBlock")
	}

	expectedValue := "import {\n\t\t\t\tto = btp_subaccount_service_binding.test_binding\n\t\t\t\tid = \"5163621f-6a1e-4fbf-af3a-0f530a0dc4d5,910e9a7d-0fb4-4428-a813-56550e683579\"\n\t\t\t  }\n"

	if importBlock != expectedValue {
		t.Errorf("got %q, wanted %q", importBlock, expectedValue)
	}

}
