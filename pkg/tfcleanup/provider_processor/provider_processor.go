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
	if len(body.Attributes()) > 0 {
		generictools.ReplaceAttribute(body, "gGlobal account subdomain", globalAccountIdentifier, variables)
		generictools.ReplaceAttribute(body, "Cloud Foundry API endpoint", cfApiEndpointIdentifier, variables)
	}

	for _, block := range body.Blocks() {
		inBlocks := append(inBlocks, block.Type())
		processProviderAttributes(block.Body(), inBlocks, variables)
	}
}
