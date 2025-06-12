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

func TestFormatSubscriptionResourceName(t *testing.T) {

	appName := "feature-flags-dashboard"
	planName := "dashboard"
	expected := "feature-flags-dashboard_dashboard"

	result := FormatSubscriptionResourceName(appName, planName)

	if result != expected {
		t.Errorf("got %q, wanted %q", result, expected)
	}
}

func TestFormatServiceInstanceResourceName(t *testing.T) {

	serviceINstanceName := "audit-log-exporter"
	planId := "a50128a9-35fc-4624-9953-c79668ef3e5b"
	expected := "audit-log-exporter_a50128a9-35fc-4624-9953-c79668ef3e5b"

	result := FormatServiceInstanceResourceName(serviceINstanceName, planId)

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

func TestFormatDirEntitlementResourceName(t *testing.T) {

	appName := "feature-flags-dashboard"
	planName := "dashboard"
	expected := "feature-flags-dashboard_dashboard"

	result := FormatDirEntitlementResourceName(appName, planName)

	if result != expected {
		t.Errorf("got %q, wanted %q", result, expected)
	}
}
