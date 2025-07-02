package resourceprocessor

import (
	"github.com/SAP/terraform-exporter-btp/internal/btpcli"
	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/SAP/terraform-exporter-btp/pkg/tfutils"
	"github.com/SAP/terraform-exporter-btp/pkg/toggles"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func ProcessResources(hclFile *hclwrite.File, level string, variables *generictools.VariableContent, dependencyAddresses *generictools.DependencyAddresses, btpClient *btpcli.ClientFacade, levelIds generictools.LevelIds) {
	processResourceAttributes(hclFile.Body(), nil, level, variables, dependencyAddresses, btpClient, levelIds)
	processDependencies(hclFile.Body(), dependencyAddresses, variables)
}

func processResourceAttributes(body *hclwrite.Body, inBlocks []string, level string, variables *generictools.VariableContent, dependencyAddresses *generictools.DependencyAddresses, btpClient *btpcli.ClientFacade, levelIds generictools.LevelIds) {
	if len(inBlocks) > 0 {

		generictools.RemoveEmptyAttributes(body)

		_, blockIdentifier, resourceAddress := generictools.ExtractBlockInformation(inBlocks)

		switch level {
		case tfutils.SubaccountLevel:
			processSubaccountLevel(body, variables, dependencyAddresses, blockIdentifier, resourceAddress, btpClient, levelIds)
		case tfutils.DirectoryLevel:
			processDirectoryLevel(body, variables, dependencyAddresses, blockIdentifier, resourceAddress, btpClient)
		case tfutils.OrganizationLevel:
			processCfOrgLevel(body, variables, dependencyAddresses, blockIdentifier, resourceAddress)
		}
	}

	blocks := body.Blocks()
	for _, block := range blocks {
		inBlocks := append(inBlocks, block.Type()+","+block.Labels()[0]+","+block.Labels()[1])
		processResourceAttributes(block.Body(), inBlocks, level, variables, dependencyAddresses, btpClient, levelIds)
	}
}

func processSubaccountLevel(body *hclwrite.Body, variables *generictools.VariableContent, dependencyAddresses *generictools.DependencyAddresses, blockIdentifier string, resourceAddress string, btpClient *btpcli.ClientFacade, levelIds generictools.LevelIds) {
	if blockIdentifier == subaccountBlockIdentifier {
		processSubaccountAttributes(body, variables, btpClient)
		dependencyAddresses.SubaccountAddress = resourceAddress
	}

	if blockIdentifier == subaccountEntitlementBlockIdentifier {
		fillSubaccountEntitlementDependencyAddresses(body, resourceAddress, dependencyAddresses)
	}

	if blockIdentifier == subscriptionBlockIdentifier {
		addEntitlementDependency(body, dependencyAddresses, btpClient, levelIds.SubaccountId)
	}

	if blockIdentifier == serviceInstanceBlockIdentifier {
		addServiceInstanceDependency(body, dependencyAddresses, btpClient, levelIds.SubaccountId)
	}

	if blockIdentifier == subaccountRoleBlockIdentifier {
		fillRoleDependencyAddresses(body, resourceAddress, dependencyAddresses)
	}

	if blockIdentifier == subaccountRoleCollectionBlockIdentifier {
		addRoleDependency(body, dependencyAddresses)
	}

	if blockIdentifier != subaccountBlockIdentifier {
		generictools.ReplaceMainDependency(body, subaccountIdentifier, dependencyAddresses.SubaccountAddress)
	}
}

func processDirectoryLevel(body *hclwrite.Body, variables *generictools.VariableContent, dependencyAddresses *generictools.DependencyAddresses, blockIdentifier string, resourceAddress string, btpClient *btpcli.ClientFacade) {
	if blockIdentifier == directoryBlockIdentifier {
		processDirectoryAttributes(body, variables, btpClient)
		dependencyAddresses.DirectoryAddress = resourceAddress
	}

	if blockIdentifier == directoryRoleBlockIdentifier {
		fillRoleDependencyAddresses(body, resourceAddress, dependencyAddresses)
	}

	if blockIdentifier == directoryRoleCollectionBlockIdentifier {
		addRoleDependency(body, dependencyAddresses)
	}

	if blockIdentifier != directoryBlockIdentifier {
		generictools.ReplaceMainDependency(body, directoryIdentifier, dependencyAddresses.DirectoryAddress)
	}
}

func processCfOrgLevel(body *hclwrite.Body, variables *generictools.VariableContent, dependencyAddresses *generictools.DependencyAddresses, blockIdentifier string, resourceAddress string) {
	extractOrgIds(body, variables)
	if blockIdentifier == spaceBlockIdentifier {
		fillSpaceDependencyAddress(body, dependencyAddresses, resourceAddress)
	}

	if blockIdentifier != spaceBlockIdentifier {
		replaceSpaceDependency(body, spaceIdentifier, dependencyAddresses.SpaceAddress)
	}
}

func processDependencies(body *hclwrite.Body, dependencyAddresses *generictools.DependencyAddresses, variables *generictools.VariableContent) {
	// Remove blocks that point to defaulted resources that get created by the platform automagically
	for _, blockToRemove := range dependencyAddresses.BlocksToRemove {
		generictools.RemoveConfigBlock(body, blockToRemove.ResourceAddress)
	}

	// Add generic module for entitlements
	if !toggles.IsEntitlementModuleGenerationDeactivated() {

		// Add module for entitlements
		addEntitlementModule(body, dependencyAddresses.SubaccountAddress)

		// Add variables to variables to be generated
		addEntitlementVariables(variables, dependencyAddresses)

		// Remove the entitlement configurations
		for _, entitlementToRemove := range dependencyAddresses.EntitlementAddress {
			generictools.RemoveConfigBlock(body, entitlementToRemove)
		}
	}

	// Add datasource for service instances is necessary - Outer loop to have the main body object available
	processedDataSources := make(map[string]bool)

	for _, datasourceInfo := range dependencyAddresses.DataSourceInfo {

		if !processedDataSources[datasourceInfo.DatasourceAddress] {
			addServicePlanDataSources(body, datasourceInfo)
			// Avoid duplicate data sources
			processedDataSources[datasourceInfo.DatasourceAddress] = true
		}
	}
}
