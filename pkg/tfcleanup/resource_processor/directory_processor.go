package resourceprocessor

import (
	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func processDirectoryAttributes(body *hclwrite.Body, variables *generictools.VariableContent) {
	attrs := body.Attributes()
	for name, attr := range attrs {
		tokens := attr.Expr().BuildTokens(nil)

		if name == parentIdentifier && len(tokens) == 3 {
			replacedTokens, parentValue := generictools.ReplaceStringToken(tokens, parentIdentifier)
			if parentValue != "" {
				(*variables)[name] = generictools.VariableInfo{
					Description: "ID of the parent of the SAP BTP directory",
					Value:       parentValue,
				}
			}
			body.SetAttributeRaw(name, replacedTokens)
		}
	}
}
