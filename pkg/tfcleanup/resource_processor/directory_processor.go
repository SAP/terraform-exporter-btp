package resourceprocessor

import (
	"github.com/SAP/terraform-exporter-btp/internal/btpcli"
	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func processDirectoryAttributes(body *hclwrite.Body, variables *generictools.VariableContent, btpClient *btpcli.ClientFacade) {
	attrs := body.Attributes()
	for name, attr := range attrs {
		tokens := attr.Expr().BuildTokens(nil)

		if name == parentIdentifier && len(tokens) == 3 {

			parentId := generictools.GetStringToken(tokens)

			if generictools.IsGlobalAccountParent(btpClient, parentId) {
				body.RemoveAttribute(name)
			} else {

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
}
