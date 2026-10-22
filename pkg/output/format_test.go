package output

import "testing"

func TestFormatResourceNameGeneric(t *testing.T) {

	input := "Application Destination Administrator"
	expected := "application_destination_administrator"

	result := FormatResourceNameGeneric(input)

	if result != expected {
		t.Errorf("got %q, wanted %q", result, expected)
	}
}
func TestFormatResourceNameGenericNoCapConversion(t *testing.T) {

	input := "Application Destination Administrator"
	expected := "Application_Destination_Administrator"

	result := FormatResourceNameGenericNoCapConversion(input)

	if result != expected {
		t.Errorf("got %q, wanted %q", result, expected)
	}
}

func TestFormatOrgRoleResourceName(t *testing.T) {

	orgRoleType := "org_manager"
	userId := "someID"
	expected := "org_manager_someID"

	result := FormatOrgRoleResourceName(orgRoleType, userId)

	if result != expected {
		t.Errorf("got %q, wanted %q", result, expected)
	}
}
