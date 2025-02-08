package providerprocessor

import (
	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func ProcessProvider(hclFile *hclwrite.File, variables *generictools.VariableContent) {
	processProviderAttributes(hclFile.Body(), nil, variables)
}

func processProviderAttributes(body *hclwrite.Body, inBlocks []string, variables *generictools.VariableContent) {
	attrs := body.Attributes()
	for name, attr := range attrs {
		tokens := attr.Expr().BuildTokens(nil)

		if name == globalAccountIdentifier && len(tokens) == 3 {
			replacedTokens, globalAccountValue := generictools.ReplaceStringToken(tokens, globalAccountIdentifier)
			if globalAccountValue != "" {
				(*variables)[name] = generictools.VariableInfo{
					Description: "Global account subdomain",
					Value:       globalAccountValue,
				}
			}
			body.SetAttributeRaw(name, replacedTokens)
		}
	}

	blocks := body.Blocks()
	for _, block := range blocks {
		inBlocks := append(inBlocks, block.Type())
		processProviderAttributes(block.Body(), inBlocks, variables)
	}
}
