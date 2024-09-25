package output

import "testing"

func TestFormatRoleResourceName(t *testing.T) {

	input := "Application Destination Administrator"
	expected := "application_destination_administrator"

	result := FormatRoleResourceName(input)

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
func TestFormatServiceBindingResourceName(t *testing.T) {

	input := "My App binding"
	expected := "my_app_binding"
	result := FormatServiceBindingResourceName(input)
	if result != expected {
		t.Errorf("got %q, wanted %q", result, expected)
	}

}
