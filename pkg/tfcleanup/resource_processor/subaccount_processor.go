package resourceprocessor

import (
	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func processSubaccountAttributes(body *hclwrite.Body, variables *generictools.VariableContent) {
	attrs := body.Attributes()
	for name, attr := range attrs {
		tokens := attr.Expr().BuildTokens(nil)

		if name == regionIdentifier && len(tokens) == 3 {
			replacedTokens, regionValue := generictools.ReplaceStringToken(tokens, regionIdentifier)
			if regionValue != "" {
				(*variables)[name] = generictools.VariableInfo{
					Description: "Region of SAP BTP subaccount",
					Value:       regionValue,
				}
			}
			body.SetAttributeRaw(name, replacedTokens)
		}
	}
}
