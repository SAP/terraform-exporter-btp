package resourceprocessor

import (
	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

const subaccountRoleBlockIdentifier = "btp_subaccount_role"
const directoryRoleBlockIdentifier = "btp_directory_role"
const roleAppIdIdentifier = "app_id"
const roleNameIdentifier = "name"
const roleTemplateNameIdentifier = "role_template_name"

func fillRoleDependencyAddresses(body *hclwrite.Body, resourceAddress string, dependencyAddresses *generictools.DependencyAddresses) {
	appIdAttr := body.GetAttribute(roleAppIdIdentifier)
	roleNameAttr := body.GetAttribute(roleNameIdentifier)
	roleTemplateNameAttr := body.GetAttribute(roleTemplateNameIdentifier)

	if appIdAttr == nil || roleNameAttr == nil || roleTemplateNameAttr == nil {
		return
	}

	appIdToken := appIdAttr.Expr().BuildTokens(nil)
	roleNameToken := roleNameAttr.Expr().BuildTokens(nil)
	roleTemplateNameToken := roleTemplateNameAttr.Expr().BuildTokens(nil)

	appId := generictools.GetStringToken(appIdToken)
	roleName := generictools.GetStringToken(roleNameToken)
	roleTemplateName := generictools.GetStringToken(roleTemplateNameToken)

	if appId != "" && roleName != "" && roleTemplateName != "" {
		key := generictools.RoleKey{
			AppId:            appId,
			Name:             roleName,
			RoleTemplateName: roleTemplateName,
		}

		(*dependencyAddresses).RoleAddress[key] = resourceAddress
	}
}
