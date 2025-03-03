package testutils

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// Define source files (initial state, target state)
type testFileStruc struct {
	sourceFilePath string
	targetFilePath string
}

type testFileMapping map[string]testFileStruc

var testFiles = testFileMapping{
	"provider_btp": {
		sourceFilePath: "provider_btp_source.tf",
		targetFilePath: "provider_btp_target.tf",
	},
	"provider_cf": {
		sourceFilePath: "provider_cf_source.tf",
		targetFilePath: "provider_cf_target.tf",
	},
	"provider_withbackend_btp": {
		sourceFilePath: "provider_btp_backend_source.tf",
		targetFilePath: "provider_btp_backend_target.tf",
	},
	"provider_wobackend_btp": {
		sourceFilePath: "provider_btp_wo_backend.tf",
		targetFilePath: "provider_btp_wo_backend.tf",
	},
	"sa_trust_config_replace": {
		sourceFilePath: "resource_subaccount_trust_configuration_replace.tf",
		targetFilePath: "resource_subaccount_trust_configuration_replace.tf",
	},
	"sa_trust_config_no_replace": {
		sourceFilePath: "resource_subaccount_trust_configuration_no_replace.tf",
		targetFilePath: "resource_subaccount_trust_configuration_no_replace.tf",
	},
	"sa_with_ga_parent": {
		sourceFilePath: "resource_subaccount_with_ga_parent_source.tf",
		targetFilePath: "resource_subaccount_with_ga_parent_target.tf",
	},
	"sa_without_ga_parent": {
		sourceFilePath: "resource_subaccount_wo_ga_parent_source.tf",
		targetFilePath: "resource_subaccount_wo_ga_parent_target.tf",
	},
	"sa_entitlement": {
		sourceFilePath: "resource_subaccount_entitlement.tf",
		targetFilePath: "resource_subaccount_entitlement.tf",
	},
	"sa_entitlement_error": {
		sourceFilePath: "resource_subaccount_entitlement_incomplete.tf",
		targetFilePath: "resource_subaccount_entitlement_incomplete.tf",
	},
	"sa_role": {
		sourceFilePath: "resource_subaccount_role.tf",
		targetFilePath: "resource_subaccount_role.tf",
	},
	"sa_role_error": {
		sourceFilePath: "resource_subaccount_role_incomplete.tf",
		targetFilePath: "resource_subaccount_role_incomplete.tf",
	},
	"sa_subscription_dependency": {
		sourceFilePath: "resource_subaccount_subscription_dep_source.tf",
		targetFilePath: "resource_subaccount_subscription_dep_target.tf",
	},
	"sa_subscription_no_dependency": {
		sourceFilePath: "resource_subaccount_subscription_no_dependency.tf",
		targetFilePath: "resource_subaccount_subscription_no_dependency.tf",
	},
	"sa_service_instance_dependency": {
		sourceFilePath: "resource_subaccount_service_instance_dep_source.tf",
		targetFilePath: "resource_subaccount_service_instance_dep_target.tf",
	},
	"sa_service_instance_no_dependency": {
		sourceFilePath: "resource_subaccount_service_instance_no_dependency.tf",
		targetFilePath: "resource_subaccount_service_instance_no_dependency.tf",
	},
	"sa_role_collection_dependency": {
		sourceFilePath: "resource_subaccount_role_collection_dep_source.tf",
		targetFilePath: "resource_subaccount_role_collection_dep_target.tf",
	},
	"sa_role_collection_single_dependency": {
		sourceFilePath: "resource_subaccount_role_collection_single_dep_source.tf",
		targetFilePath: "resource_subaccount_role_collection_single_dep_target.tf",
	},
	"sa_role_collection_no_dependency": {
		sourceFilePath: "resource_subaccount_role_collection_no_dependency.tf",
		targetFilePath: "resource_subaccount_role_collection_no_dependency.tf",
	},
	"empty_attributes": {
		sourceFilePath: "resource_empty_attributes_source.tf",
		targetFilePath: "resource_empty_attributes_target.tf",
	},
	"main_dependency": {
		sourceFilePath: "resource_main_dependency_source.tf",
		targetFilePath: "resource_main_dependency_target.tf",
	},
	"replace_attribute": {
		sourceFilePath: "resource_replace_attribute_source.tf",
		targetFilePath: "resource_replace_attribute_target.tf",
	},
	"remove_import_block": {
		sourceFilePath: "remove_import_block_source.tf",
		targetFilePath: "remove_import_block_target.tf",
	},
	"remove_config_block": {
		sourceFilePath: "remove_config_block_source.tf",
		targetFilePath: "remove_config_block_target.tf",
	},
	"cf_space_dependency": {
		sourceFilePath: "resource_cf_space_dependency_source.tf",
		targetFilePath: "resource_cf_space_dependency_target.tf",
	},
	"cf_space_no_dependency": {
		sourceFilePath: "resource_cf_space_no_dependency.tf",
		targetFilePath: "resource_cf_space_no_dependency.tf",
	},
}

// Read a hcl files from disc
func GetHclFilesById(id string) (sourceHclFile *hclwrite.File, targetHclFile *hclwrite.File) {
	// Read hcl file from disc
	sourceFilePath, targetFilePath := getFilePathsbyId(id)

	currentDir, _ := os.Getwd()
	// We navigate up one level to land in the tfcleanup directory
	currentDir = filepath.Dir(currentDir)

	srcPath := filepath.Join(currentDir, "testutils", "testdata", sourceFilePath)

	srcFile, err := os.ReadFile(srcPath)
	if err != nil {
		log.Printf("Failed to read file %q: %s", srcPath, err)
		return
	}

	trgtPath := filepath.Join(currentDir, "testutils", "testdata", targetFilePath)
	trgtFile, err := os.ReadFile(trgtPath)
	if err != nil {
		log.Printf("Failed to read file %q: %s", trgtPath, err)
		return
	}

	sourceHclFile, diags := hclwrite.ParseConfig(srcFile, srcPath, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		for _, diag := range diags {
			if diag.Subject != nil {
				log.Printf("[%s:%d] %s: %s", diag.Subject.Filename, diag.Subject.Start.Line, diag.Summary, diag.Detail)
			} else {
				log.Printf("%s: %s", diag.Summary, diag.Detail)
			}
		}
		return
	}

	targetHclFile, diags = hclwrite.ParseConfig(trgtFile, trgtPath, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		for _, diag := range diags {
			if diag.Subject != nil {
				log.Printf("[%s:%d] %s: %s", diag.Subject.Filename, diag.Subject.Start.Line, diag.Summary, diag.Detail)
			} else {
				log.Printf("%s: %s", diag.Summary, diag.Detail)
			}
		}
		return
	}

	return sourceHclFile, targetHclFile
}

func AreHclFilesEqual(testResultHclFile *hclwrite.File, targetHclFile *hclwrite.File) error {

	hclwrite.Format(testResultHclFile.Bytes())

	if bytes.Equal(hclwrite.Format(testResultHclFile.Bytes()), hclwrite.Format(targetHclFile.Bytes())) {
		return nil
	} else {
		return fmt.Errorf("HCL files are not equal")
	}

}

func getFilePathsbyId(id string) (sourceFilePath string, targetFilePath string) {
	return testFiles[id].sourceFilePath, testFiles[id].targetFilePath
}

func GetGlobalAccountMockParentData(parentId string) bool {
	return parentId == "GlobalAccountSubdomain"
}

func GetServiceMockData(planId string) (planName string, serviceName string) {
	if planId == "0815" {
		return "standard", "alert-notification"
	}
	return "premium", "alert-notification-dash"
}
