package resourceprocessor

import (
	"github.com/SAP/terraform-exporter-btp/internal/btpcli"
	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

const directoryBlockIdentifier = "btp_directory"
const directoryIdentifier = "directory_id"

func processDirectoryAttributes(body *hclwrite.Body, variables *generictools.VariableContent, btpClient *btpcli.ClientFacade) {
	generictools.ProcessParentAttribute(body, "ID of the parent of the SAP BTP directory", btpClient, variables)
}
