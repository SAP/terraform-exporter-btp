package resourceprocessor

import (
	"github.com/SAP/terraform-exporter-btp/internal/btpcli"
	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func processSubaccountAttributes(body *hclwrite.Body, variables *generictools.VariableContent, btpClient *btpcli.ClientFacade) {
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

		if name == parentIdentifier && len(tokens) == 3 {

			parentId := generictools.GetStringToken(tokens)

			if generictools.IsGlobalAccountParent(btpClient, parentId) {
				body.RemoveAttribute(name)
			} else {
				replacedTokens, parentValue := generictools.ReplaceStringToken(tokens, parentIdentifier)
				if parentValue != "" {
					(*variables)[name] = generictools.VariableInfo{
						Description: "ID of the parent of the SAP BTP subaccount",
						Value:       parentValue,
					}
				}
				body.SetAttributeRaw(name, replacedTokens)
			}
		}
	}
}
