package resourceprocessor

import (
	"github.com/SAP/terraform-exporter-btp/internal/btpcli"
	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/SAP/terraform-exporter-btp/pkg/tfutils"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func ProcessResources(hclFile *hclwrite.File, level string, variables *generictools.VariableContent, dependencyAddresses *generictools.DepedendcyAddresses, btpClient *btpcli.ClientFacade, levelIds generictools.LevelIds) {

	processResourceAttributes(hclFile.Body(), nil, level, variables, dependencyAddresses, btpClient, levelIds)
	// Add datasource for service instances is necessary - Outer loop to have the main body object available
	for _, datasourceInfo := range dependencyAddresses.DataSourceInfo {
		addServicePlanDataSources(hclFile.Body(), datasourceInfo)
	}
}

func processResourceAttributes(body *hclwrite.Body, inBlocks []string, level string, variables *generictools.VariableContent, dependencyAddresses *generictools.DepedendcyAddresses, btpClient *btpcli.ClientFacade, levelIds generictools.LevelIds) {

	if len(inBlocks) > 0 {

		removeEmptyAttributes(body)

		_, blockIdentifier, resourceAddress := generictools.ExtractBlockInformation(inBlocks)

		switch level {
		case tfutils.SubaccountLevel:
			processSubaccountLevel(body, variables, dependencyAddresses, blockIdentifier, resourceAddress, btpClient, levelIds)
		case tfutils.DirectoryLevel:
			processDirectoryLevel(body, variables, dependencyAddresses, blockIdentifier, resourceAddress, btpClient, levelIds)
		case tfutils.OrganizationLevel:
			processCfOrgLevel(body, variables, dependencyAddresses, blockIdentifier, resourceAddress, btpClient, levelIds)
		}
	}

	blocks := body.Blocks()
	for _, block := range blocks {
		inBlocks := append(inBlocks, block.Type()+","+block.Labels()[0]+","+block.Labels()[1])
		processResourceAttributes(block.Body(), inBlocks, level, variables, dependencyAddresses, btpClient, levelIds)
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

func processSubaccountLevel(body *hclwrite.Body, variables *generictools.VariableContent, dependencyAddresses *generictools.DepedendcyAddresses, blockIdentifier string, resourceAddress string, btpClient *btpcli.ClientFacade, levelIds generictools.LevelIds) {
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

	if blockIdentifier == serviceInstanceBlockIdentifier {
		addServiceInstanceDependency(body, dependencyAddresses, btpClient, levelIds.SubaccountId)
	}

	if blockIdentifier != subaccountBlockIdentifier {
		replaceMainDependency(body, subaccountIdentifier, dependencyAddresses.SubaccountAddress)
	}
}

func processDirectoryLevel(body *hclwrite.Body, variables *generictools.VariableContent, dependencyAddresses *generictools.DepedendcyAddresses, blockIdentifier string, resourceAddress string, btpClient *btpcli.ClientFacade, levelIds generictools.LevelIds) {
	if blockIdentifier == directoryBlockIdentifier {
		processDirectoryAttributes(body, variables, btpClient)

		dependencyAddresses.DirectoryAddress = resourceAddress
	}

	if blockIdentifier != directoryBlockIdentifier {
		replaceMainDependency(body, directoryIdentifier, dependencyAddresses.DirectoryAddress)
	}
}

func processCfOrgLevel(body *hclwrite.Body, variables *generictools.VariableContent, dependencyAddresses *generictools.DepedendcyAddresses, blockIdentifier string, resourceAddress string, btpClient *btpcli.ClientFacade, levelIds generictools.LevelIds) {
	extractOrgIds(body, variables, levelIds.CfOrgId)
}
