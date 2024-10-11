package tfutils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	files "github.com/SAP/terraform-exporter-btp/pkg/files"
	output "github.com/SAP/terraform-exporter-btp/pkg/output"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/spf13/viper"
)

// Constants for TF version for Terraform providers e.g. for SAP BTP
const BtpProviderVersion = "v1.7.0"

const (
	SubaccountLevel = "subaccountLevel"
	DirectoryLevel  = "directoryLevel"
)

const (
	CmdDirectoryParameter           string = "directory"
	CmdSubaccountParameter          string = "subaccount"
	CmdEntitlementParameter         string = "sa-entitlements"
	CmdDirEntitlementParameter      string = "dir-entitlements"
	CmdEnvironmentInstanceParameter string = "sa-environment-instances"
	CmdSubscriptionParameter        string = "sa-subscriptions"
	CmdTrustConfigurationParameter  string = "sa-trust-configurations"
	CmdRoleParameter                string = "sa-roles"
	CmdDirRoleParameter             string = "dir-roles"
	CmdRoleCollectionParameter      string = "sa-role-collections"
	CmdDirRoleCollectionParameter   string = "dir-role-collections"
	CmdServiceInstanceParameter     string = "sa-service-instances"
	CmdServiceBindingParameter      string = "sa-service-bindings"
	CmdSecuritySettingParameter     string = "sa-security-settings"
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

const DataSourcesKind DocKind = "data-sources"
const ResourcesKind DocKind = "resources"

type BtpResource struct {
	Name   string
	Values []string
}

type BtpResources struct {
	BtpResources []BtpResource
}

func FetchImportConfiguration(subaccountId string, directoryId string, resourceType string, tmpFolder string) (map[string]interface{}, error) {

	dataBlock, err := readDataSource(subaccountId, directoryId, resourceType)
	if err != nil {
		return nil, fmt.Errorf("error reading data source: %v", err)
	}

	dataBlockFile := filepath.Join(tmpFolder, "main.tf")
	err = files.CreateFileWithContent(dataBlockFile, dataBlock)
	if err != nil {
		return nil, fmt.Errorf("create file %s failed: %v", dataBlockFile, err)
	}

	jsonBytes, err := getTfStateData(tmpFolder, resourceType)
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

func GetDocByResourceName(kind DocKind, resourceName string) (EntityDocs, error) {
	var choice string

	if (kind == ResourcesKind && resourceName != SubaccountSecuritySettingType) || (kind == DataSourcesKind && resourceName == SubaccountType) || (kind == DataSourcesKind && resourceName == DirectoryType) {
		// We need the singular form of the resource name for all resoucres and the subaccount data source
		choice = resourceName
	} else {
		// We need the plural form of the resource name for all other data sources and security setting resource
		choice = resourceName + "s"
	}

	doc, err := GetDocsForResource("SAP", "btp", "btp", kind, choice, BtpProviderVersion, "github.com")
	if err != nil {
		log.Fatalf("read doc failed for %s, %s: %v", kind, choice, err)
		return EntityDocs{}, err
	}

	return doc, nil
}

func TranslateResourceParamToTechnicalName(resource string) string {
	switch resource {
	case CmdSubaccountParameter:
		return SubaccountType
	case CmdEntitlementParameter:
		return SubaccountEntitlementType
	case CmdEnvironmentInstanceParameter:
		return SubaccountEnvironmentInstanceType
	case CmdSubscriptionParameter:
		return SubaccountSubscriptionType
	case CmdTrustConfigurationParameter:
		return SubaccountTrustConfigurationType
	case CmdRoleParameter:
		return SubaccountRoleType
	case CmdRoleCollectionParameter:
		return SubaccountRoleCollectionType
	case CmdServiceInstanceParameter:
		return SubaccountServiceInstanceType
	case CmdServiceBindingParameter:
		return SubaccountServiceBindingType
	case CmdSecuritySettingParameter:
		return SubaccountSecuritySettingType
	case CmdDirectoryParameter:
		return DirectoryType
	case CmdDirEntitlementParameter:
		return DirectoryEntitlementType
	case CmdDirRoleParameter:
		return DirectoryRoleType
	case CmdDirRoleCollectionParameter:
		return DirectoryRoleCollectionType
	}
	return ""
}

func ReadDataSources(subaccountId string, directoryId string, resourceList []string) (btpResources BtpResources, err error) {

	var btpResourcesList []BtpResource
	for _, resource := range resourceList {
		values, err := generateDataSourcesForList(subaccountId, directoryId, resource)
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

func readDataSource(subaccountId string, directoryId string, resourceName string) (string, error) {

	doc, err := GetDocByResourceName(DataSourcesKind, resourceName)
	if err != nil {
		return "", err
	}

	var dataBlock string

	level, _ := GetExecutionLevelAndId(subaccountId, directoryId)

	switch level {
	case SubaccountLevel:
		if resourceName == SubaccountType {
			dataBlock = strings.Replace(doc.Import, "The ID of the subaccount", subaccountId, -1)
		} else {
			dataBlock = strings.Replace(doc.Import, doc.Attributes["subaccount_id"], subaccountId, -1)
		}
	case DirectoryLevel:
		if resourceName == DirectoryType {
			dataBlock = strings.Replace(doc.Import, "The ID of the directory.", directoryId, -1)
		} else {
			dataBlock = strings.Replace(doc.Import, doc.Attributes["directory_id"], directoryId, -1)
		}
	}
	return dataBlock, nil
}

func getTfStateData(configDir string, resourceName string) ([]byte, error) {
	execPath, err := exec.LookPath("terraform")
	if err != nil {
		log.Fatalf("error finding Terraform: %v", err)
		return nil, err
	}
	// create a new Terraform instance
	tf, err := tfexec.NewTerraform(configDir, execPath)
	if err != nil {
		log.Fatalf("error running NewTerraform: %v", err)
		return nil, err
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		log.Fatalf("error running Init: %v", err)
		return nil, err
	}
	err = tf.Apply(context.Background())
	if err != nil {
		log.Fatalf("error running Apply: %v", err)
		return nil, err
	}

	state, err := tf.Show(context.Background())
	if err != nil {
		log.Fatalf("error running Show: %v", err)
		return nil, err
	}

	// distinguish if the resourceName is entitlelement or different via case
	var jsonBytes []byte
	switch resourceName {
	case SubaccountEntitlementType:
		jsonBytes, err = json.Marshal(state.Values.RootModule.Resources[0].AttributeValues["values"])
	default:
		jsonBytes, err = json.Marshal(state.Values.RootModule.Resources[0].AttributeValues)
	}

	if err != nil {
		log.Fatalf("error json.Marshal: %v", err)
		return nil, err
	}

	return jsonBytes, nil
}

func transformDataToStringArray(btpResource string, data map[string]interface{}) []string {
	var stringArr []string

	switch btpResource {
	case CmdSubaccountParameter:
		stringArr = []string{fmt.Sprintf("%v", data["name"])}
	case CmdDirectoryParameter:
		stringArr = []string{fmt.Sprintf("%v", data["name"])}
	case CmdEntitlementParameter:
		for key := range data {
			key := strings.Replace(key, ":", "_", -1)
			stringArr = append(stringArr, key)
		}
	case CmdSubscriptionParameter:
		subscriptions := data["values"].([]interface{})
		for _, value := range subscriptions {
			subscription := value.(map[string]interface{})
			if fmt.Sprintf("%v", subscription["state"]) != "NOT_SUBSCRIBED" {
				stringArr = append(stringArr, output.FormatSubscriptionResourceName(fmt.Sprintf("%v", subscription["app_name"]), fmt.Sprintf("%v", subscription["plan_name"])))
			}
		}
	case CmdEnvironmentInstanceParameter:
		environmentInstances := data["values"].([]interface{})
		for _, value := range environmentInstances {
			environmentInstance := value.(map[string]interface{})
			stringArr = append(stringArr, fmt.Sprintf("%v", environmentInstance["environment_type"]))
		}
	case CmdTrustConfigurationParameter:
		trusts := data["values"].([]interface{})
		for _, value := range trusts {
			trust := value.(map[string]interface{})
			stringArr = append(stringArr, fmt.Sprintf("%v", trust["origin"]))
		}
	case CmdRoleParameter:
		roles := data["values"].([]interface{})
		for _, value := range roles {
			role := value.(map[string]interface{})
			stringArr = append(stringArr, output.FormatResourceNameGeneric(fmt.Sprintf("%v", role["name"])))
		}
	case CmdRoleCollectionParameter:
		roleCollections := data["values"].([]interface{})
		for _, value := range roleCollections {
			roleCollection := value.(map[string]interface{})
			stringArr = append(stringArr, output.FormatResourceNameGeneric(fmt.Sprintf("%v", roleCollection["name"])))
		}
	case CmdServiceInstanceParameter:
		instances := data["values"].([]interface{})
		for _, value := range instances {
			instance := value.(map[string]interface{})
			stringArr = append(stringArr, output.FormatServiceInstanceResourceName(fmt.Sprintf("%v", instance["name"]), fmt.Sprintf("%v", instance["serviceplan_id"])))
		}
	case CmdServiceBindingParameter:
		bindings := data["values"].([]interface{})
		for _, value := range bindings {
			binding := value.(map[string]interface{})
			stringArr = append(stringArr, fmt.Sprintf("%v", binding["name"]))
		}
	case CmdSecuritySettingParameter:
		stringArr = []string{fmt.Sprintf("%v", data["subaccount_id"])}

	}
	return stringArr
}

func generateDataSourcesForList(subaccountId string, directoryId string, resourceName string) ([]string, error) {
	dataBlockFile := filepath.Join(TmpFolder, "main.tf")
	var jsonBytes []byte

	btpResourceType := TranslateResourceParamToTechnicalName(resourceName)

	dataBlock, err := readDataSource(subaccountId, directoryId, btpResourceType)
	if err != nil {
		error := fmt.Errorf("error reading data source: %s", err)
		return nil, error
	}

	err = files.CreateFileWithContent(dataBlockFile, dataBlock)
	if err != nil {
		error := fmt.Errorf("error creating file %s", dataBlockFile)
		return nil, error
	}

	jsonBytes, err = getTfStateData(TmpFolder, btpResourceType)
	if err != nil {
		error := fmt.Errorf("error fetching Terraform data: %s", err)
		return nil, error
	}

	var data map[string]interface{}

	err = json.Unmarshal(jsonBytes, &data)
	if err != nil {
		log.Fatalf("error: %s", err)
		return nil, err
	}

	return transformDataToStringArray(resourceName, data), nil
}

func runTerraformCommand(args ...string) error {

	verbose := viper.GetViper().GetBool("verbose")
	cmd := exec.Command("terraform", args...)
	if verbose {
		cmd.Stdout = os.Stdout
	} else {
		cmd.Stdout = nil
	}

	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func GetExecutionLevelAndId(subaccountID string, directoryID string) (level string, id string) {
	if subaccountID != "" {
		return SubaccountLevel, subaccountID
	} else if directoryID != "" {
		return DirectoryLevel, directoryID
	}
	return "", ""
}
