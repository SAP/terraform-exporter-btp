package cmd

import (
	"btptfexport/tfutils"
	"context"
	"encoding/json"
	"log"
	"os/exec"
	"strings"

	"github.com/hashicorp/terraform-exec/tfexec"
)

// Constants for TF version for Terraform providers e.g. for SAP BTP
const BtpProviderVersion = "v1.6.0"

type ResourceName string

const (
	SubaccountType                    ResourceName = "btp_subaccount"
	SubaccountEntitlementType         ResourceName = "btp_subaccount_entitlement"
	SubaccountEnvironmentInstanceType ResourceName = "btp_subaccount_environment_instance"
	SubaccountSubscriptionType        ResourceName = "btp_subaccount_subscription"
	SubaccountTrustConfigurationType  ResourceName = "btp_subaccount_trust_configuration"
)

const DataSourcesKind tfutils.DocKind = "data-sources"
const ResourcesKind tfutils.DocKind = "resources"

func GetTfStateData(configDir string, resourceName ResourceName) ([]byte, error) {
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

func getDocByResourceName(kind tfutils.DocKind, resourceName ResourceName) (tfutils.EntityDocs, error) {
	var choice string

	if kind == ResourcesKind || (kind == DataSourcesKind && resourceName == SubaccountType) {
		// We need the singular form of the resource name for all resoucres and the subaccount data source
		choice = string(resourceName)
	} else {
		// We need the plural form of the resource name for all other data sources
		choice = string(resourceName) + "s"
	}

	doc, err := tfutils.GetDocsForResource("SAP", "btp", "btp", kind, choice, BtpProviderVersion, "github.com")
	if err != nil {
		log.Fatalf("read doc failed for %s, %s: %v", kind, choice, err)
		return tfutils.EntityDocs{}, err
	}

	return doc, nil
}

func readDataSource(subaccountId string, resourceName ResourceName) (string, error) {

	doc, err := getDocByResourceName(DataSourcesKind, resourceName)
	if err != nil {
		return "", err
	}

	var dataBlock string
	if resourceName == SubaccountType {
		dataBlock = strings.Replace(doc.Import, "The ID of the subaccount", subaccountId, -1)
	} else {
		dataBlock = strings.Replace(doc.Import, doc.Attributes["subaccount_id"], subaccountId, -1)
	}
	return dataBlock, nil
}

func translateResourceParamToTechnicalName(resource string) ResourceName {
	switch resource {
	case "subaccount":
		return SubaccountType
	case "entitlements":
		return SubaccountEntitlementType
	case "environment-instances":
		return SubaccountEnvironmentInstanceType
	case "subscriptions":
		return SubaccountSubscriptionType
	case "trust-configurations":
		return SubaccountTrustConfigurationType
	}
	return ""
}
