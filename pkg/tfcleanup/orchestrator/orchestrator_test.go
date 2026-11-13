package orchestrator

import (
	"os"
	"path/filepath"
	"testing"

	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/SAP/terraform-exporter-btp/pkg/tfutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImportsOnlyCleanup(t *testing.T) {
	for _, level := range []string{tfutils.SubaccountLevel, tfutils.DirectoryLevel, tfutils.OrganizationLevel} {
		t.Run(level, func(t *testing.T) {
			dir := t.TempDir()
			resources := `resource "btp_subaccount" "export_0" {
  parent_id = "parent-id"
  description = null
}
`
			provider := `terraform {
  required_providers {
    btp = { source = "SAP/btp" }
  }
}
provider "btp" {
  globalaccount = "global-account"
}
`
			imports := "import {\n  to = btp_subaccount.export_0\n  id = \"account-id\"\n}\n"
			for name, content := range map[string]string{"btp_resources.tf": resources, "provider.tf": provider, "btp_subaccount_import.tf": imports} {
				require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0600))
			}
			counts := map[string]int{tfutils.SubaccountType: 1}
			backend := tfutils.BackendConfig{BackendType: "local", BackendConfig: []string{"path=state.tfstate"}}
			require.NoError(t, orchestrateCodeCleanup(dir, level, generictools.LevelIds{}, &counts, backend, true))
			// Resource refinements would remove description and query the parent.
			unchanged, err := os.ReadFile(filepath.Join(dir, "btp_resources.tf"))
			require.NoError(t, err)
			assert.Equal(t, resources, string(unchanged))
			unchanged, err = os.ReadFile(filepath.Join(dir, "btp_subaccount_import.tf"))
			require.NoError(t, err)
			assert.Equal(t, imports, string(unchanged))
			assert.Equal(t, 1, counts[tfutils.SubaccountType])
			actualProvider, err := os.ReadFile(filepath.Join(dir, "provider.tf"))
			require.NoError(t, err)
			assert.Contains(t, string(actualProvider), "globalaccount = var.globalaccount")
			assert.Contains(t, string(actualProvider), `backend "local"`)
			assert.Contains(t, string(actualProvider), `path = "state.tfstate"`)
			variables, err := os.ReadFile(filepath.Join(dir, "variables.tf"))
			require.NoError(t, err)
			assert.Contains(t, string(variables), `variable "globalaccount"`)
		})
	}
}
