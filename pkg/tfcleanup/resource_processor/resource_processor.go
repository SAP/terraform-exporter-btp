package resourceprocessor

import (
	"log"

	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/SAP/terraform-exporter-btp/pkg/tfutils"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func ProcessResources(hclFile *hclwrite.File, variables *generictools.VariableContent, level string) {

	processResourceAttributes(hclFile.Body(), nil, variables, level)
}

func processResourceAttributes(body *hclwrite.Body, inBlocks []string, variables *generictools.VariableContent, level string) {

	if len(inBlocks) > 0 {
		// remove empty values for all resources
		removeEmptyAttributes(body)

		switch level {
		case tfutils.SubaccountLevel:
			if inBlocks[0] == subaccountBlockIdentifier {
				processSubaccountAttributes(body, variables)
			}

			/*	if inBlocks[0] == subscriptionBlockIdentifier {
					attrs := body.Attributes()
					for name, attr := range attrs {
						tokens := attr.Expr().BuildTokens(nil)

						fmt.Println("Name: ", name)
						fmt.Println("Tokens: ", tokens)
						fmt.Println("Value of first token: ", string(tokens[0].Bytes))
						fmt.Println("=====================================")
					}
				}
			*/

		case tfutils.DirectoryLevel:
			log.Println("Directory level is not supported yet")
		case tfutils.OrganizationLevel:
			log.Println("Organization level is not supported yet")
		}
	}
	blocks := body.Blocks()
	for _, block := range blocks {
		inBlocks := append(inBlocks, block.Type()+"_"+block.Labels()[0])
		processResourceAttributes(block.Body(), inBlocks, variables, level)
	}
}

func removeEmptyAttributes(body *hclwrite.Body) {
	attrs := body.Attributes()
	for name, attr := range attrs {
		tokens := attr.Expr().BuildTokens(nil)

		if len(tokens) == 1 && string(tokens[0].Bytes) == generictools.EmptyString {
			body.RemoveAttribute(name)
		}

		if len(tokens) == 5 {
			var combinedString string
			for _, token := range tokens {
				combinedString += string(token.Bytes)
			}
			if combinedString == generictools.EmptyString {
				body.RemoveAttribute(name)
			}
		}
	}
}
