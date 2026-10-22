package resourceprocessor

import (
	"testing"

	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/assert"
)

func TestProcessDirectoryLevelRestoresMissingDirectoryID(t *testing.T) {
	f := hclwrite.NewEmptyFile()
	block := f.Body().AppendNewBlock("resource", []string{directoryRoleBlockIdentifier, "role"})

	dependencies := generictools.NewDependencyAddresses()
	processDirectoryLevel(block.Body(), &generictools.VariableContent{}, &dependencies, directoryRoleBlockIdentifier, directoryRoleBlockIdentifier+".role", nil, "directory-guid")

	assert.Contains(t, string(f.Bytes()), "directory_id = \"directory-guid\"")
}
