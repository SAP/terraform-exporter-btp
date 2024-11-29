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

func TestFormatResourceNameGenericWithDots(t *testing.T) {

	input := "ApiManagement.SelfService.Administrator"
	expected := "apimanagement_selfservice_administrator"

	result := FormatResourceNameGeneric(input)

	if result != expected {
		t.Errorf("got %q, wanted %q", result, expected)
	}
}

func TestFormatResourceNameGenericMisc(t *testing.T) {

	input := "SomeResource withSpaces.AndDots"
	expected := "someresource_withspaces_anddots"

	result := FormatResourceNameGeneric(input)

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
