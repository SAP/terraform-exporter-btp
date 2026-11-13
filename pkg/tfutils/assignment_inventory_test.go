package tfutils

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoleCollectionAssignmentCatalog(t *testing.T) {
	const command = "role-collection-assignments"
	assert.Contains(t, GetResourcesList("all", SubaccountLevel), command)
	assert.Equal(t, []string{command}, GetResourcesList(command, SubaccountLevel))
	assert.NotContains(t, GetValidResourcesByLevel(DirectoryLevel), command)
	assert.NotContains(t, GetValidResourcesByLevel(OrganizationLevel), command)
	spec, err := GetListResourceSpec(command, SubaccountLevel)
	require.NoError(t, err)
	assert.Equal(t, "btp_subaccount_role_collection_assignment", spec.ResourceType)
	query := renderListQuery(spec, queryScope(spec.ResourceType, "sa", "", "", ""))
	assert.Contains(t, query, `list "btp_subaccount_role_collection_assignment" "export"`)
	assert.Contains(t, query, `subaccount_id = "sa"`)
	assert.Contains(t, query, "include_resource = true")
	assert.NotContains(t, query, "role_collection_name")
	assert.NotContains(t, query, "origin")
}

func TestRoleCollectionAssignmentInventory(t *testing.T) {
	tests := []struct {
		name   string
		target map[string]any
		want   string
	}{
		{"user", map[string]any{"user_name": "Alice@example.com"}, "Subaccount%20Viewer|ldap|user|Alice@example.com"},
		{"group", map[string]any{"group_name": "Auditors"}, "Subaccount%20Viewer|ldap|group|Auditors"},
		{"attribute", map[string]any{"attribute_name": "cost_center", "attribute_value": "1234"}, "Subaccount%20Viewer|ldap|attribute|cost_center|1234"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			identity := assignmentTestIdentity()
			maps.Copy(identity, tt.target)
			assert.Equal(t, tt.want, inventoryValue(SubaccountRoleCollectionAssignmentType, QueryResult{Identity: identity}))
		})
	}
}

func TestRoleCollectionAssignmentInventoryDoesNotMergeIdentities(t *testing.T) {
	seen := make(map[string]bool)
	for _, collection := range []string{"A B", "a_b", "A|B", "A%7CB"} {
		for _, origin := range []string{"ldap", "LDAP", "custom|idp", ""} {
			for _, target := range []string{"user_name", "group_name", "attribute_name"} {
				identity := assignmentTestIdentity()
				identity["role_collection_name"], identity["origin"] = collection, origin
				identity[target] = "A|B/%_ ü"
				if target == "attribute_name" {
					identity["attribute_value"] = "1|2"
				}
				value := roleCollectionAssignmentInventoryValue(identity)
				require.NotEmpty(t, value)
				assert.False(t, seen[value], "duplicate inventory identifier: %s", value)
				seen[value] = true
				parts := strings.Split(value, "|")
				for i, expected := range map[int]string{0: collection, 1: origin, 3: "A|B/%_ ü"} {
					decoded, err := url.PathUnescape(parts[i])
					require.NoError(t, err)
					assert.Equal(t, expected, decoded)
				}
			}
		}
	}
}

func TestRoleCollectionAssignmentRejectsIncompleteIdentity(t *testing.T) {
	tests := []map[string]any{
		{},
		{"user_name": "alice", "origin": nil},
		{"user_name": "alice", "subaccount_id": ""},
		{"user_name": "alice", "role_collection_name": nil},
		{"user_name": 123},
		{"user_name": "alice", "group_name": "team"},
		{"attribute_name": "cost_center"},
		{"attribute_value": "1234"},
		{"user_name": "alice", "attribute_value": "1234"},
	}
	spec, err := GetListResourceSpec(CmdRoleCollectionAssignmentParameter, SubaccountLevel)
	require.NoError(t, err)
	for _, overrides := range tests {
		identity := assignmentTestIdentity()
		maps.Copy(identity, overrides)
		_, _, err := parseQueryOutput(assignmentTestEvent(t, identity, "", ""), spec)
		require.ErrorContains(t, err, "without an inventory identity")
	}
}

func TestRoleCollectionAssignmentQuerySelection(t *testing.T) {
	spec, err := GetListResourceSpec(CmdRoleCollectionAssignmentParameter, SubaccountLevel)
	require.NoError(t, err)
	var events []string
	for _, origin := range []string{"ldap", "custom-idp"} {
		identity := assignmentTestIdentity()
		identity["user_name"], identity["origin"] = "Alice@example.com", origin
		events = append(events, string(assignmentTestEvent(t, identity, "", "")))
	}
	results, diagnostics, err := parseQueryOutput([]byte(strings.Join(events, "\n")), spec)
	require.NoError(t, err)
	assert.Empty(t, diagnostics)
	require.Len(t, results, 2)
	all, err := filterQueryResults(spec, results, "sa", "", nil)
	require.NoError(t, err)
	assert.Len(t, all, 2, "assignments to SAP-provided collections must be included")
	selected, err := filterQueryResults(spec, results, "sa", "", []string{results[1].InventoryValue})
	require.NoError(t, err)
	require.Len(t, selected, 1)
	assert.Equal(t, "custom-idp", selected[0].Identity["origin"])
	_, err = filterQueryResults(spec, results, "sa", "", []string{"missing"})
	require.ErrorContains(t, err, "was not found")
}

func TestRoleCollectionAssignmentExport(t *testing.T) {
	spec, err := GetListResourceSpec(CmdRoleCollectionAssignmentParameter, SubaccountLevel)
	require.NoError(t, err)
	for _, kind := range []string{"user", "group", "attribute"} {
		t.Run(kind, func(t *testing.T) {
			identity := assignmentTestIdentity()
			if kind == "attribute" {
				identity["attribute_name"], identity["attribute_value"] = "cost_center", "1234"
			} else {
				identity[kind+"_name"] = "Alice@example.com"
			}
			var targets, importTargets strings.Builder
			for _, field := range []string{"user_name", "group_name", "attribute_name", "attribute_value"} {
				if value := identity[field]; value != nil {
					fmt.Fprintf(&targets, "  %s = %q\n", field, value)
					fmt.Fprintf(&importTargets, "    %s = %q\n", field, value)
				} else {
					fmt.Fprintf(&targets, "  %s = null\n", field)
				}
			}
			config := `resource "btp_subaccount_role_collection_assignment" "export_0" {
  subaccount_id = "sa"
  role_collection_name = "Subaccount Viewer"
  origin = "ldap"
` + targets.String() + "}"
			importConfig := `import {
  to = btp_subaccount_role_collection_assignment.export_0
  identity = {
    subaccount_id = "sa"
    role_collection_name = "Subaccount Viewer"
    origin = "ldap"
` + importTargets.String() + "  }\n}"
			results, _, err := parseQueryOutput(assignmentTestEvent(t, identity, config, importConfig), spec)
			require.NoError(t, err)
			directory := t.TempDir()
			currentDir, err := os.Getwd()
			require.NoError(t, err)
			relativeDir, err := filepath.Rel(currentDir, directory)
			require.NoError(t, err)
			require.NoError(t, WriteQueryResults(relativeDir, "assignments.tf", spec.ResourceType, results))
			generated, err := os.ReadFile(filepath.Join(directory, "assignments.tf"))
			require.NoError(t, err)
			assert.Equal(t, config+"\n\n", string(generated))
			importsPath := filepath.Join(directory, spec.ResourceType+"_import.tf")
			imports, err := os.ReadFile(importsPath)
			require.NoError(t, err)
			assert.Equal(t, importConfig+"\n\n", string(imports))
			require.NoError(t, applyExecutionTargetVariables(directory, ExecutionTarget{SubaccountID: "sa"}))
			require.NoError(t, normalizeGeneratedResourceAddresses(directory))
			imports, err = os.ReadFile(importsPath)
			require.NoError(t, err)
			assert.Contains(t, string(imports), "btp_subaccount_role_collection_assignment.subaccount_role_collection_assignment_0")
			assert.Contains(t, string(imports), "var.subaccount_id")
			for field, value := range identity {
				if value != nil && field != "subaccount_id" {
					assert.Contains(t, string(imports), fmt.Sprintf("%q", value))
				}
			}
		})
	}
}

func TestRoleCollectionAssignmentEmptyQueryAndDiagnostics(t *testing.T) {
	spec, err := GetListResourceSpec(CmdRoleCollectionAssignmentParameter, SubaccountLevel)
	require.NoError(t, err)
	results, diagnostics, err := parseQueryOutput([]byte(`{"type":"list_start"}`), spec)
	require.NoError(t, err)
	assert.Empty(t, results)
	assert.Empty(t, diagnostics)
	_, diagnostics, err = parseQueryOutput([]byte(`{"type":"diagnostic","diagnostic":{"severity":"error","summary":"Unable to Resolve Role Collection Assignment Origin","detail":"multiple trust configurations match"}}`), spec)
	require.NoError(t, err)
	assert.Contains(t, diagnostics, "multiple trust configurations match")
}

func assignmentTestIdentity() map[string]any {
	return map[string]any{
		"subaccount_id": "sa", "role_collection_name": "Subaccount Viewer", "origin": "ldap",
		"user_name": nil, "group_name": nil, "attribute_name": nil, "attribute_value": nil,
	}
}

// Synthetic query events follow the provider 1.28.0 identity schema; they are
// deliberately independent of a live subaccount or its assignments.
func assignmentTestEvent(t *testing.T, identity map[string]any, config, importConfig string) []byte {
	t.Helper()
	data, err := json.Marshal(map[string]any{
		"type": "list_resource_found",
		"list_resource_found": map[string]any{
			"identity": identity, "resource_object": identity,
			"config": config, "import_config": importConfig,
		},
	})
	require.NoError(t, err)
	return data
}
