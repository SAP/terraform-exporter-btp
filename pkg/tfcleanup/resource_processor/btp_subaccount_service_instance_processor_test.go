package resourceprocessor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/SAP/terraform-exporter-btp/internal/btpcli"
	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/testutils"
	"github.com/SAP/terraform-exporter-btp/pkg/tfutils"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
)

func TestAddServiceInstanceDependency(t *testing.T) {
	srcFileServiceInstanceDep, trgtFileServiceInstanceDep := testutils.GetHclFilesById("sa_service_instance_dependency")
	srcFileServiceInstanceNoDep, trgtFileServiceInstanceNoDep := testutils.GetHclFilesById("sa_service_instance_no_dependency")

	defaultTestDependencies := getNewServiceInstanceDepTemplate()

	tests := []struct {
		name         string
		src          *hclwrite.File
		trgt         *hclwrite.File
		dependencies generictools.DependencyAddresses
	}{
		{
			name:         "Test Service Instance Dependency",
			src:          srcFileServiceInstanceDep,
			trgt:         trgtFileServiceInstanceDep,
			dependencies: defaultTestDependencies,
		},
		{
			name:         "Test No Service Instance Dependency",
			src:          srcFileServiceInstanceNoDep,
			trgt:         trgtFileServiceInstanceNoDep,
			dependencies: generictools.DependencyAddresses{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			blocks := tt.src.Body().Blocks()
			// we assume one resource entry in the blocks file
			addServiceInstanceDependency(blocks[0].Body(), &tt.dependencies, false)
			assert.NoError(t, testutils.AreHclFilesEqual(tt.src, tt.trgt))
		})
	}

}

func TestAddServiceInstanceDependencyKeepsQueryGeneratedPlanID(t *testing.T) {
	f := hclwrite.NewEmptyFile()
	block := f.Body().AppendNewBlock("resource", []string{serviceInstanceBlockIdentifier, "instance"})
	body := block.Body()
	body.SetAttributeValue("name", cty.StringVal("my-instance"))
	body.SetAttributeValue(serviceInstancePlanIdentifier, cty.StringVal("plan-id"))

	dependencies := generictools.NewDependencyAddresses()
	require.NotPanics(t, func() {
		addServiceInstanceDependency(body, &dependencies, false)
	})

	planID := body.GetAttribute(serviceInstancePlanIdentifier)
	require.NotNil(t, planID)
	assert.Equal(t, "plan-id", generictools.GetStringToken(planID.Expr().BuildTokens(nil)))
}

func getNewServiceInstanceDepTemplate() generictools.DependencyAddresses {
	defaultTestDependencies := generictools.NewDependencyAddresses()

	entitlementKey := generictools.EntitlementKey{
		ServiceName: "alert-notification",
		PlanName:    "standard",
	}

	defaultTestDependencies.EntitlementAddress = make(map[generictools.EntitlementKey]generictools.EntitlementInfo)
	defaultTestDependencies.EntitlementAddress[entitlementKey] = generictools.EntitlementInfo{Address: "btp_subaccount_entitlement.entitlement_0"}
	defaultTestDependencies.SubaccountAddress = "btp_subaccount.subaccount_0"

	return defaultTestDependencies
}

func TestProcessServiceInstance(t *testing.T) {
	tests := []struct {
		name, source, plan, offering string
		lookupErr                    error
		wantCalls                    int
		converted                    bool
	}{
		{name: "ID only", source: "serviceplan_id = \"plan-id\"", plan: "standard", offering: "alert-notification", wantCalls: 1, converted: true},
		{name: "partial names", source: "serviceplan_id = \"plan-id\"\nserviceplan_name = \"stale\"", plan: "standard", offering: "alert-notification", wantCalls: 1, converted: true},
		{name: "partial offering", source: "serviceplan_id = \"plan-id\"\nservice_offering_name = \"stale\"", plan: "standard", offering: "alert-notification", wantCalls: 1, converted: true},
		{name: "existing names", source: "serviceplan_id = \"plan-id\"\nserviceplan_name = \"standard\"\nservice_offering_name = \"alert-notification\"", converted: true},
		{name: "missing ID"},
		{name: "null ID", source: "serviceplan_id = null"},
		{name: "empty ID", source: "serviceplan_id = \"\""},
		{name: "reference ID", source: "serviceplan_id = var.plan_id"},
		{name: "interpolated ID", source: "serviceplan_id = \"${var.plan_id}\""},
		{name: "numeric ID", source: "serviceplan_id = 42"},
		{name: "lookup error", source: "serviceplan_id = \"plan-id\"", lookupErr: fmt.Errorf("lookup failed"), wantCalls: 1},
		{name: "partial names on error", source: "serviceplan_id = \"plan-id\"\nserviceplan_name = \"existing\"", lookupErr: fmt.Errorf("lookup failed"), wantCalls: 1},
		{name: "empty plan", source: "serviceplan_id = \"plan-id\"", offering: "alert-notification", wantCalls: 1},
		{name: "empty offering", source: "serviceplan_id = \"plan-id\"", plan: "standard", wantCalls: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, diags := hclwrite.ParseConfig([]byte(tt.source+"\n"), "test.tf", hcl.InitialPos)
			require.False(t, diags.HasErrors(), diags.Error())
			before := string(f.Bytes())
			calls := 0
			lookup := func(subaccountID, planID string) (string, string, error) {
				calls++
				assert.Equal(t, "subaccount-id", subaccountID)
				assert.Equal(t, "plan-id", planID)
				return tt.plan, tt.offering, tt.lookupErr
			}
			deps := generictools.NewDependencyAddresses()
			processServiceInstance(f.Body(), &deps, "subaccount-id", lookup, false)
			assert.Equal(t, tt.wantCalls, calls)
			if !tt.converted {
				assert.Equal(t, before, string(f.Bytes()))
				return
			}
			assert.Nil(t, f.Body().GetAttribute(serviceInstancePlanIdentifier))
			assert.Equal(t, "standard", serviceInstanceLiteral(f.Body(), serviceInstancePlanNameIdentifier))
			assert.Equal(t, "alert-notification", serviceInstanceLiteral(f.Body(), serviceInstanceOfferingNameIdentifier))
			assert.Nil(t, f.Body().GetAttribute("depends_on"))
			after := string(f.Bytes())
			processServiceInstance(f.Body(), &deps, "subaccount-id", lookup, false)
			assert.Equal(t, tt.wantCalls, calls)
			assert.Equal(t, after, string(f.Bytes()))
		})
	}
}

func TestProcessServiceInstanceEntitlementDependency(t *testing.T) {
	for _, module := range []bool{false, true} {
		t.Run(fmt.Sprint(module), func(t *testing.T) {
			want := "btp_subaccount_entitlement.entitlement_0"
			if module {
				want = genericEntitlementModuleAddress
			}
			f := hclwrite.NewEmptyFile()
			f.Body().SetAttributeValue(serviceInstancePlanIdentifier, cty.StringVal("plan-id"))
			deps := getNewServiceInstanceDepTemplate()
			processServiceInstance(f.Body(), &deps, "subaccount-id", func(_, _ string) (string, string, error) {
				return "standard", "alert-notification", nil
			}, module)
			require.Nil(t, f.Body().GetAttribute(serviceInstancePlanIdentifier))
			dep := f.Body().GetAttribute("depends_on")
			require.NotNil(t, dep)
			assert.Contains(t, string(dep.Expr().BuildTokens(nil).Bytes()), want)
		})
	}
}

func TestProcessServiceInstanceEscapedNames(t *testing.T) {
	f := hclwrite.NewEmptyFile()
	f.Body().SetAttributeValue(serviceInstancePlanIdentifier, cty.StringVal("plan-id"))
	deps := generictools.NewDependencyAddresses()
	plan, offering := "plan\"name", "offering${literal}"
	processServiceInstance(f.Body(), &deps, "subaccount-id", func(_, _ string) (string, string, error) {
		return plan, offering, nil
	}, false)
	assert.Nil(t, f.Body().GetAttribute(serviceInstancePlanIdentifier))
	assert.Equal(t, plan, serviceInstanceLiteral(f.Body(), serviceInstancePlanNameIdentifier))
	assert.Equal(t, offering, serviceInstanceLiteral(f.Body(), serviceInstanceOfferingNameIdentifier))
}

// Exercise the public cleanup entry point with a local CLI server only.
func TestProcessResourcesResolvesServiceInstance(t *testing.T) {
	for _, failOffering := range []bool{false, true} {
		t.Run(fmt.Sprint(failOffering), func(t *testing.T) {
			calls := 0
			srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				var payload struct{ ParamValues map[string]string }
				require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
				assert.Equal(t, "subaccount-id", payload.ParamValues["subaccount"])
				assert.Equal(t, "get", r.URL.RawQuery)
				w.Header().Set(btpcli.HeaderCLIBackendStatus, "200")
				response := map[string]string{}
				switch {
				case strings.HasSuffix(r.URL.Path, "/services/plan"):
					assert.Equal(t, "plan-id", payload.ParamValues["id"])
					response = map[string]string{"name": "standard", "service_offering_id": "offering-id"}
				case strings.HasSuffix(r.URL.Path, "/services/offering"):
					assert.Equal(t, "offering-id", payload.ParamValues["id"])
					response = map[string]string{"name": "alert-notification"}
					if failOffering {
						w.Header().Set(btpcli.HeaderCLIBackendStatus, "403")
						response = map[string]string{"message": "access denied"}
					}
				default:
					t.Errorf("unexpected request: %s", r.URL)
				}
				assert.NoError(t, json.NewEncoder(w).Encode(response))
			}))
			defer srv.Close()
			endpoint, err := url.Parse(srv.URL)
			require.NoError(t, err)
			client := btpcli.NewClientFacade(btpcli.NewV2ClientWithHttpClient(srv.Client(), endpoint))
			f := hclwrite.NewEmptyFile()
			body := f.Body().AppendNewBlock("resource", []string{serviceInstanceBlockIdentifier, "instance"}).Body()
			body.SetAttributeValue("name", cty.StringVal("my-instance"))
			body.SetAttributeValue(serviceInstancePlanIdentifier, cty.StringVal("plan-id"))
			body.SetAttributeValue(serviceInstancePlanNameIdentifier, cty.NullVal(cty.String))
			body.SetAttributeValue(serviceInstanceOfferingNameIdentifier, cty.NullVal(cty.String))
			deps := generictools.NewDependencyAddresses()
			variables := make(generictools.VariableContent)
			ProcessResources(f, tfutils.SubaccountLevel, &variables, &deps, client, generictools.LevelIds{SubaccountId: "subaccount-id"})
			assert.Equal(t, 2, calls)
			assert.Equal(t, "my-instance", serviceInstanceLiteral(body, "name"))
			assert.Equal(t, "subaccount-id", serviceInstanceLiteral(body, "subaccount_id"))
			if failOffering {
				assert.Equal(t, "plan-id", serviceInstanceLiteral(body, serviceInstancePlanIdentifier))
				assert.Nil(t, body.GetAttribute(serviceInstancePlanNameIdentifier))
				assert.Nil(t, body.GetAttribute(serviceInstanceOfferingNameIdentifier))
			} else {
				assert.Nil(t, body.GetAttribute(serviceInstancePlanIdentifier))
				assert.Equal(t, "standard", serviceInstanceLiteral(body, serviceInstancePlanNameIdentifier))
				assert.Equal(t, "alert-notification", serviceInstanceLiteral(body, serviceInstanceOfferingNameIdentifier))
			}
		})
	}
}
