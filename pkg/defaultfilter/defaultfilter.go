package defaultfilter

import (
	"fmt"
	"slices"

	"github.com/SAP/terraform-exporter-btp/internal/btpcli"
	"github.com/SAP/terraform-exporter-btp/pkg/toggles"
)

func FetchDefaultRoleCollectionsBySubaccount(subaccountId string) []string {
	if toggles.IsRoleCollectionFilterDeactived() {
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
	if toggles.IsRoleCollectionFilterDeactived() {
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
