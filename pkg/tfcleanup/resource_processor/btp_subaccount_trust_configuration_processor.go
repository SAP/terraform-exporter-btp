package resourceprocessor

import (
	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

const trustConfigBlockIdentifier = "btp_subaccount_trust_configuration"
const trustNameIdentifier = "name"
const trustDefaultIdentifier = "sap.default"

func processTrustConfigurationAttributes(body *hclwrite.Body, blockIdentifier string, resourceAddress string, dependencyAddresses *generictools.DepedendcyAddresses) {
	removeEntry := false
	attrs := body.Attributes()
	for name, attr := range attrs {
		tokens := attr.Expr().BuildTokens(nil)

		if name == trustNameIdentifier {
			identityProviderName := generictools.GetStringToken(tokens)
			if identityProviderName == trustDefaultIdentifier {
				removeEntry = true
				break
			}
		}
	}

	if removeEntry {
		identifier := generictools.BlockSpecifier{
			BlockIdentifier: blockIdentifier,
			ResourceAddress: resourceAddress,
		}

		(*dependencyAddresses).BlocksToRemove = append((*dependencyAddresses).BlocksToRemove, identifier)
	}
}
