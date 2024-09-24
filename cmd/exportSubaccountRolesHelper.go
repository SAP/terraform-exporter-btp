package cmd

import (
	"btptfexport/files"
	"btptfexport/output"
	"btptfexport/tfutils"
	"fmt"
	"log"
	"slices"
	"strings"
)

func exportSubaccountRoles(subaccountID string, configDir string, filterValues []string) {

	output.AddNewLine()
	spinner, err := output.StartSpinner("crafting import block for " + strings.ToUpper(tfutils.SubaccountRoleType))
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	data, err := tfutils.FetchImportConfiguration(subaccountID, tfutils.SubaccountRoleType, tfutils.TmpFolder)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	importBlock, err := getSubaccountRolesImportBlock(data, subaccountID, filterValues)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	if len(importBlock) == 0 {
		log.Println("no entitlement found for the given subaccount")
		return
	}

	err = files.WriteImportConfiguration(configDir, tfutils.SubaccountRoleType, importBlock)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	err = output.StopSpinner(spinner)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}
}

func getSubaccountRolesImportBlock(data map[string]interface{}, subaccountId string, filterValues []string) (string, error) {

	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.SubaccountRoleType)
	if err != nil {
		log.Fatalf("read doc failed!")
		return "", err
	}

	var importBlock string
	roles := data["values"].([]interface{})

	if len(filterValues) != 0 {
		var subaccountAllRoles []string

		for _, value := range roles {
			role := value.(map[string]interface{})
			resourceName := output.FormatRoleResourceName(fmt.Sprintf("%v", role["name"]))
			subaccountAllRoles = append(subaccountAllRoles, resourceName)
			if slices.Contains(filterValues, resourceName) {
				importBlock += templateEnvironmentInstanceImport(role, subaccountId, resourceDoc)
			}
		}

		missingRole, subset := isSubset(subaccountAllRoles, filterValues)

		if !subset {
			return "", fmt.Errorf("role %s not found in the subaccount. Please adjust it in the provided file", missingRole)
		}

	} else {
		for _, value := range roles {
			role := value.(map[string]interface{})
			importBlock += templateRoleImport(role, subaccountId, resourceDoc)
		}
	}

	return importBlock, nil
}

func templateRoleImport(role map[string]interface{}, subaccountId string, resourceDoc tfutils.EntityDocs) string {

	resourceDoc.Import = strings.Replace(resourceDoc.Import, "'", "", -1)
	resourceName := output.FormatRoleResourceName(fmt.Sprintf("%v", role["name"]))
	template := strings.Replace(resourceDoc.Import, "<resource_name>", resourceName, -1)
	template = strings.Replace(template, "<subaccount_id>", subaccountId, -1)
	template = strings.Replace(template, "<name>", fmt.Sprintf("%v", role["name"]), -1)
	template = strings.Replace(template, "<role_template_name>", fmt.Sprintf("%v", role["role_template_name"]), -1)
	template = strings.Replace(template, "<app_id>", fmt.Sprintf("%v", role["app_id"]), -1)
	return template + "\n"
}
