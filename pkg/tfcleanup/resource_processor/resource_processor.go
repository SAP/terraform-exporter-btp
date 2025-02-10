package resourceprocessor

import (
	"log"

	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/SAP/terraform-exporter-btp/pkg/tfutils"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func ProcessResources(hclFile *hclwrite.File, level string, variables *generictools.VariableContent, dependencyAddresses *generictools.DepedendcyAddresses) {

	processResourceAttributes(hclFile.Body(), nil, level, variables, dependencyAddresses)
}

func processResourceAttributes(body *hclwrite.Body, inBlocks []string, level string, variables *generictools.VariableContent, dependencyAddresses *generictools.DepedendcyAddresses) {

	if len(inBlocks) > 0 {

		removeEmptyAttributes(body)

		blockIdentifier, resourceAddress := generictools.ExtractBlockInformation(inBlocks)

		switch level {
		case tfutils.SubaccountLevel:
			processSubaccountLevel(body, variables, dependencyAddresses, blockIdentifier, resourceAddress)
		case tfutils.DirectoryLevel:
			processDirectoryLevel(body, variables, dependencyAddresses, blockIdentifier, resourceAddress)
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
	if mainAddress == "" {
		return
	}

	attrs := body.Attributes()
	for name, attr := range attrs {
		tokens := attr.Expr().BuildTokens(nil)

		if name == mainIdentifier && len(tokens) == 3 {
			replacedTokens := generictools.ReplaceDependency(tokens, mainAddress)
			body.SetAttributeRaw(name, replacedTokens)
		}
	}
}

func processSubaccountLevel(body *hclwrite.Body, variables *generictools.VariableContent, dependencyAddresses *generictools.DepedendcyAddresses, blockIdentifier string, resourceAddress string) {
	if blockIdentifier == subaccountBlockIdentifier {
		processSubaccountAttributes(body, variables)

		dependencyAddresses.SubaccountAddress = resourceAddress
	}

	if blockIdentifier != subaccountBlockIdentifier {
		replaceMainDependency(body, subaccountIdentifier, dependencyAddresses.SubaccountAddress)
	}

	if blockIdentifier == subaccountEntitlementBlockIdentifier {
		fillSubaccountEntitlementDependencyAddresses(body, resourceAddress, dependencyAddresses)
	}

	if blockIdentifier == subscriptionBlockIdentifier {
		addEntitlementDependency(body, dependencyAddresses)
	}
}

func processDirectoryLevel(body *hclwrite.Body, variables *generictools.VariableContent, dependencyAddresses *generictools.DepedendcyAddresses, blockIdentifier string, resourceAddress string) {
	if blockIdentifier == directoryBlockIdentifier {
		processDirectoryAttributes(body, variables)

		dependencyAddresses.DirectoryAddress = resourceAddress
	}

	if blockIdentifier != directoryBlockIdentifier {
		replaceMainDependency(body, directoryIdentifier, dependencyAddresses.DirectoryAddress)
	}
}
