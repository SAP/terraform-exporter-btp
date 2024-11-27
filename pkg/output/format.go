package output

import "strings"

func FormatResourceNameGeneric(name string) string {
	return strings.ToLower(strings.Replace(name, " ", "_", -1))
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

func FormatOrgRoleResourceName(orgRoleType string, userId string) string {
	return orgRoleType + "_" + userId
}
