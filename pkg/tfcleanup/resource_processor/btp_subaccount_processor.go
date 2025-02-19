package resourceprocessor

import (
	"github.com/SAP/terraform-exporter-btp/internal/btpcli"
	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

const subaccountBlockIdentifier = "btp_subaccount"
const subaccountIdentifier = "subaccount_id"
const regionIdentifier = "region"

func processSubaccountAttributes(body *hclwrite.Body, variables *generictools.VariableContent, btpClient *btpcli.ClientFacade) {
	generictools.ReplaceAttribute(body, regionIdentifier, "Region of SAP BTP subaccount", variables)
	generictools.ProcessParentAttribute(body, "ID of the parent of the SAP BTP subaccount", btpClient, variables)
}

