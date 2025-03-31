package tfutils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/SAP/terraform-exporter-btp/internal/cfcli"
	"github.com/SAP/terraform-exporter-btp/pkg/defaultfilter"
	files "github.com/SAP/terraform-exporter-btp/pkg/files"
	output "github.com/SAP/terraform-exporter-btp/pkg/output"
)

// Constants for TF version for Terraform providers
const BtpProviderVersion = "v1.11.0"
const CfProviderVersion = "v1.4.0"

const (
	SubaccountLevel   = "subaccountLevel"
	DirectoryLevel    = "directoryLevel"
	OrganizationLevel = "organizationLevel"
	SpaceLevel        = "spaceLevel"
)

const (
	CmdDirectoryParameter           string = "directory"
	CmdSubaccountParameter          string = "subaccount"
	CmdEntitlementParameter         string = "entitlements"
	CmdEnvironmentInstanceParameter string = "environment-instances"
	CmdSubscriptionParameter        string = "subscriptions"
	CmdTrustConfigurationParameter  string = "trust-configurations"
	CmdRoleParameter                string = "roles"
	CmdRoleCollectionParameter      string = "role-collections"
	CmdServiceInstanceParameter     string = "service-instances"
	CmdServiceBindingParameter      string = "service-bindings"
	CmdSecuritySettingParameter     string = "security-settings"
	CmdCfSpaceParameter             string = "spaces"
	CmdCfUserParameter              string = "users_cf"
	CmdCfDomainParamater            string = "domains"
	CmdCfOrgRoleParameter           string = "org-roles"
	CmdCfRouteParameter             string = "routes"
	CmdCfSpaceQuotaParameter        string = "space-quotas"
	CmdCfServiceInstanceParameter   string = "cf-service-instances"
	CmdCfSpaceRoleParameter         string = "space-roles"
)

const (
	SubaccountType                    string = "btp_subaccount"
	SubaccountEntitlementType         string = "btp_subaccount_entitlement"
	SubaccountEnvironmentInstanceType string = "btp_subaccount_environment_instance"
	SubaccountSubscriptionType        string = "btp_subaccount_subscription"
	SubaccountTrustConfigurationType  string = "btp_subaccount_trust_configuration"
	SubaccountRoleType                string = "btp_subaccount_role"
	SubaccountRoleCollectionType      string = "btp_subaccount_role_collection"
	SubaccountServiceInstanceType     string = "btp_subaccount_service_instance"
	SubaccountServiceBindingType      string = "btp_subaccount_service_binding"
	SubaccountSecuritySettingType     string = "btp_subaccount_security_setting"
)

const (
	DirectoryType               string = "btp_directory"
	DirectoryEntitlementType    string = "btp_directory_entitlement"
	DirectoryRoleType           string = "btp_directory_role"
	DirectoryRoleCollectionType string = "btp_directory_role_collection"
)

const (
	CfSpaceType           string = "cloudfoundry_space"
	CfUserType            string = "cloudfoundry_user_cf"
	CfUserTypeRead        string = "cloudfoundry_user"
	CfOrgRoleType         string = "cloudfoundry_org_role"
	CfDomainType          string = "cloudfoundry_domain"
	CfRouteType           string = "cloudfoundry_route"
	CfSpaceQuotaType      string = "cloudfoundry_space_quota"
	CfServiceInstanceType string = "cloudfoundry_service_instance"
	CfSpaceRoleType       string = "cloudfoundry_space_role"
)

const DirectoryFeatureDefault string = "DEFAULT"
const DirectoryFeatureEntitlements string = "ENTITLEMENTS"
const DirectoryFeatureRoles string = "AUTHORIZATIONS"

const DataSourcesKind DocKind = "data-sources"
const ResourcesKind DocKind = "resources"

type BtpResource struct {
	Name   string
	Values []string
}

type BtpResources struct {
	BtpResources []BtpResource
}

func FetchImportConfiguration(subaccountId string, directoryId string, organizationId string, spaceId string, resourceType string, tmpFolder string) (map[string]interface{}, error) {

	dataBlock, err := readDataSource(subaccountId, directoryId, organizationId, spaceId, resourceType)
	if err != nil {
		return nil, fmt.Errorf("error reading data source: %v", err)
	}

	dataBlockFile := filepath.Join(tmpFolder, "main.tf")
	err = files.CreateFileWithContent(dataBlockFile, dataBlock)
	if err != nil {
		return nil, fmt.Errorf("create file %s failed: %v", dataBlockFile, err)
	}

	_, iD := GetExecutionLevelAndId(subaccountId, directoryId, organizationId, spaceId)

	jsonBytes, err := getTfStateData(tmpFolder, resourceType, iD)
	if err != nil {
		return nil, fmt.Errorf("error getting Terraform state data: %v", err)
	}

	var data map[string]interface{}
	err = json.Unmarshal(jsonBytes, &data)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return data, nil
}

func GetDocByResourceName(kind DocKind, resourceName string, level string) (EntityDocs, error) {
	var choice string

	// Special handling for resource cloudfoundry_user_cf as the corresponding data source is cloudfoundry_user
	if resourceName == CfUserType && kind == DataSourcesKind {
		resourceName = CfUserTypeRead
	}

	if (kind == ResourcesKind && resourceName != SubaccountSecuritySettingType) || (kind == DataSourcesKind && resourceName == SubaccountType) || (kind == DataSourcesKind && resourceName == DirectoryType) {
		// We need the singular form of the resource name for all resoucres and the subaccount data source
		choice = resourceName
	} else {
		// We need the plural form of the resource name for all other data sources and security setting resource
		choice = resourceName + "s"
	}

	var ghOrg string
	var provider string
	var resourcePrefix string
	var providerVersion string

	if level == OrganizationLevel || level == SpaceLevel {
		ghOrg = "cloudfoundry"
		provider = "cloudfoundry"
		resourcePrefix = "cloudfoundry"
		providerVersion = CfProviderVersion
	} else {
		ghOrg = "SAP"
		provider = "btp"
		resourcePrefix = "btp"
		providerVersion = BtpProviderVersion
	}

	doc, err := GetDocsForResource(ghOrg, provider, resourcePrefix, kind, choice, providerVersion, "github.com")
	if err != nil {
		fmt.Print("\r\n")
		log.Fatalf("read doc failed for %s, %s: %v", kind, choice, err)
		return EntityDocs{}, err
	}

	return doc, nil
}

func TranslateResourceParamToTechnicalName(resource string, level string) string {
	switch resource {
	case CmdSubaccountParameter:
		return SubaccountType
	case CmdEntitlementParameter:
		return getEntitlementTechnicalNameByLevel(level)
	case CmdEnvironmentInstanceParameter:
		return SubaccountEnvironmentInstanceType
	case CmdSubscriptionParameter:
		return SubaccountSubscriptionType
	case CmdTrustConfigurationParameter:
		return SubaccountTrustConfigurationType
	case CmdRoleParameter:
		return getRoleTechnicalNameByLevel(level)
	case CmdRoleCollectionParameter:
		return getRoleCollectionTechnicalNameByLevel(level)
	case CmdServiceInstanceParameter:
		return SubaccountServiceInstanceType
	case CmdServiceBindingParameter:
		return SubaccountServiceBindingType
	case CmdSecuritySettingParameter:
		return SubaccountSecuritySettingType
	case CmdDirectoryParameter:
		return DirectoryType
	case CmdCfSpaceParameter:
		return CfSpaceType
	case CmdCfUserParameter:
		return CfUserType
	case CmdCfDomainParamater:
		return CfDomainType
	case CmdCfOrgRoleParameter:
		return CfOrgRoleType
	case CmdCfRouteParameter:
		return CfRouteType
	case CmdCfSpaceQuotaParameter:
		return CfSpaceQuotaType
	case CmdCfServiceInstanceParameter:
		return CfServiceInstanceType
	case CmdCfSpaceRoleParameter:
		return CfSpaceRoleType
	}
	return ""
}

func getEntitlementTechnicalNameByLevel(level string) string {
	switch level {
	case SubaccountLevel:
		return SubaccountEntitlementType
	case DirectoryLevel:
		return DirectoryEntitlementType
	default:
		return ""
	}
}

func getRoleTechnicalNameByLevel(level string) string {
	switch level {
	case SubaccountLevel:
		return SubaccountRoleType
	case DirectoryLevel:
		return DirectoryRoleType
	default:
		return ""
	}
}

func getRoleCollectionTechnicalNameByLevel(level string) string {
	switch level {
	case SubaccountLevel:
		return SubaccountRoleCollectionType
	case DirectoryLevel:
		return DirectoryRoleCollectionType
	default:
		return ""
	}
}

func ReadDataSources(subaccountId string, directoryId string, organizationId string, resourceList []string) (btpResources BtpResources, err error) {

	var btpResourcesList []BtpResource
	var featureListMemory []string

	level, _ := GetExecutionLevelAndId(subaccountId, directoryId, organizationId, "")

	for _, resource := range resourceList {
		var values []string
		var featureList []string
		var err error
		var spaceId = ""

		if !resourceIsProcessable(level, resource, featureListMemory) {
			continue
		}

		if resource == CmdCfSpaceRoleParameter {
			var spaces map[string]string
			spaces, err = cfcli.GetSpaceList(organizationId)
			for _, spaceID := range spaces {
				var spaceRoles []string
				spaceRoles, featureList, err = generateDataSourcesForList(subaccountId, directoryId, organizationId, spaceID, resource)
				values = append(values, spaceRoles...)
			}
		} else {
			values, featureList, err = generateDataSourcesForList(subaccountId, directoryId, organizationId, spaceId, resource)
		}

		if resource == CmdDirectoryParameter {
			// Store the features of the directory for later use
			featureListMemory = featureList
		}

		if err != nil {
			error := fmt.Errorf("error generating data sources: %v", err)
			return BtpResources{}, error
		}

		if len(values) != 0 {
			// Only append existing resources to avoid confusion
			data := BtpResource{Name: resource, Values: values}
			btpResourcesList = append(btpResourcesList, data)
		}
	}

	btpResources = BtpResources{BtpResources: btpResourcesList}
	return btpResources, nil
}

func readDataSource(subaccountId string, directoryId string, organizationId string, spaceId string, resourceName string) (string, error) {
	level, _ := GetExecutionLevelAndId(subaccountId, directoryId, organizationId, spaceId)

	doc, err := GetDocByResourceName(DataSourcesKind, resourceName, level)
	if err != nil {
		return "", err
	}

	var dataBlock string

	switch level {
	case SubaccountLevel:
		if resourceName == SubaccountType {
			dataBlock = strings.ReplaceAll(doc.Import, "The ID of the subaccount", subaccountId)
		} else {
			dataBlock = strings.ReplaceAll(doc.Import, doc.Attributes["subaccount_id"], subaccountId)
		}
	case DirectoryLevel:
		if resourceName == DirectoryType {
			dataBlock = strings.ReplaceAll(doc.Import, "The ID of the directory.", directoryId)
		} else {
			dataBlock = strings.ReplaceAll(doc.Import, doc.Attributes["directory_id"], directoryId)
		}
	case OrganizationLevel:
		if resourceName == CfUserType || resourceName == CfUserTypeRead || resourceName == CfDomainType || resourceName == CfRouteType || resourceName == CfServiceInstanceType {
			dataBlock = strings.ReplaceAll(doc.Import, "The ID of the organization", organizationId)
		} else {
			dataBlock = strings.ReplaceAll(doc.Import, doc.Attributes["org"], organizationId)
		}
	case SpaceLevel:
		if resourceName == CfSpaceRoleType {
			dataBlock = strings.ReplaceAll(doc.Import, doc.Attributes["space"], spaceId)
		}
	}

	return dataBlock, nil
}

func getTfStateData(configDir string, resourceName string, identifier string) ([]byte, error) {

	chDir := fmt.Sprintf("-chdir=%s", configDir)
	// Set custom user agent for call of TF Provider via exporter
	addUserAgent()
	defer removeUserAgent()

	err := runTfCmdGeneric(chDir, "init", "-upgrade")
	if err != nil {
		removeUserAgent()
		fmt.Print("\r\n")
		log.Fatalf("error running Init: %v", err)
		return nil, err
	}

	err = runTfCmdGeneric(chDir, "apply", "-auto-approve")
	if err != nil {
		err = handleNotFoundError(err, resourceName, identifier)
		removeUserAgent()
		fmt.Print("\r\n")
		log.Fatalf("error running Apply: %v", err)
		return nil, err
	}

	state, err := runTfShowJson(configDir)
	if err != nil {
		removeUserAgent()
		fmt.Print("\r\n")
		log.Fatalf("error running Show: %v", err)
		return nil, err
	}

	// distinguish if the resourceName is entitlelement or different via case
	var jsonBytes []byte
	switch resourceName {
	case SubaccountEntitlementType, DirectoryEntitlementType:
		jsonBytes, err = json.Marshal(state.Values.RootModule.Resources[0].AttributeValues["values"])
	default:
		jsonBytes, err = json.Marshal(state.Values.RootModule.Resources[0].AttributeValues)
	}

	if err != nil {
		removeUserAgent()
		fmt.Print("\r\n")
		log.Fatalf("error json.Marshal: %v", err)
		return nil, err
	}

	return jsonBytes, nil
}

func transformDataToStringArray(btpResource string, data map[string]interface{}) []string {
	var stringArr []string

	switch btpResource {
	case SubaccountType:
		stringArr = []string{fmt.Sprintf("%v", data["name"])}
	case DirectoryType:
		stringArr = []string{fmt.Sprintf("%v", data["name"])}
	case SubaccountEntitlementType, DirectoryEntitlementType:
		transformEntitlementStringArray(data, &stringArr)
	case SubaccountSubscriptionType:
		transformSubscriptionsStringArray(data, &stringArr)
	case SubaccountEnvironmentInstanceType:
		transformDataToStringArrayGeneric(data, &stringArr, "values", "environment_type")
	case SubaccountTrustConfigurationType:
		transformDataToStringArrayGeneric(data, &stringArr, "values", "origin")
	case SubaccountRoleType, DirectoryRoleType:
		transformDataToStringArrayGeneric(data, &stringArr, "values", "name")
	case SubaccountRoleCollectionType, DirectoryRoleCollectionType:
		transformDataToStringArrayGeneric(data, &stringArr, "values", "name")
	case SubaccountServiceInstanceType:
		transformServiceInstanceStringArray(data, &stringArr)
	case SubaccountServiceBindingType:
		transformDataToStringArrayGeneric(data, &stringArr, "values", "name")
	case SubaccountSecuritySettingType:
		stringArr = []string{fmt.Sprintf("%v", data["subaccount_id"])}
	case CfSpaceType:
		transformDataToStringArrayGeneric(data, &stringArr, "spaces", "name")
	case CfUserType, CfUserTypeRead:
		transformDataToStringArrayGeneric(data, &stringArr, "users", "username")
	case CfDomainType:
		transformDataToStringArrayGeneric(data, &stringArr, "domains", "name")
	case CfOrgRoleType:
		transformOrgRolesStringArray(data, &stringArr)
	case CfRouteType:
		transformDataToStringArrayGeneric(data, &stringArr, "routes", "url")
	case CfSpaceQuotaType:
		transformDataToStringArrayGeneric(data, &stringArr, "space_quotas", "name")
	case CfServiceInstanceType:
		transformCfServiceInstanceStringArray(data, &stringArr)
	case CfSpaceRoleType:
		transformCfSpaceRolesStringArray(data, &stringArr)
	}
	return stringArr
}

func generateDataSourcesForList(subaccountId string, directoryId string, organizationId string, spaceID string, resourceName string) ([]string, []string, error) {
	dataBlockFile := filepath.Join(TmpFolder, "main.tf")
	var jsonBytes []byte

	level, iD := GetExecutionLevelAndId(subaccountId, directoryId, organizationId, spaceID)

	btpResourceType := TranslateResourceParamToTechnicalName(resourceName, level)

	if btpResourceType == CfUserType {
		// For CF Users we must use the data source cloudfoundry_user, but the resource name is cloudfoundry_user_cf
		btpResourceType = CfUserTypeRead
	}

	dataBlock, err := readDataSource(subaccountId, directoryId, organizationId, spaceID, btpResourceType)
	if err != nil {
		error := fmt.Errorf("error reading data source: %s", err)
		return nil, nil, error
	}

	err = files.CreateFileWithContent(dataBlockFile, dataBlock)
	if err != nil {
		error := fmt.Errorf("error creating file %s", dataBlockFile)
		return nil, nil, error
	}

	jsonBytes, err = getTfStateData(TmpFolder, btpResourceType, iD)
	if err != nil {
		error := fmt.Errorf("error fetching Terraform data: %s", err)
		return nil, nil, error
	}

	var data map[string]interface{}

	err = json.Unmarshal(jsonBytes, &data)
	if err != nil {
		error := fmt.Errorf("error unmarshelling JSON: %s", err)
		return nil, nil, error
	}

	data = filterDefaultValues(subaccountId, directoryId, btpResourceType, data)

	return transformDataToStringArray(btpResourceType, data), extractFeatureList(data, btpResourceType), nil
}

func GetExecutionLevelAndId(subaccountID string, directoryID string, organizationID string, spaceID string) (level string, id string) {
	if subaccountID != "" {
		return SubaccountLevel, subaccountID
	} else if directoryID != "" {
		return DirectoryLevel, directoryID
	} else if organizationID != "" {
		if spaceID != "" {
			return SpaceLevel, spaceID
		}
		return OrganizationLevel, organizationID
	}
	return "", ""
}

func handleNotFoundError(err error, resourceName string, iD string) error {
	if strings.Contains(err.Error(), "404") {
		// if it is a 404 error it is probably thw wrong ID, so we return a more readible error message
		switch resourceName {
		case DirectoryType:
			return fmt.Errorf("the directory with ID %s was not found. Check that the values for directory ID and globalaccount subdomain are valid", iD)
		case SubaccountType:
			return fmt.Errorf("the subaccount with ID %s was not found. Check that the values for subaccount ID and globalaccount subdomain are valid", iD)
		}
	}
	return err
}

func extractFeatureList(data map[string]interface{}, resourceName string) []string {

	featureList := []string{}

	if resourceName == DirectoryType {
		features := data["features"].([]interface{})
		var featureList []string
		for _, feature := range features {
			featureList = append(featureList, fmt.Sprintf("%v", feature.(string)))
		}
		return featureList
	}

	return featureList

}

func resourceIsProcessable(level string, resource string, featureList []string) bool {

	// Only relevant for directory resources due to unmanaged and managed directories
	if level != DirectoryLevel {
		return true
	}

	// Check if the resource is processable based on the feature list
	if resource == CmdEntitlementParameter && !slices.Contains(featureList, DirectoryFeatureEntitlements) {
		return false
	}

	if (resource == CmdRoleParameter || resource == CmdRoleCollectionParameter) && !slices.Contains(featureList, DirectoryFeatureRoles) {
		return false
	}

	return true
}

func transformDataToStringArrayGeneric(data map[string]interface{}, stringArr *[]string, dataSourceListKey string, resourceKey string) {
	entities := data[dataSourceListKey].([]interface{})
	for _, value := range entities {
		entity := value.(map[string]interface{})
		*stringArr = append(*stringArr, output.FormatResourceNameGeneric(fmt.Sprintf("%v", entity[resourceKey])))
	}
}

func transformEntitlementStringArray(data map[string]interface{}, stringArr *[]string) {
	for key := range data {
		key := strings.ReplaceAll(key, ":", "_")
		*stringArr = append(*stringArr, key)
	}
}

func transformServiceInstanceStringArray(data map[string]interface{}, stringArr *[]string) {
	instances := data["values"].([]interface{})
	for _, value := range instances {
		instance := value.(map[string]interface{})
		*stringArr = append(*stringArr, output.FormatServiceInstanceResourceName(fmt.Sprintf("%v", instance["name"]), fmt.Sprintf("%v", instance["serviceplan_id"])))
	}
}

func transformSubscriptionsStringArray(data map[string]interface{}, stringArr *[]string) {
	subscriptions := data["values"].([]interface{})
	for _, value := range subscriptions {
		subscription := value.(map[string]interface{})
		if fmt.Sprintf("%v", subscription["state"]) != "NOT_SUBSCRIBED" {
			*stringArr = append(*stringArr, output.FormatSubscriptionResourceName(fmt.Sprintf("%v", subscription["app_name"]), fmt.Sprintf("%v", subscription["plan_name"])))
		}
	}
}

func transformOrgRolesStringArray(data map[string]interface{}, stringArr *[]string) {
	roles := data["roles"].([]interface{})
	for _, value := range roles {
		role := value.(map[string]interface{})
		*stringArr = append(*stringArr, output.FormatOrgRoleResourceName(fmt.Sprintf("%v", role["type"]), fmt.Sprintf("%v", role["user"])))
	}
}

func transformCfServiceInstanceStringArray(data map[string]interface{}, stringArr *[]string) {
	instances := data["service_instances"].([]interface{})
	for _, value := range instances {
		instance := value.(map[string]interface{})
		*stringArr = append(*stringArr, output.FormatServiceInstanceResourceName(fmt.Sprintf("%v", instance["name"]), fmt.Sprintf("%v", instance["service_plan"])))
	}
}

func transformCfSpaceRolesStringArray(data map[string]interface{}, stringArr *[]string) {
	roles := data["roles"].([]interface{})
	for _, value := range roles {
		role := value.(map[string]interface{})
		*stringArr = append(*stringArr, output.FormatSpaceRoleResourceName(fmt.Sprintf("%v", role["type"]), fmt.Sprintf("%v", role["space"]), fmt.Sprintf("%v", role["user"])))
	}
}

func addUserAgent() {
	_ = os.Setenv("BTP_APPEND_USER_AGENT", "terraform_exporter_btp")
}

func removeUserAgent() {
	_ = os.Unsetenv("BTP_APPEND_USER_AGENT")
}

func filterDefaultValues(subaccountId string, directoryId string, btpResourceType string, data map[string]interface{}) map[string]any {

	switch btpResourceType {
	case SubaccountRoleCollectionType, DirectoryRoleCollectionType:
		return defaultfilter.FilterDefaultRoleCollectionsFromJsonData(subaccountId, directoryId, data)
	case SubaccountRoleType, DirectoryRoleType:
		return defaultfilter.FilterDefaultRolesFromJsonData(subaccountId, directoryId, data)
	case SubaccountTrustConfigurationType:
		return defaultfilter.FilterDefaultIdpJsonData(data)
	case SubaccountEntitlementType:
		return defaultfilter.FilterDefaultEntitlementsFromJsonData(data)
	default:
		return data
	}

}
