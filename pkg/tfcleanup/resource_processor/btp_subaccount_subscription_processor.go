package resourceprocessor

import (
	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

const subscriptionBlockIdentifier = "btp_subaccount_subscription"
const subscriptionAppNameIdentifier = "app_name"
const subscriptionPlanNameIdentifier = "plan_name"

func addEntitlementDependency(body *hclwrite.Body, dependencyAddresses *generictools.DepedendcyAddresses) {
	var appName string
	var planName string

	subscriptionAppNameAttr := body.GetAttribute(subscriptionAppNameIdentifier)
	subscriptionPlanNameAttr := body.GetAttribute(subscriptionPlanNameIdentifier)

	if subscriptionAppNameAttr == nil || subscriptionPlanNameAttr == nil {
		return
	}

	subscriptionAppNameAttrTokens := subscriptionAppNameAttr.Expr().BuildTokens(nil)
	subscriptionPlanNameAttrTokens := subscriptionPlanNameAttr.Expr().BuildTokens(nil)

	if len(subscriptionAppNameAttrTokens) == 3 {
		appName = generictools.GetStringToken(subscriptionAppNameAttrTokens)
	}

	if len(subscriptionPlanNameAttrTokens) == 3 {
		planName = generictools.GetStringToken(subscriptionPlanNameAttrTokens)
	}

	if appName != "" && planName != "" {
		key := generictools.EntitlementKey{
			ServiceName: appName,
			PlanName:    planName,
		}

		dependencyAddress := (*dependencyAddresses).EntitlementAddress[key]

		if dependencyAddress != "" {
			body.SetAttributeRaw("depends_on", hclwrite.Tokens{
				{
					Type:  hclsyntax.TokenOBrack,
					Bytes: []byte("["),
				},
				{Type: hclsyntax.TokenStringLit,
					Bytes: []byte(dependencyAddress),
				},
				{
					Type:  hclsyntax.TokenCBrack,
					Bytes: []byte("]"),
				},
			})
		}
	}
}
