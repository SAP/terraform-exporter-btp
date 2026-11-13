package resourceprocessor

import (
	"testing"

	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/testutils"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollectEntitlementImportsWithoutResourceRefinements(t *testing.T) {
	source := []byte(`resource "btp_subaccount_entitlement" "export_0" {
  service_name = "service"
  plan_name = "plan"
  amount = 2
  plan_unique_identifier = null
}
resource "btp_subaccount_subscription" "export_0" {
  app_name = "app"
}
`)
	file, diagnostics := hclwrite.ParseConfig(source, "resources.tf", hcl.InitialPos)
	require.False(t, diagnostics.HasErrors(), diagnostics.Error())
	original := string(file.Bytes())
	dependencies := generictools.NewDependencyAddresses()
	CollectEntitlementImports(file, &dependencies)
	assert.Equal(t, original, string(file.Bytes()), "resource configuration must not be refined")
	require.Len(t, dependencies.EntitlementAddress, 1)
	require.Len(t, dependencies.BlocksToRemove, 1)
	assert.Equal(t, "btp_subaccount_entitlement.export_0", dependencies.BlocksToRemove[0].ResourceAddress)

	imports, diagnostics := hclwrite.ParseConfig([]byte(`import {
  to = btp_subaccount_entitlement.export_0
  identity = { subaccount_id = "account-id", service_name = "service", plan_name = "plan" }
}
`), "imports.tf", hcl.InitialPos)
	require.False(t, diagnostics.HasErrors(), diagnostics.Error())
	counts := map[string]int{subaccountEntitlementBlockIdentifier: 1}
	for _, block := range dependencies.BlocksToRemove {
		generictools.RemoveImportBlock(imports.Body(), block.ResourceAddress, &counts)
	}
	appendEntitlementModuleImports(imports.Body(), dependencies.EntitlementAddress, generictools.LevelIds{SubaccountId: "account-id"})
	require.Len(t, imports.Body().Blocks(), 1)
	assert.Equal(t, "import", imports.Body().Blocks()[0].Type())
	assert.Contains(t, string(imports.Bytes()), `module.btp_subaccount_entitlement.btp_subaccount_entitlement.entitlement["service-plan"]`)
	assert.Contains(t, string(imports.Bytes()), `"account-id,service,plan"`)
	assert.Equal(t, 1, counts[subaccountEntitlementBlockIdentifier])
}

func TestFillSubaccountEntitlementDependencyAddresses(t *testing.T) {

	srcFileEntitlement, _ := testutils.GetHclFilesById("sa_entitlement")
	srcFileEntitlementIncomplete, _ := testutils.GetHclFilesById("sa_entitlement_error")

	emptyTestDependencies := generictools.NewDependencyAddresses()
	defaultTestDependencies := generictools.NewDependencyAddresses()

	defaultTestDependencies.EntitlementAddress = make(map[generictools.EntitlementKey]generictools.EntitlementInfo)

	defaultTestDependencies.EntitlementAddress[generictools.EntitlementKey{
		ServiceName: "feature-flags-dashboard",
		PlanName:    "dashboard",
	}] = generictools.EntitlementInfo{Address: "btp_subaccount_entitlement.entitlement_2", Amount: 1}

	tests := []struct {
		name             string
		src              *hclwrite.File
		resourceAddress  string
		trgtDependencies *generictools.DependencyAddresses
	}{
		{
			name:             "Test Entitlement Dependency Address",
			src:              srcFileEntitlement,
			resourceAddress:  "btp_subaccount_entitlement.entitlement_2",
			trgtDependencies: &defaultTestDependencies,
		},
		{
			name:             "Test Entitlement Incomplete",
			src:              srcFileEntitlementIncomplete,
			resourceAddress:  "btp_subaccount_entitlement.entitlement_2",
			trgtDependencies: &emptyTestDependencies,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			dependencies := generictools.NewDependencyAddresses()
			blocks := tt.src.Body().Blocks()
			// we assume one resource entry in the blocks file
			fillSubaccountEntitlementDependencyAddresses(blocks[0].Body(), tt.resourceAddress, &dependencies)

			assert.Equal(t, tt.trgtDependencies.EntitlementAddress, dependencies.EntitlementAddress)
		})
	}
}
