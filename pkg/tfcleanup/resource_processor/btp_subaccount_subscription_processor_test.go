package resourceprocessor

import (
	"testing"

	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/testutils"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/assert"
)

func TestAddEntitlementDependency(t *testing.T) {

	srcFileSubscriptionDep, trgtFileSubscriptionDep := testutils.GetHclFilesById("sa_subscription_dependency")
	srcFileSubscriptionNoDep, trgtFileSubscriptionNoDep := testutils.GetHclFilesById("sa_subscription_no_dependency")

	emptyTestDependencies := generictools.NewDependencyAddresses()
	defaultTestDependencies := generictools.NewDependencyAddresses()

	entitlementKey := generictools.EntitlementKey{
		ServiceName: "feature-flags-dashboard",
		PlanName:    "dashboard",
	}

	defaultTestDependencies.EntitlementAddress = make(map[generictools.EntitlementKey]string)
	defaultTestDependencies.EntitlementAddress[entitlementKey] = "btp_subaccount_entitlement.entitlement_2"

	tests := []struct {
		name         string
		src          *hclwrite.File
		trgt         *hclwrite.File
		dependencies *generictools.DepedendcyAddresses
	}{
		{
			name:         "Test Subscription Dependency",
			src:          srcFileSubscriptionDep,
			trgt:         trgtFileSubscriptionDep,
			dependencies: &defaultTestDependencies,
		},
		{
			name:         "Test No Subscription Dependency",
			src:          srcFileSubscriptionNoDep,
			trgt:         trgtFileSubscriptionNoDep,
			dependencies: &emptyTestDependencies,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			blocks := tt.src.Body().Blocks()
			// we assume one resource entry in the blocks file
			addEntitlementDependency(blocks[0].Body(), tt.dependencies)
			assert.NoError(t, testutils.AreHclFilesEqual(tt.src, tt.trgt))
		})
	}
}
