package tfutils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/SAP/terraform-exporter-btp/pkg/defaultfilter"
	"github.com/SAP/terraform-exporter-btp/pkg/files"
	"github.com/SAP/terraform-exporter-btp/pkg/output"
)

const (
	queryResultLimit       = 10000
	queryGeneratedFileName = "query-generated.tf"
)

type ListResourceSpec struct {
	CommandName  string
	ResourceType string
	Provider     string
	Level        string
}

type QueryResult struct {
	Identity       map[string]any
	Resource       map[string]any
	Config         string
	ImportConfig   string
	InventoryValue string
}

type queryEvent struct {
	Type              string              `json:"type"`
	Diagnostic        *queryDiagnostic    `json:"diagnostic"`
	ListResourceFound *queryResourceFound `json:"list_resource_found"`
}

type queryDiagnostic struct {
	Severity string `json:"severity"`
	Summary  string `json:"summary"`
	Detail   string `json:"detail"`
}

type queryResourceFound struct {
	Identity       map[string]any `json:"identity"`
	ResourceObject map[string]any `json:"resource_object"`
	Config         string         `json:"config"`
	ImportConfig   string         `json:"import_config"`
}

var listResourceSpecs = []ListResourceSpec{
	{CmdDirectoryParameter, DirectoryType, "btp", DirectoryLevel},
	{CmdEntitlementParameter, DirectoryEntitlementType, "btp", DirectoryLevel},
	{CmdRoleParameter, DirectoryRoleType, "btp", DirectoryLevel},
	{CmdRoleCollectionParameter, DirectoryRoleCollectionType, "btp", DirectoryLevel},
	{CmdSubaccountParameter, SubaccountType, "btp", SubaccountLevel},
	{CmdEntitlementParameter, SubaccountEntitlementType, "btp", SubaccountLevel},
	{CmdEnvironmentInstanceParameter, SubaccountEnvironmentInstanceType, "btp", SubaccountLevel},
	{CmdSubscriptionParameter, SubaccountSubscriptionType, "btp", SubaccountLevel},
	{CmdTrustConfigurationParameter, SubaccountTrustConfigurationType, "btp", SubaccountLevel},
	{CmdRoleParameter, SubaccountRoleType, "btp", SubaccountLevel},
	{CmdRoleCollectionParameter, SubaccountRoleCollectionType, "btp", SubaccountLevel},
	{CmdServiceBindingParameter, SubaccountServiceBindingType, "btp", SubaccountLevel},
	{CmdServiceInstanceParameter, SubaccountServiceInstanceType, "btp", SubaccountLevel},
	{CmdSecuritySettingParameter, SubaccountSecuritySettingType, "btp", SubaccountLevel},
	{CmdCfSpaceParameter, CfSpaceType, "cloudfoundry", OrganizationLevel},
	{CmdCfUserParameter, CfUserType, "cloudfoundry", OrganizationLevel},
	{CmdCfDomainParamater, CfDomainType, "cloudfoundry", OrganizationLevel},
	{CmdCfOrgRoleParameter, CfOrgRoleType, "cloudfoundry", OrganizationLevel},
	{CmdCfRouteParameter, CfRouteType, "cloudfoundry", OrganizationLevel},
	{CmdCfSpaceQuotaParameter, CfSpaceQuotaType, "cloudfoundry", OrganizationLevel},
	{CmdCfServiceInstanceParameter, CfServiceInstanceType, "cloudfoundry", OrganizationLevel},
	{CmdCfSpaceRoleParameter, CfSpaceRoleType, "cloudfoundry", SpaceLevel},
}

func GetListResourceSpec(commandName, level string) (ListResourceSpec, error) {
	for _, spec := range listResourceSpecs {
		if spec.CommandName == commandName && (spec.Level == level || (spec.Level == SpaceLevel && level == OrganizationLevel)) {
			return spec, nil
		}
	}
	return ListResourceSpec{}, fmt.Errorf("resource %q is not supported at %s", commandName, level)
}

func ReadListResources(subaccountID, directoryID, organizationID string, resourceList []string) (BtpResources, error) {
	level, _ := GetExecutionLevelAndId(subaccountID, directoryID, organizationID, "")
	if err := ValidateTerraformQuerySupport(); err != nil {
		return BtpResources{}, err
	}
	if err := PrepareQueryDirectory(TmpFolder); err != nil {
		return BtpResources{}, err
	}

	resources := make([]BtpResource, 0, len(resourceList))
	var directoryFeatures []string
	for _, commandName := range resourceList {
		if level == DirectoryLevel && !directoryResourceEnabled(commandName, directoryFeatures) {
			continue
		}
		_, results, err := QueryResource(TmpFolder, commandName, subaccountID, directoryID, organizationID, "", nil, false)
		if err != nil {
			return BtpResources{}, err
		}
		if commandName == CmdDirectoryParameter && len(results) == 1 {
			directoryFeatures = stringSlice(results[0].Resource["features"])
		}
		values := make([]string, 0, len(results))
		for _, result := range results {
			values = append(values, result.InventoryValue)
		}
		sort.Strings(values)
		if len(values) > 0 {
			resources = append(resources, BtpResource{Name: commandName, Values: values})
		}
	}
	return BtpResources{BtpResources: resources}, nil
}

func QueryResource(configDir, commandName, subaccountID, directoryID, organizationID, spaceID string, selected []string, generate bool) (string, []QueryResult, error) {
	level, _ := GetExecutionLevelAndId(subaccountID, directoryID, organizationID, spaceID)
	lookupLevel := level
	if commandName == CmdCfSpaceRoleParameter {
		lookupLevel = OrganizationLevel
	}
	spec, err := GetListResourceSpec(commandName, lookupLevel)
	if err != nil {
		return "", nil, err
	}

	if commandName == CmdCfSpaceRoleParameter && spaceID == "" {
		_, spaces, err := QueryResource(configDir, CmdCfSpaceParameter, "", "", organizationID, "", nil, false)
		if err != nil {
			return spec.ResourceType, nil, err
		}
		var all []QueryResult
		for _, space := range spaces {
			id := firstString(space.Identity, "space_guid", "id")
			if id == "" {
				id = firstString(space.Resource, "id")
			}
			if id == "" {
				return spec.ResourceType, nil, fmt.Errorf("space query returned a result without an identity")
			}
			results, err := runListQuery(configDir, spec, map[string]string{"space": id}, generate)
			if err != nil {
				return spec.ResourceType, nil, err
			}
			all = append(all, results...)
		}
		results, err := filterQueryResults(spec, all, subaccountID, directoryID, selected)
		return spec.ResourceType, results, err
	}

	scope := queryScope(spec.ResourceType, subaccountID, directoryID, organizationID, spaceID)
	results, err := runListQuery(configDir, spec, scope, generate)
	if err != nil {
		return spec.ResourceType, nil, err
	}
	results = filterTargetResource(spec, results, subaccountID, directoryID)
	results, err = filterQueryResults(spec, results, subaccountID, directoryID, selected)
	return spec.ResourceType, results, err
}

func WriteQueryResults(configDir, resourceFileName, resourceType string, results []QueryResult) error {
	var configs, imports strings.Builder
	for _, result := range results {
		if result.Config == "" || result.ImportConfig == "" {
			return fmt.Errorf("terraform query did not generate configuration for %s", resourceType)
		}
		config := result.Config
		importConfig := result.ImportConfig
		if resourceType == CfUserExportType {
			config = replaceGeneratedResourceType(config, CfUserType, CfUserExportType)
			importConfig = replaceGeneratedResourceType(importConfig, CfUserType, CfUserExportType)
		}
		configs.WriteString(config)
		configs.WriteString("\n\n")
		imports.WriteString(importConfig)
		imports.WriteString("\n\n")
	}
	if err := files.CreateFileWithContent(filepath.Join(configDir, resourceFileName), configs.String()); err != nil {
		return err
	}
	return files.WriteImportConfiguration(configDir, resourceType, imports.String())
}

func ExportResourceType(queryResourceType string) string {
	if queryResourceType == CfUserType {
		return CfUserExportType
	}
	return queryResourceType
}

func replaceGeneratedResourceType(content, sourceType, targetType string) string {
	content = strings.ReplaceAll(content, `resource "`+sourceType+`"`, `resource "`+targetType+`"`)
	return strings.ReplaceAll(content, sourceType+".", targetType+".")
}

func PrepareQueryDirectory(directory string) error {
	addUserAgent()
	defer removeUserAgent()
	chdir := fmt.Sprintf("-chdir=%s", directory)
	if stderr, err := runTfCmdGeneric(chdir, "init", "-backend=false"); err != nil {
		return fmt.Errorf("error initializing Terraform query directory: %v: %s", err, stderr)
	}
	return nil
}

func runListQuery(directory string, spec ListResourceSpec, scope map[string]string, generate bool) ([]QueryResult, error) {
	queryPath := filepath.Join(directory, "resources.tfquery.hcl")
	if err := files.CreateFileWithContent(queryPath, renderListQuery(spec, scope)); err != nil {
		return nil, err
	}

	args := []string{fmt.Sprintf("-chdir=%s", directory), "query", "-json", "-no-color"}
	if generate {
		generatedPath := filepath.Join(directory, queryGeneratedFileName)
		if err := removeQueryGeneratedFile(generatedPath); err != nil {
			return nil, err
		}
		defer func() { _ = removeQueryGeneratedFile(generatedPath) }()
		args = append(args, "-generate-config-out="+queryGeneratedFileName)
	}
	cmd := exec.Command("terraform", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	addUserAgent()
	err := cmd.Run()
	removeUserAgent()
	results, diagnostics, parseErr := parseQueryOutput(stdout.Bytes(), spec)
	if parseErr != nil {
		return nil, parseErr
	}
	if err != nil {
		return nil, fmt.Errorf("terraform query failed: %v: %s%s", err, stderr.String(), diagnostics)
	}
	if diagnostics != "" {
		return nil, fmt.Errorf("terraform query failed: %s", diagnostics)
	}
	return results, nil
}

func removeQueryGeneratedFile(path string) error {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("cannot remove Terraform query scratch file %q: %w", path, err)
	}
	return nil
}

func renderListQuery(spec ListResourceSpec, scope map[string]string) string {
	var b strings.Builder
	label := "export"
	if spec.ResourceType == CfSpaceRoleType {
		label = SanitizeHclLabel("space_" + scope["space"])
	}
	fmt.Fprintf(&b, "list %q %q {\n  provider         = %s\n  include_resource = true\n  limit            = %d\n", spec.ResourceType, label, spec.Provider, queryResultLimit)
	if len(scope) > 0 {
		b.WriteString("  config {\n")
		keys := make([]string, 0, len(scope))
		for key := range scope {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			fmt.Fprintf(&b, "    %s = %q\n", key, scope[key])
		}
		b.WriteString("  }\n")
	}
	b.WriteString("}\n")
	return b.String()
}

func parseQueryOutput(data []byte, spec ListResourceSpec) ([]QueryResult, string, error) {
	var results []QueryResult
	var diagnostics strings.Builder
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Buffer(make([]byte, 64*1024), 16*1024*1024)
	for scanner.Scan() {
		var event queryEvent
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			return nil, "", fmt.Errorf("invalid Terraform query JSON: %w", err)
		}
		if event.Diagnostic != nil && event.Diagnostic.Severity == "error" {
			fmt.Fprintf(&diagnostics, "%s: %s\n", event.Diagnostic.Summary, event.Diagnostic.Detail)
		}
		if found := event.ListResourceFound; found != nil {
			result := QueryResult{Identity: found.Identity, Resource: found.ResourceObject, Config: found.Config, ImportConfig: found.ImportConfig}
			result.InventoryValue = inventoryValue(spec.ResourceType, result)
			if result.InventoryValue == "" {
				return nil, "", fmt.Errorf("%s query returned a result without an inventory identity", spec.ResourceType)
			}
			results = append(results, result)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, "", err
	}
	return results, diagnostics.String(), nil
}

// queryScope contains only the mandatory account or space arguments required to
// retrieve resources at the requested level. Optional provider-side filters are
// deliberately not used; exclusions happen after the query returns.
func queryScope(resourceType, subaccountID, directoryID, organizationID, spaceID string) map[string]string {
	switch resourceType {
	case DirectoryEntitlementType, DirectoryRoleType, DirectoryRoleCollectionType:
		return map[string]string{"directory_id": directoryID}
	case SubaccountEntitlementType, SubaccountEnvironmentInstanceType, SubaccountSubscriptionType, SubaccountTrustConfigurationType, SubaccountRoleType, SubaccountRoleCollectionType, SubaccountServiceInstanceType, SubaccountServiceBindingType, SubaccountSecuritySettingType:
		return map[string]string{"subaccount_id": subaccountID}
	case CfSpaceType, CfOrgRoleType, CfDomainType, CfRouteType, CfSpaceQuotaType, CfServiceInstanceType, CfUserType:
		return map[string]string{"org": organizationID}
	case CfSpaceRoleType:
		return map[string]string{"space": spaceID}
	default:
		return nil
	}
}

func inventoryValue(resourceType string, result QueryResult) string {
	r, i := result.Resource, result.Identity
	switch resourceType {
	case DirectoryType, SubaccountType:
		return firstString(r, "name", "username")
	case DirectoryRoleType, DirectoryRoleCollectionType, SubaccountRoleType, SubaccountRoleCollectionType, SubaccountServiceBindingType, CfSpaceType, CfDomainType, CfRouteType, CfSpaceQuotaType:
		return output.FormatResourceNameGeneric(firstString(r, "name", "username", "url"))
	case CfUserType:
		return cfUserInventoryValue(r, i)
	case DirectoryEntitlementType, SubaccountEntitlementType:
		return firstString(i, "service_name") + "_" + firstString(i, "plan_name")
	case SubaccountEnvironmentInstanceType:
		return output.FormatResourceNameGeneric(firstString(r, "environment_type", "name"))
	case SubaccountSubscriptionType:
		return firstString(i, "app_name") + "_" + firstString(i, "plan_name")
	case SubaccountTrustConfigurationType:
		return output.FormatResourceNameGenericNoCapConversion(firstString(i, "origin"))
	case SubaccountSecuritySettingType:
		return firstString(i, "subaccount_id")
	case SubaccountServiceInstanceType:
		return joinNonEmpty("_", firstString(r, "name"), firstString(r, "serviceplan_id", "service_plan_id"))
	case CfOrgRoleType:
		return output.FormatOrgRoleResourceName(firstString(r, "type"), firstString(r, "user", "username"))
	case CfSpaceRoleType:
		return output.FormatSpaceRoleResourceName(firstString(r, "type"), firstString(r, "space"), firstString(r, "user"))
	case CfServiceInstanceType:
		return joinNonEmpty("_", firstString(r, "name"), firstString(r, "service_plan", "service_plan_id"))
	default:
		return firstMapString(i)
	}
}

// cfUserInventoryValue distinguishes identities that share a username but
// belong to different identity-provider origins. The origin is part of the
// cloudfoundry_user resource schema and is returned as resource data by the
// list query. It may also be included in the identity, so use that as a
// fallback for provider versions that omit it from the resource object.
func cfUserInventoryValue(resource, identity map[string]any) string {
	username := output.FormatResourceNameGeneric(firstString(resource, "username", "name"))
	origin := firstString(resource, "origin")
	if origin == "" {
		origin = firstString(identity, "origin")
	}
	if origin == "" {
		return username
	}
	return username + "__" + output.FormatResourceNameGeneric(origin)
}

func filterTargetResource(spec ListResourceSpec, results []QueryResult, subaccountID, directoryID string) []QueryResult {
	var key, target string
	switch spec.ResourceType {
	case SubaccountType:
		key, target = "subaccount_id", subaccountID
	case DirectoryType:
		key, target = "directory_id", directoryID
	default:
		return results
	}
	return slices.DeleteFunc(results, func(result QueryResult) bool {
		identity := firstString(result.Identity, key, "id")
		if identity == "" {
			identity = firstString(result.Resource, key, "id")
		}
		return identity != target
	})
}

func filterQueryResults(spec ListResourceSpec, results []QueryResult, subaccountID, directoryID string, selected []string) ([]QueryResult, error) {
	// Preserve the Version 1 JSON-import behavior: explicit inventory entries
	// are matched against every discovered result before automatic exclusions.
	// create-json and unfiltered exports pass no selection and still receive all
	// required post-query exclusions below.
	if len(selected) > 0 {
		return selectQueryResults(spec, results, selected)
	}

	var defaultNames []string
	switch spec.ResourceType {
	case SubaccountRoleType:
		defaultNames = defaultfilter.FetchDefaultRolesBySubaccount(subaccountID)
	case DirectoryRoleType:
		defaultNames = defaultfilter.FetchDefaultRolesByDirectory(directoryID)
	case SubaccountRoleCollectionType:
		defaultNames = defaultfilter.FetchDefaultRoleCollectionsBySubaccount(subaccountID)
	case DirectoryRoleCollectionType:
		defaultNames = defaultfilter.FetchDefaultRoleCollectionsByDirectory(directoryID)
	}
	results = slices.DeleteFunc(results, func(result QueryResult) bool {
		name := firstString(result.Resource, "name")
		if slices.Contains(defaultNames, name) {
			return true
		}
		if spec.ResourceType == SubaccountTrustConfigurationType && result.InventoryValue == "sap.default" {
			return true
		}
		if spec.ResourceType == SubaccountEntitlementType && defaultfilter.IsDefaultEntitlement(firstString(result.Identity, "service_name"), firstString(result.Identity, "plan_name")) {
			return true
		}
		if spec.ResourceType == SubaccountServiceInstanceType {
			return !isSapCPServiceInstance(result.Resource)
		}
		return false
	})
	return results, nil
}

func selectQueryResults(spec ListResourceSpec, results []QueryResult, selected []string) ([]QueryResult, error) {
	available := make([]string, 0, len(results))
	filtered := make([]QueryResult, 0, len(selected))
	for _, result := range results {
		available = append(available, result.InventoryValue)
		if slices.Contains(selected, result.InventoryValue) {
			filtered = append(filtered, result)
		}
	}
	for _, value := range selected {
		if !slices.Contains(available, value) {
			return nil, fmt.Errorf("%s %q was not found; adjust the provided inventory", spec.CommandName, value)
		}
	}
	return filtered, nil
}

func isSapCPServiceInstance(resource map[string]any) bool {
	contextValue := firstString(resource, "context")
	if contextValue == "" {
		return true
	}
	var contextData map[string]any
	if json.Unmarshal([]byte(contextValue), &contextData) != nil {
		return false
	}
	return firstString(contextData, "origin") == "sapcp"
}

func directoryResourceEnabled(resource string, features []string) bool {
	if len(features) == 0 || resource == CmdDirectoryParameter {
		return true
	}
	if resource == CmdEntitlementParameter {
		return slices.Contains(features, DirectoryFeatureEntitlements)
	}
	if resource == CmdRoleParameter || resource == CmdRoleCollectionParameter {
		return slices.Contains(features, DirectoryFeatureRoles)
	}
	return true
}

func firstString(values map[string]any, keys ...string) string {
	for _, key := range keys {
		if value, ok := values[key]; ok && value != nil {
			text := fmt.Sprintf("%v", value)
			if text != "" && text != "<nil>" {
				return text
			}
		}
	}
	return ""
}

func firstMapString(values map[string]any) string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%v", values[key]))
	}
	return strings.Join(parts, "_")
}

func joinNonEmpty(separator string, values ...string) string {
	nonEmpty := values[:0]
	for _, value := range values {
		if value != "" {
			nonEmpty = append(nonEmpty, value)
		}
	}
	return strings.Join(nonEmpty, separator)
}

func stringSlice(value any) []string {
	items, ok := value.([]any)
	if !ok {
		return nil
	}
	result := make([]string, 0, len(items))
	for _, item := range items {
		result = append(result, fmt.Sprintf("%v", item))
	}
	return result
}

func addUserAgent()    { _ = os.Setenv("BTP_APPEND_USER_AGENT", "terraform_exporter_btp") }
func removeUserAgent() { _ = os.Unsetenv("BTP_APPEND_USER_AGENT") }
