package tfutils

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFinalizeImportsOnly(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Terraform stub uses a POSIX shell")
	}
	for _, importsOnly := range []bool{false, true} {
		name := "full export"
		if importsOnly {
			name = "imports only"
		}
		t.Run(name, func(t *testing.T) {
			// Output selection must work independently of the cleanup toggle.
			t.Setenv("BTPTF_SKIP_CODECLEANUP", "1")
			root := t.TempDir()
			t.Chdir(root)
			dir := filepath.Join(root, "export")
			require.NoError(t, os.Mkdir(dir, 0700))
			fixtures := map[string]string{
				"btp_resources.tf": `resource "btp_subaccount" "export_0" {
  subaccount_id = "account-id"
}
module "btp_subaccount_entitlement" {
  source = "./entitlements"
}
`,
				"imports.tf": `import {
  to = btp_subaccount.export_0
  identity = { subaccount_id = "account-id" }
}
import {
  to = module.btp_subaccount_entitlement.btp_subaccount_entitlement.entitlement["service-plan"]
  id = "account-id,service,plan"
}
`,
				"provider.tf": `terraform {
  required_providers {
    btp = { source = "SAP/btp" }
  }
  backend "local" {}
}
provider "btp" {
  globalaccount = var.globalaccount
}
`,
				"variables.tf": `variable "globalaccount" {
  type = string
  default = "global-account"
}
`,
				"terraform.tfvars": "subaccount_id = \"account-id\"\n",
				"NextSteps.md":     "Stale instructions from an interrupted export",
			}
			for filename, content := range fixtures {
				require.NoError(t, os.WriteFile(filepath.Join(dir, filename), []byte(content), 0600))
			}
			// Stand in for Terraform; record whether resource configuration is
			// still present when initialization starts, without downloading providers.
			script := filepath.Join(root, "terraform-stub")
			require.NoError(t, os.WriteFile(script, []byte(`#!/bin/sh
if [ "$1" = init ] && [ -f btp_resources.tf ]; then
  touch resources-present-at-init
fi
printf '%s\n' "$*" >> commands.log
`), 0700))
			originalTool := getIaCTool
			getIaCTool = func() (string, error) { return script, nil }
			t.Cleanup(func() { getIaCTool = originalTool })

			FinalizeTfConfig("export", ExecutionTarget{SubaccountID: "account-id"}, importsOnly)

			for _, filename := range []string{"btp_resources.tf", "NextSteps.md", "resources-present-at-init"} {
				if importsOnly {
					assert.NoFileExists(t, filepath.Join(dir, filename))
				} else {
					assert.FileExists(t, filepath.Join(dir, filename))
				}
			}
			imports, err := os.ReadFile(filepath.Join(dir, "imports.tf"))
			require.NoError(t, err)
			assert.Contains(t, string(imports), "btp_subaccount.subaccount_0")
			assert.Contains(t, string(imports), "subaccount_id = var.subaccount_id")
			assert.Contains(t, string(imports), `module.btp_subaccount_entitlement.btp_subaccount_entitlement.entitlement["service-plan"]`)
			for _, filename := range []string{"provider.tf", "terraform.tfvars"} {
				content, err := os.ReadFile(filepath.Join(dir, filename))
				require.NoError(t, err)
				assert.Equal(t, fixtures[filename], string(content))
			}
			variables, err := os.ReadFile(filepath.Join(dir, "variables.tf"))
			require.NoError(t, err)
			assert.Contains(t, string(variables), `variable "globalaccount"`)
			assert.Contains(t, string(variables), `variable "subaccount_id"`)
			commands, err := os.ReadFile(filepath.Join(dir, "commands.log"))
			require.NoError(t, err)
			assert.Contains(t, string(commands), "init -backend=false")
			assert.Contains(t, string(commands), "fmt -recursive -list=false")
		})
	}
}

func TestImportsOnlyPostQuerySkipsTerraformFormatting(t *testing.T) {
	t.Chdir(t.TempDir())
	require.NoError(t, os.Mkdir("query", 0700))
	require.NoError(t, os.Mkdir("export", 0700))
	resources := `resource "btp_subaccount" "export_0" {}`
	imports := "import {\n  to = btp_subaccount.export_0\n  id = \"account-id\"\n}\n"
	require.NoError(t, os.WriteFile(filepath.Join("query", "btp_resources.tf"), []byte(resources), 0600))
	require.NoError(t, os.WriteFile(filepath.Join("query", "btp_subaccount_import.tf"), []byte(imports), 0600))
	originalTool := getIaCTool
	getIaCTool = func() (string, error) {
		t.Fatal("imports-only post-query processing must not run Terraform fmt")
		return "", nil
	}
	t.Cleanup(func() { getIaCTool = originalTool })

	ExecPostQuerySteps("query", "export", "btp_resources.tf", SubaccountType, true)

	assert.NoDirExists(t, "query")
	merged, err := os.ReadFile(filepath.Join("export", "btp_subaccount_import.tf"))
	require.NoError(t, err)
	assert.Equal(t, imports, string(merged))
	// Raw resources remain available for import transformations and resume
	// until finalization removes them.
	merged, err = os.ReadFile(filepath.Join("export", "btp_resources.tf"))
	require.NoError(t, err)
	assert.Contains(t, string(merged), resources)
}

func TestRemoveResourceConfiguration(t *testing.T) {
	t.Run("empty or already finalized export", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, removeResourceConfiguration(dir))
		require.NoError(t, removeResourceConfiguration(dir))
	})
	t.Run("removal failure is returned", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "btp_resources.tf")
		require.NoError(t, os.Mkdir(path, 0700))
		require.NoError(t, os.WriteFile(filepath.Join(path, "keep"), nil, 0600))
		require.ErrorContains(t, removeResourceConfiguration(dir), "remove btp_resources.tf")
	})
}
