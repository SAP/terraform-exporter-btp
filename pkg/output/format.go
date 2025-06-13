package output

import (
	"strings"

	"github.com/SAP/terraform-exporter-btp/internal/cfcli"
)

func FormatResourceNameGeneric(name string) string {
	return strings.ToLower(FormatResourceNameGenericNoCapConversion(name))
}

func FormatResourceNameGenericNoCapConversion(name string) string {
	return strings.ReplaceAll(name, " ", "_")
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

var FormatRoles = FormatSpaceRoleResourceName

func FormatSpaceRoleResourceName(spaceRoleType string, spaceId string, userId string) string {
	spaceName, _ := cfcli.GetSpaceName(spaceId)
	userName, _ := cfcli.GetUser(userId)
	cleanSpaceName := strings.ReplaceAll(spaceName, "_", "-")
	cleanUserName := strings.ReplaceAll(userName, "_", "-")
	cleanRoleType := strings.ReplaceAll(spaceRoleType, "_", "-")
	return cleanSpaceName + "_" + cleanRoleType + "_" + cleanUserName + "_" + spaceId
}
