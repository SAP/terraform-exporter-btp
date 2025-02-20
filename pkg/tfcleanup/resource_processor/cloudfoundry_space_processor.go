package resourceprocessor

import (
	cfcli "github.com/SAP/terraform-exporter-btp/internal/cfcli"
	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

const spaceBlockIdentifier = "cloudfoundry_space"
const spaceIdentifier = "space"
const spaceNameIdentifier = "name"

func ExtractSpaceId(body *hclwrite.Body) string {
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
