package resourceprocessor

import (
	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

const subaccountEntitlementBlockIdentifier = "btp_subaccount_entitlement"
const entitlementPlanNameIdentifier = "plan_name"
const entitlementServiceNameIdentifier = "service_name"

func fillSubaccountEntitlementDependencyAddresses(body *hclwrite.Body, resourceAddress string, dependencyAddresses *generictools.DepedendcyAddresses) {
	attrs := body.Attributes()

	var planName string
	var serviceName string

	for name, attr := range attrs {
		tokens := attr.Expr().BuildTokens(nil)
		if name == entitlementPlanNameIdentifier && len(tokens) == 3 {
			planName = generictools.GetStringToken(tokens)
		}

		if name == entitlementServiceNameIdentifier && len(tokens) == 3 {
			serviceName = generictools.GetStringToken(tokens)
		}
	}

	if planName != "" && serviceName != "" {
		key := generictools.EntitlementKey{
			ServiceName: serviceName,
			PlanName:    planName,
		}

		(*dependencyAddresses).EntitlementAddress[key] = resourceAddress

	}
}
