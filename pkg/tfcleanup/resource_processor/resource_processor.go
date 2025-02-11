package resourceprocessor

import (
	"log"

	"github.com/SAP/terraform-exporter-btp/internal/btpcli"
	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/SAP/terraform-exporter-btp/pkg/tfutils"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func ProcessResources(hclFile *hclwrite.File, level string, variables *generictools.VariableContent, dependencyAddresses *generictools.DepedendcyAddresses, btpClient *btpcli.ClientFacade) {

	processResourceAttributes(hclFile.Body(), nil, level, variables, dependencyAddresses, btpClient)
}

func processResourceAttributes(body *hclwrite.Body, inBlocks []string, level string, variables *generictools.VariableContent, dependencyAddresses *generictools.DepedendcyAddresses, btpClient *btpcli.ClientFacade) {

	if len(inBlocks) > 0 {

		removeEmptyAttributes(body)

		blockIdentifier, resourceAddress := generictools.ExtractBlockInformation(inBlocks)

		switch level {
		case tfutils.SubaccountLevel:
			processSubaccountLevel(body, variables, dependencyAddresses, blockIdentifier, resourceAddress, btpClient)
		case tfutils.DirectoryLevel:
			processDirectoryLevel(body, variables, dependencyAddresses, blockIdentifier, resourceAddress, btpClient)
		case tfutils.OrganizationLevel:
			log.Println("Organization level is not supported yet")
		}
	}

	blocks := body.Blocks()
	for _, block := range blocks {
		inBlocks := append(inBlocks, block.Type()+","+block.Labels()[0]+","+block.Labels()[1])
		processResourceAttributes(block.Body(), inBlocks, level, variables, dependencyAddresses, btpClient)
	}
}

func removeEmptyAttributes(body *hclwrite.Body) {
	attrs := body.Attributes()
	for name, attr := range attrs {
		tokens := attr.Expr().BuildTokens(nil)

		// Check for a NULL value
		if len(tokens) == 1 && string(tokens[0].Bytes) == generictools.EmptyString {
			body.RemoveAttribute(name)
		}

		// Check for an empty JSON encoded string or an empty Map
		var combinedString string
		if len(tokens) == 5 || len(tokens) == 2 {
			for _, token := range tokens {
				combinedString += string(token.Bytes)
			}
		}

		if combinedString == generictools.EmptyJson || combinedString == generictools.EmptyMap {
			body.RemoveAttribute(name)
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

func processSubaccountLevel(body *hclwrite.Body, variables *generictools.VariableContent, dependencyAddresses *generictools.DepedendcyAddresses, blockIdentifier string, resourceAddress string, btpClient *btpcli.ClientFacade) {
	if blockIdentifier == subaccountBlockIdentifier {
		processSubaccountAttributes(body, variables, btpClient)

		dependencyAddresses.SubaccountAddress = resourceAddress
	}

	if blockIdentifier == subaccountEntitlementBlockIdentifier {
		fillSubaccountEntitlementDependencyAddresses(body, resourceAddress, dependencyAddresses)
	}

	if blockIdentifier == subscriptionBlockIdentifier {
		addEntitlementDependency(body, dependencyAddresses)
	}

	// We add the reference to the subaccount at the end to have the subaccount ID available
	if blockIdentifier != subaccountBlockIdentifier {
		replaceMainDependency(body, subaccountIdentifier, dependencyAddresses.SubaccountAddress)
	}
}

func processDirectoryLevel(body *hclwrite.Body, variables *generictools.VariableContent, dependencyAddresses *generictools.DepedendcyAddresses, blockIdentifier string, resourceAddress string, btpClient *btpcli.ClientFacade) {
	if blockIdentifier == directoryBlockIdentifier {
		processDirectoryAttributes(body, variables, btpClient)

		dependencyAddresses.DirectoryAddress = resourceAddress
	}

	if blockIdentifier != directoryBlockIdentifier {
		replaceMainDependency(body, directoryIdentifier, dependencyAddresses.DirectoryAddress)
	}
}
