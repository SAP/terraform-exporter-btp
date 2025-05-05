package resourceprocessor

import (
	"testing"

	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/testutils"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/assert"
)

func TestAddServiceInstanceDependency(t *testing.T) {
	srcFileServiceInstanceDep, trgtFileServiceInstanceDep := testutils.GetHclFilesById("sa_service_instance_dependency")
	srcFileServiceInstanceNoDep, trgtFileServiceInstanceNoDep := testutils.GetHclFilesById("sa_service_instance_no_dependency")

	defaultTestDependencies := getNewServiceInstanceDepTemplate()
	defaultTestDependenciesCopy := getNewServiceInstanceDepTemplate()

	targetDependencies := getNewServiceInstanceDepTemplate()
	targetDependencies.DataSourceInfo = append(targetDependencies.DataSourceInfo, generictools.DataSourceInfo{
		DatasourceAddress:  "alert-notification_standard",
		SubaccountAddress:  "btp_subaccount.subaccount_0.id",
		OfferingName:       "alert-notification",
		Name:               "standard",
		EntitlementAddress: "btp_subaccount_entitlement.entitlement_0",
	})

	tests := []struct {
		name             string
		src              *hclwrite.File
		trgt             *hclwrite.File
		dependencies     generictools.DependencyAddresses
		trgtDependencies generictools.DependencyAddresses
	}{
		{
			name:             "Test Service Instance Dependency",
			src:              srcFileServiceInstanceDep,
			trgt:             trgtFileServiceInstanceDep,
			dependencies:     defaultTestDependencies,
			trgtDependencies: targetDependencies,
		},
		{
			name:             "Test No Service Instance Dependency",
			src:              srcFileServiceInstanceNoDep,
			trgt:             trgtFileServiceInstanceNoDep,
			dependencies:     defaultTestDependencies,
			trgtDependencies: defaultTestDependenciesCopy,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			blocks := tt.src.Body().Blocks()
			// we assume one resource entry in the blocks file
			addServiceInstanceDependency(blocks[0].Body(), &tt.dependencies, nil, "subaccount_id")
			assert.NoError(t, testutils.AreHclFilesEqual(tt.src, tt.trgt))
			assert.Equal(t, &tt.dependencies, &tt.trgtDependencies)
		})
	}

}

func getNewServiceInstanceDepTemplate() generictools.DependencyAddresses {
	defaultTestDependencies := generictools.NewDependencyAddresses()

	entitlementKey := generictools.EntitlementKey{
		ServiceName: "alert-notification",
		PlanName:    "standard",
	}

	defaultTestDependencies.EntitlementAddress = make(map[generictools.EntitlementKey]string)
	defaultTestDependencies.EntitlementAddress[entitlementKey] = "btp_subaccount_entitlement.entitlement_0"
	defaultTestDependencies.SubaccountAddress = "btp_subaccount.subaccount_0"

	return defaultTestDependencies
}
