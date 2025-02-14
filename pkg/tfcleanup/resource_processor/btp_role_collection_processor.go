package resourceprocessor

import (
	"encoding/json"
	"fmt"
	"strings"

	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

const subaccountRoleCollectionBlockIdentifier = "btp_subaccount_role_collection"
const directoryRoleCollectionBlockIdentifier = "btp_directory_role_collection"
const roleBlockIdentifier = "roles"

type Role struct {
	Name              string `json:"name"`
	RoleTemplateAppID string `json:"role_template_app_id"`
	RoleTemplateName  string `json:"role_template_name"`
}

func addRoleDependency(body *hclwrite.Body, dependencyAddresses *generictools.DepedendcyAddresses) {
	roleAttr := body.GetAttribute(roleBlockIdentifier)

	if roleAttr == nil {
		return
	}

	roleAttrTokens := roleAttr.Expr().BuildTokens(nil)

	var roleString string

	for _, token := range roleAttrTokens {
		roleString = roleString + string(token.Bytes)
	}

	roleBlock := preprocessString(roleString)

	var roles []Role
	err := json.Unmarshal([]byte(roleBlock), &roles)
	if err != nil {
		fmt.Println("Error unmarshaling roles assigned to role collection:", err)
		return
	}

	var builder strings.Builder
	first := true

	for _, roleEntry := range roles {
		searchKey := generictools.RoleKey{
			AppId:            roleEntry.RoleTemplateAppID,
			Name:             roleEntry.Name,
			RoleTemplateName: roleEntry.RoleTemplateName,
		}

		dependencyAddress := (*dependencyAddresses).RoleAddress[searchKey]

		if dependencyAddress != "" {
			if !first {
				builder.WriteString(", ")
			}
			builder.WriteString(dependencyAddress)
			first = false
		}
	}

	dependencies := builder.String()

	if dependencies != "" {
		body.SetAttributeRaw("depends_on", hclwrite.Tokens{
			{
				Type:  hclsyntax.TokenOBrack,
				Bytes: []byte("["),
			},
			{
				Type:  hclsyntax.TokenStringLit,
				Bytes: []byte(dependencies),
			},
			{
				Type:  hclsyntax.TokenCBrack,
				Bytes: []byte("]"),
			},
		})
	}
}

func preprocessString(input string) string {
	input = strings.ReplaceAll(input, "=", ":")
	return input
}
