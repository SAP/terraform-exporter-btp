package resourceprocessor

import (
	"fmt"
	"strings"
	"testing"

	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/SAP/terraform-exporter-btp/pkg/tfutils"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAssignmentPostProcessingRemovesUnusedTargets(t *testing.T) {
	for _, kind := range []string{"user", "group", "attribute"} {
		for _, includeCollection := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/collection=%t", kind, includeCollection), func(t *testing.T) {
				fields := map[string]string{"user_name": "null", "group_name": "null", "attribute_name": "null", "attribute_value": "null"}
				switch kind {
				case "user":
					fields["user_name"] = `"Alice@example.com"`
				case "group":
					fields["group_name"] = `"Auditors"`
				case "attribute":
					fields["attribute_name"], fields["attribute_value"] = `"cost_center"`, `"null"`
				}
				config := `resource "btp_subaccount_role_collection_assignment" "export_0" {
  subaccount_id = "sa"
  role_collection_name = "Subaccount Viewer"
  origin = "custom-idp"
`
				for name, value := range fields {
					config += fmt.Sprintf("  %s = %s\n", name, value)
				}
				config += "}\n"
				if includeCollection {
					// The collection follows the assignment, as can happen in a curated inventory.
					config += `resource "btp_subaccount_role_collection" "export_1" {
  subaccount_id = "sa"
  name = "Subaccount Viewer"
}`
				}
				file, diagnostics := hclwrite.ParseConfig([]byte(config), "test.tf", hcl.InitialPos)
				require.False(t, diagnostics.HasErrors(), diagnostics.Error())
				variables := generictools.VariableContent{}
				dependencies := generictools.NewDependencyAddresses()
				dependencies.SubaccountAddress = "btp_subaccount.export_2"
				ProcessResources(file, tfutils.SubaccountLevel, &variables, &dependencies, nil, generictools.LevelIds{SubaccountId: "sa"})
				body := file.Body().Blocks()[0].Body()
				for name, value := range fields {
					attribute := body.GetAttribute(name)
					if value == "null" {
						assert.Nil(t, attribute, "unused %s must be removed", name)
					} else {
						require.NotNil(t, attribute)
						assert.Equal(t, value, strings.TrimSpace(string(attribute.Expr().BuildTokens(nil).Bytes())))
					}
				}
				assert.Equal(t, `"custom-idp"`, strings.TrimSpace(string(body.GetAttribute("origin").Expr().BuildTokens(nil).Bytes())))
				assert.Equal(t, "btp_subaccount.export_2.id", strings.TrimSpace(string(body.GetAttribute("subaccount_id").Expr().BuildTokens(nil).Bytes())))
				wantCollection := `"Subaccount Viewer"`
				if includeCollection {
					wantCollection = "btp_subaccount_role_collection.export_1.name"
				}
				assert.Equal(t, wantCollection, strings.TrimSpace(string(body.GetAttribute("role_collection_name").Expr().BuildTokens(nil).Bytes())))
			})
		}
	}
}

func TestAssignmentCollectionReferenceMatchesScopeAndEscapedName(t *testing.T) {
	config := `resource "btp_subaccount_role_collection_assignment" "export_0" {
  subaccount_id = "sa"
  role_collection_name = "Team \"A\""
}
resource "btp_subaccount_role_collection" "correct" {
  subaccount_id = "sa"
  name = "Team \"A\""
}
resource "btp_subaccount_role_collection" "other_subaccount" {
  subaccount_id = "other-sa"
  name = "Team \"A\""
}
resource "btp_subaccount_role_collection_assignment" "unexported_collection" {
  subaccount_id = "sa"
  role_collection_name = "Subaccount Administrator"
}
resource "btp_subaccount_role_collection_assignment" "existing_reference" {
  subaccount_id = var.subaccount_id
  role_collection_name = btp_subaccount_role_collection.correct.name
}`
	file, diagnostics := hclwrite.ParseConfig([]byte(config), "test.tf", hcl.InitialPos)
	require.False(t, diagnostics.HasErrors(), diagnostics.Error())
	linkRoleCollectionAssignments(file.Body())
	blocks := file.Body().Blocks()
	assert.Equal(t, "btp_subaccount_role_collection.correct.name", strings.TrimSpace(string(blocks[0].Body().GetAttribute("role_collection_name").Expr().BuildTokens(nil).Bytes())))
	assert.Equal(t, `"Subaccount Administrator"`, strings.TrimSpace(string(blocks[3].Body().GetAttribute("role_collection_name").Expr().BuildTokens(nil).Bytes())))
	assert.Equal(t, "btp_subaccount_role_collection.correct.name", strings.TrimSpace(string(blocks[4].Body().GetAttribute("role_collection_name").Expr().BuildTokens(nil).Bytes())))
	before := string(file.Bytes())
	linkRoleCollectionAssignments(file.Body())
	assert.Equal(t, before, string(file.Bytes()))
}
