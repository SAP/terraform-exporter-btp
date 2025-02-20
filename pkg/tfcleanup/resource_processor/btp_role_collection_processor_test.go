package resourceprocessor

import (
	"testing"

	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/testutils"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/assert"
)

func TestAddRoleDependency(t *testing.T) {
	srcFileRoleCollectionDep, trgtFileRoleCollectionDep := testutils.GetHclFilesById("sa_role_collection_dependency")

	defaultTestDependencies := getNewRoleDepTemplate()

	tests := []struct {
		name         string
		src          *hclwrite.File
		trgt         *hclwrite.File
		dependencies generictools.DepedendcyAddresses
	}{
		{
			name:         "Test Role Collection Dependency",
			src:          srcFileRoleCollectionDep,
			trgt:         trgtFileRoleCollectionDep,
			dependencies: defaultTestDependencies,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blocks := tt.src.Body().Blocks()
			// we assume one resource entry in the blocks file
			addRoleDependency(blocks[0].Body(), nil)
			assert.NoError(t, testutils.AreHclFilesEqual(tt.src, tt.trgt))
		})
	}
}

func TestBuildDependencyString(t *testing.T) {
	tests := []struct {
		name                string
		roles               []Role
		dependencyAddresses *generictools.DepedendcyAddresses
		expected            string
	}{
		{
			name: "single role with dependency",
			roles: []Role{
				{
					Name:              "role1",
					RoleTemplateAppID: "app1",
					RoleTemplateName:  "template1",
				},
			},
			dependencyAddresses: &generictools.DepedendcyAddresses{
				RoleAddress: map[generictools.RoleKey]string{
					{
						AppId:            "app1",
						Name:             "role1",
						RoleTemplateName: "template1",
					}: "dependency1",
				},
			},
			expected: "dependency1",
		},
		{
			name: "multiple roles with dependencies",
			roles: []Role{
				{
					Name:              "role1",
					RoleTemplateAppID: "app1",
					RoleTemplateName:  "template1",
				},
				{
					Name:              "role2",
					RoleTemplateAppID: "app2",
					RoleTemplateName:  "template2",
				},
			},
			dependencyAddresses: &generictools.DepedendcyAddresses{
				RoleAddress: map[generictools.RoleKey]string{
					{
						AppId:            "app1",
						Name:             "role1",
						RoleTemplateName: "template1",
					}: "dependency1",
					{
						AppId:            "app2",
						Name:             "role2",
						RoleTemplateName: "template2",
					}: "dependency2",
				},
			},
			expected: "dependency1, dependency2",
		},
		{
			name: "role with no dependency",
			roles: []Role{
				{
					Name:              "role1",
					RoleTemplateAppID: "app1",
					RoleTemplateName:  "template1",
				},
			},
			dependencyAddresses: &generictools.DepedendcyAddresses{
				RoleAddress: map[generictools.RoleKey]string{},
			},
			expected: "",
		},
		{
			name: "mixed roles with and without dependencies",
			roles: []Role{
				{
					Name:              "role1",
					RoleTemplateAppID: "app1",
					RoleTemplateName:  "template1",
				},
				{
					Name:              "role2",
					RoleTemplateAppID: "app2",
					RoleTemplateName:  "template2",
				},
			},
			dependencyAddresses: &generictools.DepedendcyAddresses{
				RoleAddress: map[generictools.RoleKey]string{
					{
						AppId:            "app1",
						Name:             "role1",
						RoleTemplateName: "template1",
					}: "dependency1",
				},
			},
			expected: "dependency1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildDependencyString(tt.roles, tt.dependencyAddresses)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func getNewRoleDepTemplate() generictools.DepedendcyAddresses {
	defaultTestDependencies := generictools.NewDependencyAddresses()

	roleKey := generictools.RoleKey{
		AppId:            "app1",
		Name:             "role1",
		RoleTemplateName: "template1",
	}

	defaultTestDependencies.RoleAddress = make(map[generictools.RoleKey]string)
	defaultTestDependencies.RoleAddress[roleKey] = "dependency1"

	return defaultTestDependencies
}
