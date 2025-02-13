package resourceprocessor

import (
	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

const cfOrgIdentifier = "org"

func extractOrgIds(body *hclwrite.Body, variables *generictools.VariableContent, orgId string) {
	generictools.ReplaceAttribute(body, cfOrgIdentifier, "ID of the Cloud Foundry Organization", variables)
}
