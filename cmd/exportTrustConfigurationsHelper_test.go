package cmd

import (
	"encoding/json"
	"testing"
)

func TestGetTrustConfigurationsImportBlock(t *testing.T) {

	jsonString := "{\"id\":\"5163621f-6a1e-4fbf-af3a-0f530a0dc4d5\",\"subaccount_id\":\"5163621f-6a1e-4fbf-af3a-0f530a0dc4d5\",\"values\":[{\"origin\":\"testorigin\"}]}"
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		t.Errorf("error in unmarshalling")
	}

	importBlock, err := getSubaccountTrustConfigurationsImportBlock(data, "5163621f-6a1e-4fbf-af3a-0f530a0dc4d5", nil)
	if err != nil {
		t.Errorf("error creating importBlock")
	}

	expectedValue := "import {\n\t\t\t\tto = btp_subaccount_trust_configuration.trust0\n\t\t\t\tid = \"5163621f-6a1e-4fbf-af3a-0f530a0dc4d5,testorigin\"\n\t\t\t  }\n"

	if importBlock != expectedValue {
		t.Errorf("got %q, wanted %q", importBlock, expectedValue)
	}

}
