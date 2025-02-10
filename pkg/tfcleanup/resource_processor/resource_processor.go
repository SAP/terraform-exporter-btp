package resourceprocessor

import (
	"log"
	"strings"

	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/SAP/terraform-exporter-btp/pkg/tfutils"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func ProcessResources(hclFile *hclwrite.File, level string, variables *generictools.VariableContent, dependencyAddresses *generictools.DepedendcyAddresses) {

	processResourceAttributes(hclFile.Body(), nil, level, variables, dependencyAddresses)
}

func processResourceAttributes(body *hclwrite.Body, inBlocks []string, level string, variables *generictools.VariableContent, dependencyAddresses *generictools.DepedendcyAddresses) {

	if len(inBlocks) > 0 {
		// remove empty values for all resources
		removeEmptyAttributes(body)

		// Get the first part of the block until the comma
		blockIdentifier := strings.Split(inBlocks[0], ",")[1]
		blockAddress := strings.Split(inBlocks[0], ",")[2]
		resourceAddress := blockIdentifier + "." + blockAddress

		switch level {
		case tfutils.SubaccountLevel:
			if blockIdentifier == subaccountBlockIdentifier {
				processSubaccountAttributes(body, variables)
				//Add Address of subaccount to the dependencyAddresses
				dependencyAddresses.SubaccountAddress = resourceAddress
			}

			if blockIdentifier != subaccountBlockIdentifier {
				replaceMainDependency(body, subaccountIdentifier, dependencyAddresses.SubaccountAddress)
			}

			if inBlocks[0] == subaccountEntitlementBlockIdentifier {
				fillSubaccountEntitlementDependencyAddresses(body, resourceAddress, dependencyAddresses)
			}

			if inBlocks[0] == subscriptionBlockIdentifier {
				addEntitlementDependency(body, dependencyAddresses)
			}

		case tfutils.DirectoryLevel:
			if blockIdentifier != directoryBlockIdentifier {
				replaceMainDependency(body, directoryIdentifier, dependencyAddresses.SubaccountAddress)
			}
		case tfutils.OrganizationLevel:
			log.Println("Organization level is not supported yet")
		}
	}
	blocks := body.Blocks()
	for _, block := range blocks {
		inBlocks := append(inBlocks, block.Type()+","+block.Labels()[0]+","+block.Labels()[1])
		processResourceAttributes(block.Body(), inBlocks, level, variables, dependencyAddresses)
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
			if combinedString == generictools.EmptyJson {
				body.RemoveAttribute(name)
			}
		}
	}
}

func replaceMainDependency(body *hclwrite.Body, mainIdentifier string, mainAddress string) {
	attrs := body.Attributes()
	for name, attr := range attrs {
		tokens := attr.Expr().BuildTokens(nil)

		if name == mainIdentifier && len(tokens) == 3 {
			replacedTokens := generictools.ReplaceDependency(tokens, mainAddress)
			body.SetAttributeRaw(name, replacedTokens)
		}
	}
}
