package resourceprocessor

import (
	"log"

	"github.com/SAP/terraform-exporter-btp/internal/btpcli"
	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/SAP/terraform-exporter-btp/pkg/toggles"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

const subscriptionBlockIdentifier = "btp_subaccount_subscription"
const subscriptionAppNameIdentifier = "app_name"
const subscriptionPlanNameIdentifier = "plan_name"

func addEntitlementDependency(body *hclwrite.Body, dependencyAddresses *generictools.DependencyAddresses, btpClient *btpcli.ClientFacade, subaccountId string) {
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

		if dependencyAddress == "" {
			//Check if the app name used is the technical app name and switch to the commercial app name
			dependencyAddress = handleCommercialAppName(appName, planName, dependencyAddresses, btpClient, subaccountId)
		}

		if dependencyAddress != "" && !toggles.IsEntitlementModuleGenerationDeactivated() {
			// Add generic module dependency
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
		} else if dependencyAddress != "" && toggles.IsEntitlementModuleGenerationDeactivated() {
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

func handleCommercialAppName(appName string, planName string, dependencyAddresses *generictools.DependencyAddresses, btpClient *btpcli.ClientFacade, subaccountId string) (dependencyAddress string) {
	// Check if technical app name is different to service name/commercial app name
	technicalAppName, commercialAppName, err := btpcli.GetAppNamesBySubaccountAndApp(subaccountId, appName, btpClient)
	if err != nil {
		// Error is not critical, so we log it and continue with the flow
		log.Printf("Error fetching app names from platform for %s: %s", appName, err)
		return ""
	}

	if technicalAppName != commercialAppName {
		// Try to fetch an entry using the commercial app name
		key := generictools.EntitlementKey{
			ServiceName: commercialAppName,
			PlanName:    planName,
		}
		dependencyAddress = (*dependencyAddresses).EntitlementAddress[key]
	}
	return dependencyAddress
}
