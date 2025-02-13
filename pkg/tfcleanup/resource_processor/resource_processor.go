package resourceprocessor

import (
	"github.com/SAP/terraform-exporter-btp/internal/btpcli"
	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/SAP/terraform-exporter-btp/pkg/tfutils"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func ProcessResources(hclFile *hclwrite.File, level string, variables *generictools.VariableContent, dependencyAddresses *generictools.DepedendcyAddresses, btpClient *btpcli.ClientFacade, levelIds generictools.LevelIds) {
	processResourceAttributes(hclFile.Body(), nil, level, variables, dependencyAddresses, btpClient, levelIds)
	processDependencies(hclFile.Body(), dependencyAddresses)
}

func processResourceAttributes(body *hclwrite.Body, inBlocks []string, level string, variables *generictools.VariableContent, dependencyAddresses *generictools.DepedendcyAddresses, btpClient *btpcli.ClientFacade, levelIds generictools.LevelIds) {
	if len(inBlocks) > 0 {

		generictools.RemoveEmptyAttributes(body)

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

	if blockIdentifier == trustConfigBlockIdentifier {
		processTrustConfigurationAttributes(body, blockIdentifier, resourceAddress, dependencyAddresses)
	}

	if blockIdentifier != subaccountBlockIdentifier {
		generictools.ReplaceMainDependency(body, subaccountIdentifier, dependencyAddresses.SubaccountAddress)
	}
}

func processDirectoryLevel(body *hclwrite.Body, variables *generictools.VariableContent, dependencyAddresses *generictools.DepedendcyAddresses, blockIdentifier string, resourceAddress string, btpClient *btpcli.ClientFacade, levelIds generictools.LevelIds) {
	if blockIdentifier == directoryBlockIdentifier {
		processDirectoryAttributes(body, variables, btpClient)
		dependencyAddresses.DirectoryAddress = resourceAddress
	}

	if blockIdentifier != directoryBlockIdentifier {
		generictools.ReplaceMainDependency(body, directoryIdentifier, dependencyAddresses.DirectoryAddress)
	}
}

func processCfOrgLevel(body *hclwrite.Body, variables *generictools.VariableContent, dependencyAddresses *generictools.DepedendcyAddresses, blockIdentifier string, resourceAddress string, btpClient *btpcli.ClientFacade, levelIds generictools.LevelIds) {
	extractOrgIds(body, variables, levelIds.CfOrgId)
}

func processDependencies(body *hclwrite.Body, dependencyAddresses *generictools.DepedendcyAddresses) {
	// Remove blocks that point to defaulted resources that get created by the platform automagically
	for _, blockToRemove := range dependencyAddresses.BlocksToRemove {
		generictools.RemoveConfigBlock(body, blockToRemove.BlockIdentifier, blockToRemove.ResourceAddress)
	}
	// Add datasource for service instances is necessary - Outer loop to have the main body object available
	for _, datasourceInfo := range dependencyAddresses.DataSourceInfo {
		addServicePlanDataSources(body, datasourceInfo)
	}
}
