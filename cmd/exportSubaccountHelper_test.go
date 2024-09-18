package cmd

import (
	"encoding/json"
	"testing"
)

func TestGetSubaccountImportBlock(t *testing.T) {

	jsonString := "{\"id\":\"5163621f-6a1e-4fbf-af3a-0f530a0dc4d5\",\"name\":\"testSubaccount\"}"
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		t.Errorf("error in unmarshalling")
	}

	importBlock, err := getSubaccountImportBlock(data, "5163621f-6a1e-4fbf-af3a-0f530a0dc4d5", nil)
	if err != nil {
		t.Errorf("error creating importBlock")
	}

	expectedValue := "import {\n\t\t\t\tto = btp_subaccount.testSubaccount\n\t\t\t\tid = \"5163621f-6a1e-4fbf-af3a-0f530a0dc4d5\"\n\t\t\t  }\n"

	if importBlock != expectedValue {
		t.Errorf("got %q, wanted %q", importBlock, expectedValue)
	}

}
