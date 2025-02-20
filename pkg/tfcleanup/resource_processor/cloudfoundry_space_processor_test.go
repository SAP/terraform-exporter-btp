package resourceprocessor

import (
	"testing"

	"github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/testutils"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/assert"
)

func TestReplaceSpaceDependency(t *testing.T) {
	srcFileSpaceDep, trgtFileSpaceDep := testutils.GetHclFilesById("cf_space_dependency")
	srcFileSpaceNoDep, trgtFileSpaceNoDep := testutils.GetHclFilesById("cf_space_no_dependency")

	emptySpaceAddresses := make(map[string]string)

	tests := []struct {
		name           string
		src            *hclwrite.File
		trgt           *hclwrite.File
		spaceAddresses map[string]string
	}{
		{
			name:           "Test Space Dependency",
			src:            srcFileSpaceDep,
			trgt:           trgtFileSpaceDep,
			spaceAddresses: map[string]string{"1234567890": "cloudfoundry_space.space_dev-project-abc"},
		},
		{
			name:           "Test Space No Dependency",
			src:            srcFileSpaceNoDep,
			trgt:           trgtFileSpaceNoDep,
			spaceAddresses: emptySpaceAddresses,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			blocks := tt.src.Body().Blocks()
			// we assume one resource entry in the blocks file
			replaceSpaceDependency(blocks[0].Body(), spaceIdentifier, tt.spaceAddresses)
			assert.NoError(t, testutils.AreHclFilesEqual(tt.src, tt.trgt))
		})
	}
}
