package resourceprocessor

import (
	"testing"

	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/testutils"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/assert"
)

func TestProcessSATrustConfig(t *testing.T) {

	srcFileReplace, _ := testutils.GetHclFilesById("sa_trust_config_replace")
	srcFileNoReplace, _ := testutils.GetHclFilesById("sa_trust_config_no_replace")

	emptyTestDependencies := generictools.NewDependencyAddresses()

	defaultTrustRemoveBlock := generictools.BlockSpecifier{
		BlockIdentifier: "btp_subaccount_trust_configuration",
		ResourceAddress: "btp_subaccount_trust_configuration.trust_0",
	}

	defaultTestDependencies := generictools.NewDependencyAddresses()
	defaultTestDependencies.BlocksToRemove = append(defaultTestDependencies.BlocksToRemove, defaultTrustRemoveBlock)

	tests := []struct {
		name             string
		src              *hclwrite.File
		trgtDependencies *generictools.DepedendcyAddresses
		blockIdentifier  string
		resourceAddress  string
	}{
		{
			name:             "Test Trust Configuration Cleanup",
			src:              srcFileReplace,
			trgtDependencies: &defaultTestDependencies,
			blockIdentifier:  trustConfigBlockIdentifier,
			resourceAddress:  "btp_subaccount_trust_configuration.trust_0",
		},
		{
			name:             "Test BTP Provider Cleanup - No changes",
			src:              srcFileNoReplace,
			trgtDependencies: &emptyTestDependencies,
			blockIdentifier:  trustConfigBlockIdentifier,
			resourceAddress:  "btp_subaccount_trust_configuration.trust_1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			dependencies := generictools.NewDependencyAddresses()

			blocks := tt.src.Body().Blocks()
			// we assume one rrsource entry in the blocks file
			processTrustConfigurationAttributes(blocks[0].Body(), blocks[0].Labels()[0], blocks[0].Labels()[0]+"."+blocks[0].Labels()[1], &dependencies)
			assert.Equal(t, tt.trgtDependencies.BlocksToRemove, dependencies.BlocksToRemove)
		})
	}
}
