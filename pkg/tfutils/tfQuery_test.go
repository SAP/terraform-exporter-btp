package tfutils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemoveQueryGeneratedFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), queryGeneratedFileName)
	require.NoError(t, os.WriteFile(path, []byte("generated"), 0o600))

	require.NoError(t, removeQueryGeneratedFile(path))
	_, err := os.Stat(path)
	assert.ErrorIs(t, err, os.ErrNotExist)
	require.NoError(t, removeQueryGeneratedFile(path))
}

func TestCloudFoundryUserQueryAndExport(t *testing.T) {
	spec, err := GetListResourceSpec(CmdCfUserParameter, OrganizationLevel)
	require.NoError(t, err)
	assert.Equal(t, "cloudfoundry_user_cf", spec.ResourceType)

	scope := queryScope(spec.ResourceType, "", "", "org-guid", "")
	assert.Equal(t, map[string]string{"org": "org-guid"}, scope)
	query := renderListQuery(spec, scope)
	assert.Contains(t, query, `list "cloudfoundry_user_cf" "export"`)
	assert.Contains(t, query, `org = "org-guid"`)
	assert.Contains(t, query, "include_resource = true")

	config := `resource "cloudfoundry_user_cf" "export_0" {
  username = "user@example.com"
  origin = "ldap"
  annotations = { description = "cloudfoundry_user.literal" }
}`
	importConfig := `import {
  to = cloudfoundry_user_cf.export_0
  identity = { user_guid = "user-guid" }
}`
	input, err := json.Marshal(map[string]any{
		"type": "list_resource_found",
		"list_resource_found": map[string]any{
			"identity":        map[string]any{"user_guid": "user-guid"},
			"resource_object": map[string]any{"username": "user@example.com", "origin": "ldap"},
			"config":          config,
			"import_config":   importConfig,
		},
	})
	require.NoError(t, err)
	results, diagnostics, err := parseQueryOutput(input, spec)
	require.NoError(t, err)
	assert.Empty(t, diagnostics)
	require.Len(t, results, 1)
	assert.Equal(t, "user@example.com__ldap", results[0].InventoryValue)

	tempDir := t.TempDir()
	currentDir, err := os.Getwd()
	require.NoError(t, err)
	configDir, err := filepath.Rel(currentDir, tempDir)
	require.NoError(t, err)
	require.NoError(t, WriteQueryResults(configDir, "users.tf", spec.ResourceType, results))

	generated, err := os.ReadFile(filepath.Join(tempDir, "users.tf"))
	require.NoError(t, err)
	assert.Equal(t, config+"\n\n", string(generated))
	imports, err := os.ReadFile(filepath.Join(tempDir, "cloudfoundry_user_cf_import.tf"))
	require.NoError(t, err)
	assert.Equal(t, importConfig+"\n\n", string(imports))
}

func TestListResourceCatalogCoversEveryAllowedResource(t *testing.T) {
	testCases := []struct {
		level     string
		resources []string
	}{
		{SubaccountLevel, AllowedResourcesSubaccount},
		{DirectoryLevel, AllowedResourcesDirectory},
		{OrganizationLevel, AllowedResourcesOrganization},
	}

	for _, testCase := range testCases {
		for _, command := range testCase.resources {
			_, err := GetListResourceSpec(command, testCase.level)
			require.NoError(t, err, fmt.Sprintf("%s at %s", command, testCase.level))
		}
	}
}

func TestRenderListQuery(t *testing.T) {
	spec, err := GetListResourceSpec(CmdSubscriptionParameter, SubaccountLevel)
	require.NoError(t, err)

	query := renderListQuery(spec, map[string]string{"subaccount_id": "account-id"})

	assert.Contains(t, query, `list "btp_subaccount_subscription" "export"`)
	assert.Contains(t, query, "provider         = btp")
	assert.Contains(t, query, "include_resource = true")
	assert.Contains(t, query, `subaccount_id = "account-id"`)
}

func TestParseQueryOutput(t *testing.T) {
	input := strings.Join([]string{
		`{"type":"list_start"}`,
		`{"type":"list_resource_found","list_resource_found":{"identity":{"subaccount_id":"sa","app_name":"app","plan_name":"free"},"resource_object":{"app_name":"app"},"config":"resource \"btp_subaccount_subscription\" \"export_0\" {}","import_config":"import {}"}}`,
	}, "\n")
	spec, err := GetListResourceSpec(CmdSubscriptionParameter, SubaccountLevel)
	require.NoError(t, err)

	results, diagnostics, err := parseQueryOutput([]byte(input), spec)

	require.NoError(t, err)
	assert.Empty(t, diagnostics)
	require.Len(t, results, 1)
	assert.Equal(t, "app_free", results[0].InventoryValue)
	assert.Equal(t, "import {}", results[0].ImportConfig)
}

func TestParseQueryDiagnostic(t *testing.T) {
	input := `{"type":"diagnostic","diagnostic":{"severity":"error","summary":"query failed","detail":"denied"}}`
	spec, err := GetListResourceSpec(CmdSubaccountParameter, SubaccountLevel)
	require.NoError(t, err)

	_, diagnostics, err := parseQueryOutput([]byte(input), spec)

	require.NoError(t, err)
	assert.Contains(t, diagnostics, "query failed: denied")
}

func TestInventoryValuePreservesLegacyInventoryNames(t *testing.T) {
	testCases := []struct {
		resourceType string
		resource     map[string]any
		identity     map[string]any
		want         string
	}{
		{DirectoryRoleCollectionType, map[string]any{"name": "Directory Custom Role Collection"}, nil, "directory_custom_role_collection"},
		{SubaccountRoleType, map[string]any{"name": "Subaccount Custom Role"}, nil, "subaccount_custom_role"},
		{SubaccountEnvironmentInstanceType, map[string]any{"environment_type": "Cloud Foundry"}, nil, "cloud_foundry"},
		{SubaccountTrustConfigurationType, nil, map[string]any{"origin": "SAP ID Service"}, "SAP_ID_Service"},
		{CfUserType, map[string]any{"username": "Jane Doe"}, nil, "jane_doe"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.resourceType, func(t *testing.T) {
			got := inventoryValue(testCase.resourceType, QueryResult{Resource: testCase.resource, Identity: testCase.identity})
			assert.Equal(t, testCase.want, got)
		})
	}
}

func TestCfUserInventoryValueIncludesOrigin(t *testing.T) {
	first := inventoryValue(CfUserType, QueryResult{Resource: map[string]any{"username": "user@example.com", "origin": "ldap"}})
	second := inventoryValue(CfUserType, QueryResult{Resource: map[string]any{"username": "user@example.com", "origin": "sap.default"}})

	assert.Equal(t, "user@example.com__ldap", first)
	assert.Equal(t, "user@example.com__sap.default", second)
	assert.NotEqual(t, first, second)
}

func TestFilterTargetResource(t *testing.T) {
	spec, err := GetListResourceSpec(CmdSubaccountParameter, SubaccountLevel)
	require.NoError(t, err)
	results := []QueryResult{
		{Identity: map[string]any{"subaccount_id": "wanted"}, InventoryValue: "wanted"},
		{Identity: map[string]any{"subaccount_id": "other"}, InventoryValue: "other"},
	}

	filtered := filterTargetResource(spec, results, "wanted", "")

	require.Len(t, filtered, 1)
	assert.Equal(t, "wanted", filtered[0].InventoryValue)
}

func TestFilterSelectedResultsReportsMissingValue(t *testing.T) {
	spec, err := GetListResourceSpec(CmdCfSpaceParameter, OrganizationLevel)
	require.NoError(t, err)
	results := []QueryResult{{Resource: map[string]any{"name": "dev"}, InventoryValue: "dev"}}

	_, err = filterQueryResults(spec, results, "", "", []string{"prod"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "prod")
}

func TestExplicitRoleSelectionUsesInventoryIdentifier(t *testing.T) {
	spec, err := GetListResourceSpec(CmdRoleParameter, SubaccountLevel)
	require.NoError(t, err)
	results := []QueryResult{{
		Resource:       map[string]any{"name": "Subaccount Admin"},
		InventoryValue: "subaccount_admin",
	}}

	selected, err := filterQueryResults(spec, results, "subaccount-id", "", []string{"subaccount_admin"})

	require.NoError(t, err)
	require.Len(t, selected, 1)
	assert.Equal(t, "Subaccount Admin", selected[0].Resource["name"])
}
