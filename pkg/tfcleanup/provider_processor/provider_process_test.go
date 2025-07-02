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

	emptyBackendConfig := tfutils.BackendConfig{
		PathToBackendConfig: "",
		BackendType:         "",
		BackendConfig:       []string{},
	}

	tests := []struct {
		name          string
		src           *hclwrite.File
		trgt          *hclwrite.File
		trgtVariables *generictools.VariableContent
		backendConfig tfutils.BackendConfig
	}{
		{
			name: "Test BTP Provider Cleanup",
			src:  btpSrcFile,
			trgt: btpTrgtFile,
			trgtVariables: &generictools.VariableContent{
				"globalaccount": generictools.VariableInfo{
					Description:  "Global account subdomain",
					Type:         "string",
					DefaultValue: "my-global-account",
				},
			},
		},
		{
			name: "Test CF Provider Cleanup",
			src:  cfSrcFile,
			trgt: cfTrgtFile,
			trgtVariables: &generictools.VariableContent{
				"api_url": generictools.VariableInfo{
					Description:  "Cloud Foundry API endpoint",
					Type:         "string",
					DefaultValue: "https://api.cf.sap.hana.ondemand.com",
				},
			},
			backendConfig: emptyBackendConfig,
		},
		{
			name:          "Test BTP Provider Cleanup - No changes",
			src:           btpTrgtFile,
			trgt:          btpTrgtFile,
			trgtVariables: &emptyTestContent,
			backendConfig: emptyBackendConfig,
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

func TestProcessProviderWithBackend(t *testing.T) {

	btpSrcFileBackend, btpTrgtFileBackend := testutils.GetHclFilesById("provider_withbackend_btp")
	btpSrcFileBackend2, btpTrgtFileBackend2 := testutils.GetHclFilesById("provider_withbackend_btp")
	btpSrcFileWoBackend, btpTrgtFileWoBackend := testutils.GetHclFilesById("provider_wobackend_btp")

	emptyTestContent := make(generictools.VariableContent)

	emptyBackendConfig := tfutils.BackendConfig{
		PathToBackendConfig: "",
		BackendType:         "",
		BackendConfig:       []string{},
	}

	tests := []struct {
		name          string
		src           *hclwrite.File
		trgt          *hclwrite.File
		trgtVariables *generictools.VariableContent
		backendConfig tfutils.BackendConfig
	}{
		{
			name:          "Test BTP Provider no backend",
			src:           btpSrcFileWoBackend,
			trgt:          btpTrgtFileWoBackend,
			trgtVariables: &emptyTestContent,
			backendConfig: emptyBackendConfig,
		},
		{
			name:          "Test BTP Provider with backend file",
			src:           btpSrcFileBackend,
			trgt:          btpTrgtFileBackend,
			trgtVariables: &emptyTestContent,
			backendConfig: tfutils.BackendConfig{
				PathToBackendConfig: "../testutils/testdata/backend.tf",
				BackendType:         "",
				BackendConfig:       []string{},
			},
		},
		{
			name:          "Test BTP Provider with backend params",
			src:           btpSrcFileBackend2,
			trgt:          btpTrgtFileBackend2,
			trgtVariables: &emptyTestContent,
			backendConfig: tfutils.BackendConfig{
				PathToBackendConfig: "",
				BackendType:         "azurerm",
				BackendConfig:       []string{"resource_group_name=rg-terraform-state", "storage_account_name=terraformstatestorage", "container_name=tfstate", "key=terraform.tfstate"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			contentToCreate := make(generictools.VariableContent)

			ProcessProvider(tt.src, &contentToCreate, tt.backendConfig)

			assert.NoError(t, testutils.AreHclFilesEqual(tt.trgt, tt.src))

		})
	}
}
