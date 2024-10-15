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
	CmdEntitlementParameter         string = "entitlements"
	CmdEnvironmentInstanceParameter string = "environment-instances"
	CmdSubscriptionParameter        string = "subscriptions"
	CmdTrustConfigurationParameter  string = "trust-configurations"
	CmdRoleParameter                string = "roles"
	CmdRoleCollectionParameter      string = "role-collections"
	CmdServiceInstanceParameter     string = "service-instances"
	CmdServiceBindingParameter      string = "service-bindings"
	CmdSecuritySettingParameter     string = "security-settings"
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

	_, iD := GetExecutionLevelAndId(subaccountId, directoryId)

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
		if level == SubaccountLevel {
			return SubaccountEntitlementType
		} else if level == DirectoryLevel {
			return DirectoryEntitlementType
		}
	case CmdEnvironmentInstanceParameter:
		return SubaccountEnvironmentInstanceType
	case CmdSubscriptionParameter:
		return SubaccountSubscriptionType
	case CmdTrustConfigurationParameter:
		return SubaccountTrustConfigurationType
	case CmdRoleParameter:
		if level == SubaccountLevel {
			return SubaccountRoleType
		} else if level == DirectoryLevel {
			return DirectoryRoleType
		}
	case CmdRoleCollectionParameter:
		if level == SubaccountLevel {
			return SubaccountRoleCollectionType
		} else if level == DirectoryLevel {
			return DirectoryRoleCollectionType
		}
	case CmdServiceInstanceParameter:
		return SubaccountServiceInstanceType
	case CmdServiceBindingParameter:
		return SubaccountServiceBindingType
	case CmdSecuritySettingParameter:
		return SubaccountSecuritySettingType
	case CmdDirectoryParameter:
		return DirectoryType
	}
	return ""
}

func ReadDataSources(subaccountId string, directoryId string, resourceList []string) (btpResources BtpResources, err error) {

	var btpResourcesList []BtpResource
	for _, resource := range resourceList {

		//TODO for Directories: Get the features of the directory
		// introduce a continue depending on the supported features and the iterated resource type

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

func getTfStateData(configDir string, resourceName string, identifier string) ([]byte, error) {
	execPath, err := exec.LookPath("terraform")
	if err != nil {
		fmt.Print("\r\n")
		log.Fatalf("error finding Terraform: %v", err)
		return nil, err
	}
	// create a new Terraform instance
	tf, err := tfexec.NewTerraform(configDir, execPath)
	if err != nil {
		fmt.Print("\r\n")
		log.Fatalf("error running NewTerraform: %v", err)
		return nil, err
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		fmt.Print("\r\n")
		log.Fatalf("error running Init: %v", err)
		return nil, err
	}
	err = tf.Apply(context.Background())
	if err != nil {
		err = handleNotFoundError(err, resourceName, identifier)
		fmt.Print("\r\n")
		log.Fatalf("error running Apply: %v", err)
		return nil, err
	}

	state, err := tf.Show(context.Background())
	if err != nil {
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
		for key := range data {
			key := strings.Replace(key, ":", "_", -1)
			stringArr = append(stringArr, key)
		}
	case SubaccountSubscriptionType:
		subscriptions := data["values"].([]interface{})
		for _, value := range subscriptions {
			subscription := value.(map[string]interface{})
			if fmt.Sprintf("%v", subscription["state"]) != "NOT_SUBSCRIBED" {
				stringArr = append(stringArr, output.FormatSubscriptionResourceName(fmt.Sprintf("%v", subscription["app_name"]), fmt.Sprintf("%v", subscription["plan_name"])))
			}
		}
	case SubaccountEnvironmentInstanceType:
		environmentInstances := data["values"].([]interface{})
		for _, value := range environmentInstances {
			environmentInstance := value.(map[string]interface{})
			stringArr = append(stringArr, fmt.Sprintf("%v", environmentInstance["environment_type"]))
		}
	case SubaccountTrustConfigurationType:
		trusts := data["values"].([]interface{})
		for _, value := range trusts {
			trust := value.(map[string]interface{})
			stringArr = append(stringArr, fmt.Sprintf("%v", trust["origin"]))
		}
	case SubaccountRoleType, DirectoryRoleType:
		roles := data["values"].([]interface{})
		for _, value := range roles {
			role := value.(map[string]interface{})
			stringArr = append(stringArr, output.FormatResourceNameGeneric(fmt.Sprintf("%v", role["name"])))
		}
	case SubaccountRoleCollectionType, DirectoryRoleCollectionType:
		roleCollections := data["values"].([]interface{})
		for _, value := range roleCollections {
			roleCollection := value.(map[string]interface{})
			stringArr = append(stringArr, output.FormatResourceNameGeneric(fmt.Sprintf("%v", roleCollection["name"])))
		}
	case SubaccountServiceInstanceType:
		instances := data["values"].([]interface{})
		for _, value := range instances {
			instance := value.(map[string]interface{})
			stringArr = append(stringArr, output.FormatServiceInstanceResourceName(fmt.Sprintf("%v", instance["name"]), fmt.Sprintf("%v", instance["serviceplan_id"])))
		}
	case SubaccountServiceBindingType:
		bindings := data["values"].([]interface{})
		for _, value := range bindings {
			binding := value.(map[string]interface{})
			stringArr = append(stringArr, fmt.Sprintf("%v", binding["name"]))
		}
	case SubaccountSecuritySettingType:
		stringArr = []string{fmt.Sprintf("%v", data["subaccount_id"])}

	}
	return stringArr
}

func generateDataSourcesForList(subaccountId string, directoryId string, resourceName string) ([]string, error) {
	dataBlockFile := filepath.Join(TmpFolder, "main.tf")
	var jsonBytes []byte

	level, iD := GetExecutionLevelAndId(subaccountId, directoryId)

	btpResourceType := TranslateResourceParamToTechnicalName(resourceName, level)

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

	jsonBytes, err = getTfStateData(TmpFolder, btpResourceType, iD)
	if err != nil {
		error := fmt.Errorf("error fetching Terraform data: %s", err)
		return nil, error
	}

	var data map[string]interface{}

	err = json.Unmarshal(jsonBytes, &data)
	if err != nil {
		fmt.Print("\r\n")
		log.Fatalf("error: %s", err)
		return nil, err
	}
	// ToDo surface the features of the directory stored in data["features"].([]interface{}) analogy to subscription in transform method
	return transformDataToStringArray(btpResourceType, data), nil
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

func handleNotFoundError(err error, resourceName string, iD string) error {
	if strings.Contains(err.Error(), "404") {
		// if it is a 404 error it is probably thw wrong ID, so we return a more readible error message
		if resourceName == DirectoryType {
			return fmt.Errorf("the directory with ID %s was not found. Check that the values for directory ID and globalaccount subdomain are valid", iD)

		} else if resourceName == SubaccountType {
			return fmt.Errorf("the subaccount with ID %s was not found. Check that the values for subaccount ID and globalaccount subdomain are valid", iD)
		}
	}
	return err
}
