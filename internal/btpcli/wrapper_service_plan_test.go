package btpcli

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetServiceDataByPlanId(t *testing.T) {
	tests := []struct {
		name, planName, offeringID, offeringName, failCommand, wantErr string
		wantCalls                                                      int
	}{
		{name: "success", planName: "standard", offeringID: "offering-id", offeringName: "alert-notification", wantCalls: 2},
		{name: "plan failure", failCommand: "services/plan", wantErr: "error getting service plan name", wantCalls: 1},
		{name: "offering failure", planName: "standard", offeringID: "offering-id", failCommand: "services/offering", wantErr: "error getting service offering name", wantCalls: 2},
		{name: "missing plan name", offeringID: "offering-id", wantErr: "missing its name or offering ID", wantCalls: 1},
		{name: "missing offering ID", planName: "standard", wantErr: "missing its name or offering ID", wantCalls: 1},
		{name: "missing offering name", planName: "standard", offeringID: "offering-id", wantErr: "missing its name", wantCalls: 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			client, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				command, id := "services/plan", "plan-id"
				result := map[string]string{"name": tt.planName, "service_offering_id": tt.offeringID}
				if calls == 2 {
					command, id = "services/offering", "offering-id"
					result = map[string]string{"name": tt.offeringName}
				}
				assertCall(t, r, command, ActionGet, map[string]string{"subaccount": "subaccount-id", "id": id})
				if command == tt.failCommand {
					w.Header().Set(HeaderCLIBackendStatus, "403")
					result = map[string]string{"message": "access denied"}
				}
				assert.NoError(t, json.NewEncoder(w).Encode(result))
			}))
			defer srv.Close()
			plan, offering, err := GetServiceDataByPlanId(client, "subaccount-id", "plan-id")
			assert.Equal(t, tt.wantCalls, calls)
			if tt.wantErr != "" {
				require.ErrorContains(t, err, tt.wantErr)
				assert.Empty(t, plan)
				assert.Empty(t, offering)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.planName, plan)
				assert.Equal(t, tt.offeringName, offering)
			}
		})
	}
}

func TestGetServiceDataByPlanIdInvalidInput(t *testing.T) {
	_, _, err := GetServiceDataByPlanId(nil, "subaccount-id", "plan-id")
	require.ErrorContains(t, err, "client is unavailable")
	for _, ids := range []string{"/plan-id", "subaccount-id/"} {
		parts := strings.Split(ids, "/")
		_, _, err := GetServiceDataByPlanId(&ClientFacade{}, parts[0], parts[1])
		require.ErrorContains(t, err, "IDs are required")
	}
}
