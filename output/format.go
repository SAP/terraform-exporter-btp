package output

import "strings"

func FormatRoleResourceName(name string) string {
	return strings.ToLower(strings.Replace(name, " ", "_", -1))
}

func FormatSubscriptionResourceName(appName string, planName string) string {
	return appName + "_" + planName
}

func FormatRoleCollectionResourceName(name string) string {
	return strings.ToLower(strings.Replace(name, " ", "_", -1))
}

func FormatServiceBindingResourceName(name string) string {
	return strings.ToLower(strings.Replace(name, " ", "_", -1))
}
