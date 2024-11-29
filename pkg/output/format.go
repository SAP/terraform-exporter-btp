package output

import "strings"

func FormatResourceNameGeneric(name string) string {
	return strings.ToLower(strings.Replace(name, " ", "_", -1))
}

func FormatRoleResourceName(name string) string {
	nameGeneric := FormatResourceNameGeneric(name)
	// We replace a dot with two underscores to avoid conflicts with same role names that already contain a underscore
	return strings.Replace(nameGeneric, ".", "__", -1)
}

func FormatDirEntitlementResourceName(appName string, planName string) string {
	return appName + "_" + planName
}

func FormatSubscriptionResourceName(appName string, planName string) string {
	return appName + "_" + planName
}

func FormatServiceInstanceResourceName(serviceInstanceName string, planId string) string {
	return serviceInstanceName + "_" + planId
}
