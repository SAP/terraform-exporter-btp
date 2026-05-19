package resourceprocessor

import (
	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/SAP/terraform-exporter-btp/pkg/toggles"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

const serviceInstanceBlockIdentifier = "btp_subaccount_service_instance"
const serviceInstancePlanIdentifier = "serviceplan_id"
const serviceInstancePlanNameIdentifier = "serviceplan_name"
const serviceInstanceOfferingNameIdentifier = "service_offering_name"

func addServiceInstanceDependency(body *hclwrite.Body, dependencyAddresses *generictools.DependencyAddresses, subaccountId string) {

	// 1: Iterate over attributes to find the value for the service name and the plan ID
	// 2: Fetch the plan name via ID using the BTP CLI client
	// 3: Check if there is a dependency address for the plan name in the entitlements
	// 4: If we find such a dependenceny: Create a data source for the service plan that depends on the entitlement
	// 5: Exchange the explicit plan ID with the data source reference
	serviceOfferingNameAttr := body.GetAttribute(serviceInstanceOfferingNameIdentifier)
	serviceInstancePlanNameIdentifierAttr := body.GetAttribute(serviceInstancePlanNameIdentifier)

	var serviceOfferingName string
	var servicePlanName string

	serviceOfferingNameAttrTokens := serviceOfferingNameAttr.Expr().BuildTokens(nil)
	serviceInstancePlanNameAttrTokens := serviceInstancePlanNameIdentifierAttr.Expr().BuildTokens(nil)

	if len(serviceOfferingNameAttrTokens) == 3 {
		serviceOfferingName = generictools.GetStringToken(serviceOfferingNameAttrTokens)
	}

	if len(serviceInstancePlanNameAttrTokens) == 3 {
		servicePlanName = generictools.GetStringToken(serviceInstancePlanNameAttrTokens)
	}

	// Remove the service plan id
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

	if !toggles.IsEntitlementModuleGenerationDeactivated() {
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
