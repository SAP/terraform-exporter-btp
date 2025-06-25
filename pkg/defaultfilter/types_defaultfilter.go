package defaultfilter

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
)

//go:embed default-entitlements.json
var servicesJSON []byte

// To Parse the JSON file
type ServiceEntry struct {
	ServiceName string `json:"service"`
	PlanName    string `json:"plan"`
}

// EntitlementFilterData is now a slice of ServiceEntry
type EntitlementFilterData []ServiceEntry

var DefaultEntitlements EntitlementFilterData

func init() {
	err := json.Unmarshal(servicesJSON, &DefaultEntitlements)
	if err != nil {
		fmt.Print("\r\n")
		log.Fatalf("error while unmarshalling default-entitlements: %v", err)
		return
	}
}
