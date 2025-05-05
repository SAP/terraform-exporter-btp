package resourceprocessor

import (
	"github.com/SAP/terraform-exporter-btp/internal/btpcli"
	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

const serviceInstanceBlockIdentifier = "btp_subaccount_service_instance"
const serviceInstancePlanIdentifier = "serviceplan_id"

func addServiceInstanceDependency(body *hclwrite.Body, dependencyAddresses *generictools.DependencyAddresses, btpClient *btpcli.ClientFacade, subaccountId string) {

	// 1: Iterate over attributes to find the value for the service name and the plan ID
	// 2: Fetch the plan name via ID using the BTP CLI client
	// 3: Check if there is a dependency address for the plan name in the entitlements
	// 4: If we find such a dependenceny: Create a data source for the service plan that depends on the entitlement
	// 5: Exchange the explicit plan ID with the data source reference
	planIdAttr := body.GetAttribute(serviceInstancePlanIdentifier)
	if planIdAttr == nil {
		// No plan ID found, no further action will be taken
		return
	}

	var planId string

	planIdAttrTokens := planIdAttr.Expr().BuildTokens(nil)

	if len(planIdAttrTokens) == 3 {
		planId = generictools.GetStringToken(planIdAttrTokens)
	}

	if planId == "" {
		// Nothing found, no further action will be taken
		return
	}

	planName, serviceName, err := btpcli.GetServiceDataByPlanId(btpClient, subaccountId, planId)
	if err != nil {
		// No plan name found, no refinement of the code will be done
		return
	}

	key := generictools.EntitlementKey{
		ServiceName: serviceName,
		PlanName:    planName,
	}

	dependencyAddress := (*dependencyAddresses).EntitlementAddress[key]
	if dependencyAddress == "" {
		//No entitlement exported that fits the service instance
		return
	}

	datasourceAddress := serviceName + "_" + planName

	// replce the plan ID with the data source reference
	body.SetAttributeRaw(serviceInstancePlanIdentifier, hclwrite.Tokens{
		{
			Type:  hclsyntax.TokenStringLit,
			Bytes: []byte("data.btp_subaccount_service_plan." + datasourceAddress + ".id"),
		},
	})

	(*dependencyAddresses).DataSourceInfo = append((*dependencyAddresses).DataSourceInfo, generictools.DataSourceInfo{
		DatasourceAddress:  datasourceAddress,
		SubaccountAddress:  (*dependencyAddresses).SubaccountAddress + ".id",
		OfferingName:       serviceName,
		Name:               planName,
		EntitlementAddress: dependencyAddress,
	},
	)
}

func addServicePlanDataSources(body *hclwrite.Body, datasourceInfo generictools.DataSourceInfo) {
	body.AppendNewline()

	dsBlock := body.AppendNewBlock("data", []string{"btp_subaccount_service_plan", datasourceInfo.DatasourceAddress})

	dsBlock.Body().SetAttributeRaw("subaccount_id", hclwrite.Tokens{
		{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte(datasourceInfo.SubaccountAddress),
		},
	})

	dsBlock.Body().SetAttributeRaw("offering_name", hclwrite.Tokens{
		{
			Type:  hclsyntax.TokenOQuote,
			Bytes: []byte((`"`)),
		},
		{
			Type:  hclsyntax.TokenStringLit,
			Bytes: []byte(datasourceInfo.OfferingName),
		},
		{
			Type:  hclsyntax.TokenOQuote,
			Bytes: []byte((`"`)),
		},
	})

	dsBlock.Body().SetAttributeRaw("name", hclwrite.Tokens{
		{
			Type:  hclsyntax.TokenOQuote,
			Bytes: []byte(`"`),
		},
		{
			Type:  hclsyntax.TokenStringLit,
			Bytes: []byte(datasourceInfo.Name),
		},
		{
			Type:  hclsyntax.TokenOQuote,
			Bytes: []byte(`"`),
		},
	})

	dsBlock.Body().SetAttributeRaw("depends_on", hclwrite.Tokens{
		{
			Type:  hclsyntax.TokenOBrack,
			Bytes: []byte("["),
		},
		{Type: hclsyntax.TokenStringLit,
			Bytes: []byte(datasourceInfo.EntitlementAddress),
		},
		{
			Type:  hclsyntax.TokenCBrack,
			Bytes: []byte("]"),
		},
	})
}
