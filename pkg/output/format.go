package output

import "strings"

func FormatResourceNameGeneric(name string) string {
	return strings.ToLower(strings.Replace(name, " ", "_", -1))
}

func FormatRoleResourceName(name string) string {
	nameGeneric := FormatResourceNameGeneric(name)
	return strings.Replace(nameGeneric, ".", "_", -1)
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
