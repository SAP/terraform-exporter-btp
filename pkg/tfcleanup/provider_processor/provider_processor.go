package providerprocessor

import (
	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

const globalAccountIdentifier = "globalaccount"
const cfApiEndpointIdentifier = "api_url"

func ProcessProvider(hclFile *hclwrite.File, variables *generictools.VariableContent) {
	processProviderAttributes(hclFile.Body(), nil, variables)
}

func processProviderAttributes(body *hclwrite.Body, inBlocks []string, variables *generictools.VariableContent) {
	attributes := body.Attributes()

	if len(attributes) > 0 {
		generictools.ReplaceAttribute(body, globalAccountIdentifier, "Global account subdomain", variables)
		generictools.ReplaceAttribute(body, cfApiEndpointIdentifier, "Cloud Foundry API endpoint", variables)
	}

	for _, block := range body.Blocks() {
		inBlocks := append(inBlocks, block.Type())
		processProviderAttributes(block.Body(), inBlocks, variables)
	}
}
