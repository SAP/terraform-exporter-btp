package resourceprocessor

import (
	"testing"

	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/testutils"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/assert"
)

func TestFillRoleDependencyAddresses(t *testing.T) {

	srcFileRole, _ := testutils.GetHclFilesById("sa_role")
	srcFileRoleIncomplete, _ := testutils.GetHclFilesById("sa_role_error")

	emptyTestDependencies := generictools.NewDependencyAddresses()
	defaultTestDependencies := generictools.NewDependencyAddresses()

	defaultTestDependencies.RoleAddress = make(map[generictools.RoleKey]string)
	defaultTestDependencies.RoleAddress[generictools.RoleKey{
		AppId:            "destination-xsappname!b62",
		Name:             "Destination Administrator Instance",
		RoleTemplateName: "Destination_Administrator_Instance",
	}] = "btp_subaccount_role.role_6"

	tests := []struct {
		name             string
		src              *hclwrite.File
		resourceAddress  string
		trgtDependencies *generictools.DependencyAddresses
	}{
		{
			name:             "Test Role Dependency Address",
			src:              srcFileRole,
			resourceAddress:  "btp_subaccount_role.role_6",
			trgtDependencies: &defaultTestDependencies,
		},
		{
			name:             "Test Role Incomplete",
			src:              srcFileRoleIncomplete,
			resourceAddress:  "btp_subaccount_role.role_6",
			trgtDependencies: &emptyTestDependencies,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			dependencies := generictools.NewDependencyAddresses()
			blocks := tt.src.Body().Blocks()
			// we assume one resource entry in the blocks file
			fillRoleDependencyAddresses(blocks[0].Body(), tt.resourceAddress, &dependencies)

			assert.Equal(t, tt.trgtDependencies.RoleAddress, dependencies.RoleAddress)
		})
	}
}
