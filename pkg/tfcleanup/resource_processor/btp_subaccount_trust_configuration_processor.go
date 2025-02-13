package resourceprocessor

import (
	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

const trustConfigBlockIdentifier = "btp_subaccount_trust_configuration"
const trustNameIdentifier = "name"
const trustDefaultIdentifier = "sap.default"

func processTrustConfigurationAttributes(body *hclwrite.Body, blockIdentifier string, resourceAddress string, dependencyAddresses *generictools.DepedendcyAddresses) {
	trustNameAttr := body.GetAttribute(trustNameIdentifier)
	if trustNameAttr == nil {
		return
	}

	identityProviderName := generictools.GetStringToken(trustNameAttr.Expr().BuildTokens(nil))

	if identityProviderName == trustDefaultIdentifier {
		identifier := generictools.BlockSpecifier{
			BlockIdentifier: blockIdentifier,
			ResourceAddress: resourceAddress,
		}

		(*dependencyAddresses).BlocksToRemove = append((*dependencyAddresses).BlocksToRemove, identifier)
	}
}
