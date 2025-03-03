package providerprocessor

import (
	"testing"

	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/testutils"
	"github.com/SAP/terraform-exporter-btp/pkg/tfutils"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/assert"
)

func TestProcessProvider(t *testing.T) {

	btpSrcFile, btpTrgtFile := testutils.GetHclFilesById("provider_btp")
	cfSrcFile, cfTrgtFile := testutils.GetHclFilesById("provider_cf")

	emptyTestContent := make(generictools.VariableContent)

	tests := []struct {
		name          string
		src           *hclwrite.File
		trgt          *hclwrite.File
		trgtVariables *generictools.VariableContent
	}{
		{
			name: "Test BTP Provider Cleanup",
			src:  btpSrcFile,
			trgt: btpTrgtFile,
			trgtVariables: &generictools.VariableContent{
				"globalaccount": generictools.VariableInfo{
					Description: "Global account subdomain",
					Value:       "my-global-account",
				},
			},
		},
		{
			name: "Test CF Provider Cleanup",
			src:  cfSrcFile,
			trgt: cfTrgtFile,
			trgtVariables: &generictools.VariableContent{
				"api_url": generictools.VariableInfo{
					Description: "Cloud Foundry API endpoint",
					Value:       "https://api.cf.sap.hana.ondemand.com",
				},
			},
		},
		{
			name:          "Test BTP Provider Cleanup - No changes",
			src:           btpTrgtFile,
			trgt:          btpTrgtFile,
			trgtVariables: &emptyTestContent,
		},
		{
			name:          "Test CF Provider Cleanup - No changes",
			src:           cfTrgtFile,
			trgt:          cfTrgtFile,
			trgtVariables: &emptyTestContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			contentToCreate := make(generictools.VariableContent)
			backendConfig := tfutils.BackendConfig{
				PathToBackendConfig: "",
				BackendType:         "",
				BackendConfig:       []string{},
			}
			ProcessProvider(tt.src, &contentToCreate, backendConfig)

			assert.NoError(t, testutils.AreHclFilesEqual(tt.trgt, tt.src))
			assert.Equal(t, tt.trgtVariables, &contentToCreate)

		})
	}
}
