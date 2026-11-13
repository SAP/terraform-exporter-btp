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

func FormatOrgRoleResourceName(orgRoleType string, userId string) string {
	return orgRoleType + "_" + userId
}

func FormatSpaceRoleResourceName(spaceRoleType string, spaceId string, userId string) string {
	spaceName, _ := cfcli.GetSpaceName(spaceId)
	userName, _ := cfcli.GetUser(userId)
	cleanSpaceName := strings.ReplaceAll(spaceName, "_", "-")
	cleanUserName := strings.ReplaceAll(userName, "_", "-")
	cleanRoleType := strings.ReplaceAll(spaceRoleType, "_", "-")
	// Usernames can be shared across identity-provider origins. Retain the user
	// GUID so each assignment has a distinct inventory identifier.
	return cleanSpaceName + "_" + cleanRoleType + "_" + cleanUserName + "_" + spaceId + "_" + userId
}
