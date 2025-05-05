package resourceprocessor

import (
	"encoding/json"
	"fmt"
	"regexp"
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

func addRoleDependency(body *hclwrite.Body, dependencyAddresses *generictools.DependencyAddresses) {
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

	dependencies := buildDependencyString(roles, dependencyAddresses)

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
	// We must process the raw string extracted from the HCL file to make it a valid JSON string
	input = strings.ReplaceAll(input, "=", ":")

	// Add double quotes around keys and values
	re := regexp.MustCompile(`(\w+):`)
	input = re.ReplaceAllString(input, `"$1":`)

	re = regexp.MustCompile(`:"([^"]+)"`)
	input = re.ReplaceAllString(input, `:"$1"`)

	// Replace newlines with commas
	input = strings.ReplaceAll(input, "\n", ",")
	input = strings.ReplaceAll(input, ",,", ",")

	// Remove trailing commas before closing braces }
	re = regexp.MustCompile(`,(\s*})`)
	input = re.ReplaceAllString(input, `$1`)

	// Remove trailing commas after opening braces {
	re = regexp.MustCompile(`({\s*),`)
	input = re.ReplaceAllString(input, `$1`)

	// Remove trailing commas before closing braces ]
	re = regexp.MustCompile(`,(\s*])`)
	input = re.ReplaceAllString(input, `$1`)

	// remove a comma after [ if it exists
	re = regexp.MustCompile(`\[\s*,`)
	input = re.ReplaceAllString(input, `[`)

	// remove anything that comes after ] as this is the end of the array
	re = regexp.MustCompile(`\].*`)
	input = re.ReplaceAllString(input, `]`)

	return input
}

func buildDependencyString(roles []Role, dependencyAddresses *generictools.DependencyAddresses) string {
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

	return builder.String()
}
