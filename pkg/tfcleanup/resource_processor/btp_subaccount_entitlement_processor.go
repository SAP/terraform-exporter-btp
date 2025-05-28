package resourceprocessor

import (
	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

const subaccountEntitlementBlockIdentifier = "btp_subaccount_entitlement"
const entitlementPlanNameIdentifier = "plan_name"
const entitlementServiceNameIdentifier = "service_name"

func fillSubaccountEntitlementDependencyAddresses(body *hclwrite.Body, resourceAddress string, dependencyAddresses *generictools.DependencyAddresses) {
	planNameAttr := body.GetAttribute(entitlementPlanNameIdentifier)
	serviceNameAttr := body.GetAttribute(entitlementServiceNameIdentifier)

	if planNameAttr == nil || serviceNameAttr == nil {
		return
	}

	planNameTokens := planNameAttr.Expr().BuildTokens(nil)
	serviceNameTokens := serviceNameAttr.Expr().BuildTokens(nil)

	var planName string
	var serviceName string

	if len(planNameTokens) == 3 {
		planName = generictools.GetStringToken(planNameTokens)
	}

	if len(serviceNameTokens) == 3 {
		serviceName = generictools.GetStringToken(serviceNameTokens)
	}

	if planName != "" && serviceName != "" {
		key := generictools.EntitlementKey{
			ServiceName: serviceName,
			PlanName:    planName,
		}

		(*dependencyAddresses).EntitlementAddress[key] = resourceAddress
	}
}
