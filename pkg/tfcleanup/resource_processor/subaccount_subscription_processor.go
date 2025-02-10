package resourceprocessor

import (
	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func addEntitlementDependency(body *hclwrite.Body, dependencyAddresses *generictools.DepedendcyAddresses) {
	attrs := body.Attributes()

	var appName string
	var planName string

	for name, attr := range attrs {
		tokens := attr.Expr().BuildTokens(nil)
		if name == subscriptionAppNameIdentifier && len(tokens) == 3 {
			appName = generictools.GetStringToken(tokens)
		}

		if name == subscriptionPlanNameIdentifier && len(tokens) == 3 {
			planName = generictools.GetStringToken(tokens)
		}
	}

	if appName != "" && planName != "" {
		key := generictools.EntilementKey{
			ServiceName: appName,
			PlanName:    planName,
		}

		dependencyAddress := (*dependencyAddresses).EntitlementAddress[key]

		if dependencyAddress != "" {
			// Add depends_on to the subscription

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
			},
			)

		}
	}
}
