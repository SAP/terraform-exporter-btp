package resourceprocessor

import (
	cfcli "github.com/SAP/terraform-exporter-btp/internal/cfcli"
	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

const spaceBlockIdentifier = "cloudfoundry_space"
const spaceIdentifier = "space"
const spaceNameIdentifier = "name"

func fillSpaceDependencyAddress(body *hclwrite.Body, dependencyAddresses *generictools.DependencyAddresses, resourceAddress string) {
	spaceId := extractSpaceId(body)
	if spaceId == "" {
		return
	}
	dependencyAddresses.SpaceAddress[spaceId] = resourceAddress

}

func extractSpaceId(body *hclwrite.Body) string {
	spaceNameAttribute := body.GetAttribute(spaceNameIdentifier)
	if spaceNameAttribute == nil {
		return ""
	}

	spaceNameToken := spaceNameAttribute.Expr().BuildTokens(nil)

	spaceName := generictools.GetStringToken(spaceNameToken)
	if spaceName == "" {
		return ""
	}

	spaceId, _ := cfcli.GetSpaceId(spaceName)
	return spaceId
}

func replaceSpaceDependency(body *hclwrite.Body, spaceIdentifier string, spaceAddresses map[string]string) {
	for spaceId, spaceAddress := range spaceAddresses {
		if spaceAddress == "" {
			continue
		}

		for name, attr := range body.Attributes() {
			tokens := attr.Expr().BuildTokens(nil)

			if name == spaceIdentifier && len(tokens) == 3 && spaceId == generictools.GetStringToken(tokens) {
				replacedTokens := generictools.ReplaceDependency(tokens, spaceAddress)
				body.SetAttributeRaw(name, replacedTokens)
			}
		}
	}
}
