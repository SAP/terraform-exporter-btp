package output

import "strings"

func FormatResourceNameGeneric(name string) string {
	nameNoSpaces := strings.ToLower(strings.Replace(name, " ", "_", -1))
	return strings.Replace(nameNoSpaces, ".", "_", -1)
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
