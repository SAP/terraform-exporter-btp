package resourceprocessor

import (
	"strconv"

	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

const subaccountEntitlementBlockIdentifier = "btp_subaccount_entitlement"
const entitlementPlanNameIdentifier = "plan_name"
const entitlementServiceNameIdentifier = "service_name"
const entitlementAmountIdentifier = "amount"
const moduleName = "btp_subaccount_entitlement"
const genericEntitlementModuleAddress = "module." + moduleName
const entitlementModuleName = "aydin-ozcan/sap-btp-entitlements/btp"
const entitlementModuleVersion = "1.0.1"

func fillSubaccountEntitlementDependencyAddresses(body *hclwrite.Body, resourceAddress string, dependencyAddresses *generictools.DependencyAddresses) {
	planNameAttr := body.GetAttribute(entitlementPlanNameIdentifier)
	serviceNameAttr := body.GetAttribute(entitlementServiceNameIdentifier)
	amountAttr := body.GetAttribute(entitlementAmountIdentifier)

	if planNameAttr == nil || serviceNameAttr == nil {
		return
	}

	planNameTokens := planNameAttr.Expr().BuildTokens(nil)
	serviceNameTokens := serviceNameAttr.Expr().BuildTokens(nil)

	var planName string
	var serviceName string
	var amount int

	if len(planNameTokens) == 3 {
		planName = generictools.GetStringToken(planNameTokens)
	}

	if len(serviceNameTokens) == 3 {
		serviceName = generictools.GetStringToken(serviceNameTokens)
	}

	if amountAttr != nil {
		amountTokens := amountAttr.Expr().BuildTokens(nil)
		if len(amountTokens) == 1 {
			amount, _ = strconv.Atoi(string(amountTokens[0].Bytes))
		} else {
			amount = 0
		}
	}

	if planName != "" && serviceName != "" {
		key := generictools.EntitlementKey{
			ServiceName: serviceName,
			PlanName:    planName,
			Amount:      amount,
		}

		(*dependencyAddresses).EntitlementAddress[key] = resourceAddress
	}
}

func addEntitlementModule(body *hclwrite.Body, subaccountAddress string) {
	body.AppendNewline()

	moduleBlock := body.AppendNewBlock("module", []string{moduleName})

	moduleBlock.Body().SetAttributeRaw("source", hclwrite.Tokens{
		{
			Type:  hclsyntax.TokenOQuote,
			Bytes: []byte(`"`),
		},
		{
			Type:  hclsyntax.TokenStringLit,
			Bytes: []byte(entitlementModuleName),
		},
		{
			Type:  hclsyntax.TokenOQuote,
			Bytes: []byte(`"`),
		},
	})

	moduleBlock.Body().SetAttributeRaw("version", hclwrite.Tokens{
		{
			Type:  hclsyntax.TokenOQuote,
			Bytes: []byte(`"`),
		},
		{
			Type:  hclsyntax.TokenStringLit,
			Bytes: []byte(entitlementModuleVersion),
		},
		{
			Type:  hclsyntax.TokenOQuote,
			Bytes: []byte(`"`),
		},
	})

	moduleBlock.Body().SetAttributeRaw("subaccount_id", hclwrite.Tokens{
		{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte(subaccountAddress),
		},
	})

	moduleBlock.Body().SetAttributeRaw("entitlements", hclwrite.Tokens{
		{
			Type:  hclsyntax.TokenStringLit,
			Bytes: []byte("var.btp_subaccount_entitlements"),
		},
	})
}

func addEntitlementVariables(variablesToCreate *generictools.VariableContent, dependencyAddresses *generictools.DependencyAddresses) {
	if len(dependencyAddresses.EntitlementAddress) == 0 {
		return
	}

	variableName := "btp_subaccount_entitlements"
	variableInfo := "Object of entitlements to be assigned to the subaccount."
	variableType := "map(list(string))"

	defaultValue := "{\n"
	for key := range dependencyAddresses.EntitlementAddress {
		serviceName := key.ServiceName
		planName := key.PlanName
		amount := key.Amount

		if amount > 0 {
			defaultValue += " \"" + serviceName + "\" = [\"" + planName + "=" + strconv.Itoa(amount) + "\"]\n"
		} else {
			defaultValue += " \"" + serviceName + "\" = [\"" + planName + "\"]\n"
		}
	}
	defaultValue += "}"

	(*variablesToCreate)[variableName] = generictools.VariableInfo{
		Description:  variableInfo,
		Type:         variableType,
		DefaultValue: defaultValue,
	}
}
