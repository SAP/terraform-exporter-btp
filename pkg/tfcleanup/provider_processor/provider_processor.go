package providerprocessor

import (
	"strings"

	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/SAP/terraform-exporter-btp/pkg/tfutils"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

const globalAccountIdentifier = "globalaccount"
const cfApiEndpointIdentifier = "api_url"

func ProcessProvider(hclFile *hclwrite.File, variables *generictools.VariableContent, backendConfig tfutils.BackendConfig) {
	processProviderAttributes(hclFile.Body(), nil, variables)
	addBackendBlock(hclFile.Body(), backendConfig)
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

func addBackendBlock(body *hclwrite.Body, backendConfig tfutils.BackendConfig) {

	if backendConfig.PathToBackendConfig != "" {
		appendBackendBlockByFile(backendConfig.PathToBackendConfig, body)
		return
	}

	if backendConfig.BackendType != "" && backendConfig.BackendConfig != nil {
		appendBackendBlockByParams(backendConfig.BackendType, backendConfig.BackendConfig, body)
		return
	}

}

func appendBackendBlockByFile(filePath string, body *hclwrite.Body) {
	backendFile := generictools.GetHclFile(filePath)
	backendBody := backendFile.Body()

	terraformBlock := backendBody.Blocks()[0]
	backendBlock := terraformBlock.Body().Blocks()[0]

	if backendBlock.Type() == "backend" {
		body.AppendNewline()
		body.AppendBlock(backendBlock)
	}
}

func appendBackendBlockByParams(backendType string, backendConfig []string, body *hclwrite.Body) {
	body.AppendNewline()
	backendBlock := body.AppendNewBlock("backend", []string{backendType})

	backendBody := backendBlock.Body()

	for _, config := range backendConfig {
		configParts := strings.Split(config, "=")

		backendBody.SetAttributeRaw(configParts[0],
			hclwrite.Tokens{
				{
					Type:  hclsyntax.TokenIdent,
					Bytes: []byte("\"" + configParts[1] + "\""),
				},
			},
		)
	}
}
