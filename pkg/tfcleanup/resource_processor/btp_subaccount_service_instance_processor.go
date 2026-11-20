package resourceprocessor

import (
	"log"

	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

const serviceInstanceBlockIdentifier = "btp_subaccount_service_instance"
const serviceInstancePlanIdentifier = "serviceplan_id"
const serviceInstancePlanNameIdentifier = "serviceplan_name"
const serviceInstanceOfferingNameIdentifier = "service_offering_name"

// servicePlanLookup allows cleanup tests to supply explicit lookup results.
type servicePlanLookup func(subaccountID, planID string) (planName, offeringName string, err error)

func processServiceInstance(body *hclwrite.Body, dependencies *generictools.DependencyAddresses, subaccountID string, lookup servicePlanLookup, generateEntitlementModule bool) {
	planName := serviceInstanceLiteral(body, serviceInstancePlanNameIdentifier)
	offeringName := serviceInstanceLiteral(body, serviceInstanceOfferingNameIdentifier)
	if planName == "" || offeringName == "" {
		planID := serviceInstanceLiteral(body, serviceInstancePlanIdentifier)
		if planID == "" {
			return
		}
		planName, offeringName, err := lookup(subaccountID, planID)
		if err != nil {
			log.Printf("Error resolving service plan %s in subaccount %s; keeping original configuration: %s", planID, subaccountID, err)
			return
		}
		if planName == "" || offeringName == "" {
			log.Printf("Incomplete names for service plan %s in subaccount %s; keeping original configuration", planID, subaccountID)
			return
		}
		body.SetAttributeValue(serviceInstancePlanNameIdentifier, cty.StringVal(planName))
		body.SetAttributeValue(serviceInstanceOfferingNameIdentifier, cty.StringVal(offeringName))
	}
	addServiceInstanceDependency(body, dependencies, generateEntitlementModule)
}

// Only literal strings are eligible for platform lookup or entitlement matching.
// Parsing also decodes HCL escapes without treating expressions as IDs.
func serviceInstanceLiteral(body *hclwrite.Body, name string) string {
	attr := body.GetAttribute(name)
	if attr == nil {
		return ""
	}
	expr, diags := hclsyntax.ParseExpression(attr.Expr().BuildTokens(nil).Bytes(), "", hcl.InitialPos)
	if diags.HasErrors() {
		return ""
	}
	template, ok := expr.(*hclsyntax.TemplateExpr)
	if !ok || !template.IsStringLiteral() {
		return ""
	}
	value, diags := template.Value(nil)
	if diags.HasErrors() || value.IsNull() || !value.IsKnown() || value.Type() != cty.String {
		return ""
	}
	return value.AsString()
}

func addServiceInstanceDependency(body *hclwrite.Body, dependencyAddresses *generictools.DependencyAddresses, generateEntitlementModule bool) {
	serviceOfferingName := serviceInstanceLiteral(body, serviceInstanceOfferingNameIdentifier)
	servicePlanName := serviceInstanceLiteral(body, serviceInstancePlanNameIdentifier)
	if serviceOfferingName == "" || servicePlanName == "" {
		return
	}

	// The name-based configuration is complete, so remove the mutually
	// exclusive service plan ID.
	body.RemoveAttribute(serviceInstancePlanIdentifier)

	// Add dependency to entitlement if necessary
	key := generictools.EntitlementKey{
		ServiceName: serviceOfferingName,
		PlanName:    servicePlanName,
	}

	dependencyInfo := (*dependencyAddresses).EntitlementAddress[key]
	if dependencyInfo.Address == "" {
		//No entitlement exported that fits the service instance
		return
	}

	if generateEntitlementModule {
		body.SetAttributeRaw("depends_on", hclwrite.Tokens{
			{
				Type:  hclsyntax.TokenOBrack,
				Bytes: []byte("["),
			},
			{Type: hclsyntax.TokenStringLit,
				Bytes: []byte(genericEntitlementModuleAddress),
			},
			{
				Type:  hclsyntax.TokenCBrack,
				Bytes: []byte("]"),
			},
		})

	} else {
		body.SetAttributeRaw("depends_on", hclwrite.Tokens{
			{
				Type:  hclsyntax.TokenOBrack,
				Bytes: []byte("["),
			},
			{Type: hclsyntax.TokenStringLit,
				Bytes: []byte(dependencyInfo.Address),
			},
			{
				Type:  hclsyntax.TokenCBrack,
				Bytes: []byte("]"),
			},
		})
	}
}
