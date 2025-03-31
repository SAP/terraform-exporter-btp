package defaultfilter

import (
	"fmt"
	"slices"
	"strings"

	"github.com/SAP/terraform-exporter-btp/internal/btpcli"
	"github.com/SAP/terraform-exporter-btp/pkg/toggles"
)

func FetchDefaultRoleCollectionsBySubaccount(subaccountId string) []string {
	if toggles.IsRoleCollectionFilterDeactivated() {
		return []string{}
	}

	// If we run into errors, we will return an empty slice and NOT abort the processing
	defaulRoleCollections := []string{}

	btpClient, err := btpcli.GetLoggedInClient()

	if err != nil {
		return defaulRoleCollections
	}

	defaulRoleCollections, _ = btpcli.GetDefaultRoleCollectionsBySubaccount(subaccountId, btpClient)
	return defaulRoleCollections
}

func FetchDefaultRoleCollectionsByDirectory(directoryId string) []string {
	if toggles.IsRoleCollectionFilterDeactivated() {
		return []string{}
	}
	// If we run into errors, we will return an empty slice and NOT abort the processing
	defaulRoleCollections := []string{}

	btpClient, err := btpcli.GetLoggedInClient()

	if err != nil {
		return defaulRoleCollections
	}

	defaulRoleCollections, _ = btpcli.GetDefaultRoleCollectionsByDirectory(directoryId, btpClient)
	return defaulRoleCollections
}

func IsRoleCollectionInDefaultList(roleCollectionName string, defaultRoleCollections []string) bool {
	if len(defaultRoleCollections) == 0 {
		return false
	}

	if slices.Contains(defaultRoleCollections, roleCollectionName) {
		return true
	} else {
		return false
	}
}

func FilterDefaultRoleCollectionsFromJsonData(subaccountId string, directoryID string, data map[string]any) map[string]any {
	var defaultRoleCollectionNames []string
	const dataSourceListKey = "values"
	const resourceKey = "name"

	if subaccountId != "" {
		defaultRoleCollectionNames = FetchDefaultRoleCollectionsBySubaccount(subaccountId)
	} else if directoryID != "" {
		defaultRoleCollectionNames = FetchDefaultRoleCollectionsByDirectory(directoryID)
	}

	if len(defaultRoleCollectionNames) == 0 {
		return data
	}

	// Filter out the default role collections from the data
	entities := data[dataSourceListKey].([]interface{})

	entities = slices.DeleteFunc(entities, func(value interface{}) bool {
		entity := value.(map[string]interface{})
		return IsRoleCollectionInDefaultList(fmt.Sprintf("%v", entity[resourceKey]), defaultRoleCollectionNames)
	})

	data[dataSourceListKey] = entities
	return data
}

func FetchDefaultRolesBySubaccount(subaccountId string) []string {
	defaulRoles := []string{}

	if toggles.IsRoleFilterDeactivated() {
		return defaulRoles
	}

	// If we run into errors, we will return an empty slice and NOT abort the processing
	btpClient, err := btpcli.GetLoggedInClient()

	if err != nil {
		return defaulRoles
	}

	defaulRoles, _ = btpcli.GetDefaultRolesBySubaccount(subaccountId, btpClient)
	return defaulRoles
}

func FetchDefaultRolesByDirectory(directoryId string) []string {
	defaultRoles := []string{}

	if toggles.IsRoleFilterDeactivated() {
		return defaultRoles
	}

	// If we run into errors, we will return an empty slice and NOT abort the processing
	btpClient, err := btpcli.GetLoggedInClient()

	if err != nil {
		return defaultRoles
	}

	defaultRoles, _ = btpcli.GetDefaultRolesByDirectory(directoryId, btpClient)
	return defaultRoles
}

func IsRoleInDefaultList(roleName string, defaultRoles []string) bool {
	if len(defaultRoles) == 0 {
		return false
	}

	if slices.Contains(defaultRoles, roleName) {
		return true
	} else {
		return false
	}
}

func FilterDefaultRolesFromJsonData(subaccountId string, directoryID string, data map[string]any) map[string]any {
	var defaultRoles []string
	const dataSourceListKey = "values"
	const resourceKey = "name"

	if subaccountId != "" {
		defaultRoles = FetchDefaultRolesBySubaccount(subaccountId)
	} else if directoryID != "" {
		defaultRoles = FetchDefaultRolesByDirectory(directoryID)
	}

	if len(defaultRoles) == 0 {
		return data
	}

	// Filter out the default role collections from the data
	entities := data[dataSourceListKey].([]interface{})

	entities = slices.DeleteFunc(entities, func(value interface{}) bool {
		entity := value.(map[string]interface{})
		return IsRoleInDefaultList(fmt.Sprintf("%v", entity[resourceKey]), defaultRoles)
	})

	data[dataSourceListKey] = entities
	return data
}

func FilterDefaultIdpJsonData(data map[string]any) map[string]any {
	const dataSourceListKey = "values"
	const resourceKey = "origin"

	// Filter out the default role collections from the data
	entities := data[dataSourceListKey].([]interface{})

	entities = slices.DeleteFunc(entities, func(value interface{}) bool {
		entity := value.(map[string]interface{})
		return fmt.Sprintf("%v", entity[resourceKey]) == "sap.default"
	})

	data[dataSourceListKey] = entities
	return data
}

func IsIdpDefaultIdp(origin string) bool {
	return origin == "sap.default"
}

func IsDefaultEntitlement(serviceName string, planName string) bool {
	if toggles.IsEntitlementFilterDeactivated() {
		return false
	}

	return slices.Contains(DefaultEntitlements, struct {
		ServiceName string
		PlanName    string
	}{ServiceName: serviceName, PlanName: planName})
}

func FilterDefaultEntitlementsFromJsonData(data map[string]any) map[string]any {
	if toggles.IsEntitlementFilterDeactivated() {
		return data
	}

	for key := range data {
		entitlement := strings.Split(key, ":")
		serviceName := entitlement[0]
		planName := entitlement[1]
		if IsDefaultEntitlement(serviceName, planName) {
			delete(data, key)
		}
	}
	return data
}
