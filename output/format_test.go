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

func TestFormatSubscriptionResourceName(t *testing.T) {

	appName := "feature-flags-dashboard"
	planName := "dashboard"
	expected := "feature-flags-dashboard_dashboard"

	result := FormatSubscriptionResourceName(appName, planName)

	if result != expected {
		t.Errorf("got %q, wanted %q", result, expected)
	}

}
