package output

import "strings"

func FormatResourceNameGeneric(name string) string {
	return strings.ToLower(strings.Replace(name, " ", "_", -1))
}

func FormatSubscriptionResourceName(appName string, planName string) string {
	return appName + "_" + planName
}

func FormatServiceInstanceResourceName(serviceInstanceName string, planId string) string {
	return serviceInstanceName + "_" + planId
}
