package tfutils

import "testing"

func TestReadDataSource(t *testing.T) {
	//data := map[string]interface {} ["id": "5163621f-6a1e-4fbf-af3a-0f530a0dc4d4", "subaccount_id": "5163621f-6a1e-4fbf-af3a-0f530a0dc4d4", "values": struct {id string, instance_name string}{"BCC1E35C-EC86-49FC-96C5-5A3BBBB248C9","terraform-integration-prod_testcf1-zazq2cbi"}]
	dataBlock, err := readDataSource("5163621f-6a1e-4fbf-af3a-0f530a0dc4d5", "", "", SubaccountEnvironmentInstanceType)

	if err != nil {
		t.Errorf("error creating dataBlock")
	}

	expectedValue := "data \"btp_subaccount_environment_instances\" \"all\"{\n  subaccount_id = \"5163621f-6a1e-4fbf-af3a-0f530a0dc4d5\" \n}\n"

	if dataBlock != expectedValue {
		t.Errorf("got %q, wanted %q", dataBlock, expectedValue)
	}

}

func TestReadSubaccountEntilementsDataSource(t *testing.T) {

	dataBlock, err := readDataSource("5163621f-6a1e-4fbf-af3a-0f530a0dc4d5", "", "", SubaccountEntitlementType)

	if err != nil {
		t.Errorf("error creating dataBlock")
	}

	expectedValue := "data \"btp_subaccount_entitlements\" \"all\"{\n  subaccount_id = \"5163621f-6a1e-4fbf-af3a-0f530a0dc4d5\" \n}\n"

	if dataBlock != expectedValue {
		t.Errorf("got %q, wanted %q", dataBlock, expectedValue)
	}

}

func TestReadSubaccountDataSource(t *testing.T) {
	//data := map[string]interface {} ["id": "5163621f-6a1e-4fbf-af3a-0f530a0dc4d4", "subaccount_id": "5163621f-6a1e-4fbf-af3a-0f530a0dc4d4", "values": struct {id string, instance_name string}{"BCC1E35C-EC86-49FC-96C5-5A3BBBB248C9","terraform-integration-prod_testcf1-zazq2cbi"}]
	dataBlock, err := readDataSource("5163621f-6a1e-4fbf-af3a-0f530a0dc4d5", "", "", SubaccountType)

	if err != nil {
		t.Errorf("error creating dataBlock")
	}

	expectedValue := "data \"btp_subaccount\" \"all\"{\n  id = \"5163621f-6a1e-4fbf-af3a-0f530a0dc4d5\" \n}\n"

	if dataBlock != expectedValue {
		t.Errorf("got %q, wanted %q", dataBlock, expectedValue)
	}

}

func TestReadSubaccountSubscriptionDataSource(t *testing.T) {
	//data := map[string]interface {} ["id": "5163621f-6a1e-4fbf-af3a-0f530a0dc4d4", "subaccount_id": "5163621f-6a1e-4fbf-af3a-0f530a0dc4d4", "values": struct {id string, instance_name string}{"BCC1E35C-EC86-49FC-96C5-5A3BBBB248C9","terraform-integration-prod_testcf1-zazq2cbi"}]
	dataBlock, err := readDataSource("5163621f-6a1e-4fbf-af3a-0f530a0dc4d5", "", "", SubaccountSubscriptionType)

	if err != nil {
		t.Errorf("error creating dataBlock")
	}

	expectedValue := "data \"btp_subaccount_subscriptions\" \"all\"{\n  subaccount_id = \"5163621f-6a1e-4fbf-af3a-0f530a0dc4d5\" \n}\n"

	if dataBlock != expectedValue {
		t.Errorf("got %q, wanted %q", dataBlock, expectedValue)
	}

}

func TestReadSubaccountTrustConfigurationsDataSource(t *testing.T) {
	//data := map[string]interface {} ["id": "5163621f-6a1e-4fbf-af3a-0f530a0dc4d4", "subaccount_id": "5163621f-6a1e-4fbf-af3a-0f530a0dc4d4", "values": struct {id string, instance_name string}{"BCC1E35C-EC86-49FC-96C5-5A3BBBB248C9","terraform-integration-prod_testcf1-zazq2cbi"}]
	dataBlock, err := readDataSource("5163621f-6a1e-4fbf-af3a-0f530a0dc4d5", "", "", SubaccountTrustConfigurationType)

	if err != nil {
		t.Errorf("error creating dataBlock")
	}

	expectedValue := "data \"btp_subaccount_trust_configurations\" \"all\"{\n  subaccount_id = \"5163621f-6a1e-4fbf-af3a-0f530a0dc4d5\" \n}\n"

	if dataBlock != expectedValue {
		t.Errorf("got %q, wanted %q", dataBlock, expectedValue)
	}

}
