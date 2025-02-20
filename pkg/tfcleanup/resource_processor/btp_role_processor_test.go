package resourceprocessor

import (
	"testing"

	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/testutils"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/assert"
)

func TestFillSubaccountEntitlementDependencyAddresses(t *testing.T) {

	srcFileEntitlement, _ := testutils.GetHclFilesById("sa_entitlement")
	srcFileEntitlementIncomplete, _ := testutils.GetHclFilesById("sa_entitlement_error")

	emptyTestDependencies := generictools.NewDependencyAddresses()
	defaultTestDependencies := generictools.NewDependencyAddresses()

	defaultTestDependencies.EntitlementAddress = make(map[generictools.EntitlementKey]string)

	defaultTestDependencies.EntitlementAddress[generictools.EntitlementKey{
		ServiceName: "feature-flags-dashboard",
		PlanName:    "dashboard",
	}] = "btp_subaccount_entitlement.entitlement_2"

	tests := []struct {
		name             string
		src              *hclwrite.File
		resourceAddress  string
		trgtDependencies *generictools.DepedendcyAddresses
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
