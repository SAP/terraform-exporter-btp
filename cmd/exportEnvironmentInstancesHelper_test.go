package cmd

import (
	"encoding/json"
	"testing"
)

func TestReadDataSource(t *testing.T) {
	//data := map[string]interface {} ["id": "5163621f-6a1e-4fbf-af3a-0f530a0dc4d4", "subaccount_id": "5163621f-6a1e-4fbf-af3a-0f530a0dc4d4", "values": struct {id string, instance_name string}{"BCC1E35C-EC86-49FC-96C5-5A3BBBB248C9","terraform-integration-prod_testcf1-zazq2cbi"}]
	dataBlock, err := readDataSource("5163621f-6a1e-4fbf-af3a-0f530a0dc4d5", SubaccountEnvironmentInstanceType)

	if err != nil {
		t.Errorf("error creating dataBlock")
	}

	expectedValue := "data \"btp_subaccount_environment_instances\" \"all\"{\n  subaccount_id = \"5163621f-6a1e-4fbf-af3a-0f530a0dc4d5\" \n}\n"

	if dataBlock != expectedValue {
		t.Errorf("got %q, wanted %q", dataBlock, expectedValue)
	}

}

func TestGetImportBlock(t *testing.T) {

	jsonString := "{\"id\":\"5163621f-6a1e-4fbf-af3a-0f530a0dc4d4\",\"subaccount_id\":\"5163621f-6a1e-4fbf-af3a-0f530a0dc4d4\",\"values\":[{\"environment_type\":\"cloudfoundry\",\"id\":\"BCC1E35C-EC86-49FC-96C5-5A3BBBB248C9\"}]}"
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		t.Errorf("error in unmarshalling")
	}

	importBlock, err := getSubaccountEnvironmentInstanceBlock(data, "5163621f-6a1e-4fbf-af3a-0f530a0dc4d5", nil)
	if err != nil {
		t.Errorf("error creating importBlock")
	}

	expectedValue := "import {\n\t\t\t\tto = btp_subaccount_environment_instance.cloudfoundry\n\t\t\t\tid = \"5163621f-6a1e-4fbf-af3a-0f530a0dc4d5,BCC1E35C-EC86-49FC-96C5-5A3BBBB248C9\"\n\t\t\t  }\n"

	if importBlock != expectedValue {
		t.Errorf("got %q, wanted %q", importBlock, expectedValue)
	}

}
