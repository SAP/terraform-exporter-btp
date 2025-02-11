package resourceprocessor

import (
	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func extractOrgIds(body *hclwrite.Body, variables *generictools.VariableContent, orgId string) {
	attrs := body.Attributes()
	for name, attr := range attrs {
		tokens := attr.Expr().BuildTokens(nil)

		if name == cfOrgIdentifier && len(tokens) == 3 {
			replacedTokens, _ := generictools.ReplaceStringToken(tokens, cfOrgIdentifier)
			(*variables)[name] = generictools.VariableInfo{
				Description: "ID of the Cloud Foundry Organization",
				Value:       orgId,
			}
			body.SetAttributeRaw(name, replacedTokens)
		}
	}
}
