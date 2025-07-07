package resourceprocessor

import (
	"path/filepath"
	"strconv"
	"strings"

	generictools "github.com/SAP/terraform-exporter-btp/pkg/tfcleanup/generic_tools"
	"github.com/SAP/terraform-exporter-btp/pkg/toggles"
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
		}

		(*dependencyAddresses).EntitlementAddress[key] = generictools.EntitlementInfo{Amount: amount, Address: resourceAddress}
	}
}

func addEntitlementModule(body *hclwrite.Body, subaccountAddress string, subaccountId string, entitlementAddress map[generictools.EntitlementKey]generictools.EntitlementInfo) {
	if toggles.IsEntitlementModuleGenerationDeactivated() {
		return
	}

	if len(entitlementAddress) == 0 {
		return
	}

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

	if subaccountAddress == "" {
		moduleBlock.Body().SetAttributeRaw("subaccount", hclwrite.Tokens{
			{
				Type:  hclsyntax.TokenOQuote,
				Bytes: []byte(`"`),
			},
			{
				Type:  hclsyntax.TokenStringLit,
				Bytes: []byte(subaccountId),
			},
			{
				Type:  hclsyntax.TokenOQuote,
				Bytes: []byte(`"`),
			},
		})
	} else {
		moduleBlock.Body().SetAttributeRaw("subaccount", hclwrite.Tokens{
			{
				Type:  hclsyntax.TokenIdent,
				Bytes: []byte(subaccountAddress),
			},
		})

	}

	moduleBlock.Body().SetAttributeRaw("entitlements", hclwrite.Tokens{
		{
			Type:  hclsyntax.TokenStringLit,
			Bytes: []byte("var.btp_subaccount_entitlements"),
		},
	})
}

func addEntitlementVariables(variablesToCreate *generictools.VariableContent, entitlementAddress map[generictools.EntitlementKey]generictools.EntitlementInfo) {
	if toggles.IsEntitlementModuleGenerationDeactivated() {
		return
	}

	if len(entitlementAddress) == 0 {
		return
	}

	variableName := "btp_subaccount_entitlements"
	variableInfo := "Object of entitlements to be assigned to the subaccount."
	variableType := "map(list(string))"

	defaultValue := "{\n"
	for key, info := range entitlementAddress {
		serviceName := key.ServiceName
		planName := key.PlanName
		amount := info.Amount

		if amount > 0 {
			stringForValue := serviceName + "=[\"" + planName + "=" + strconv.Itoa(amount) + "\"]\n"
			defaultValue += strings.ReplaceAll(stringForValue, " ", "")
		} else {
			defaultValue += "\"" + serviceName + "\" = [\"" + planName + "\"]\n"
		}
	}
	defaultValue += "}"

	(*variablesToCreate)[variableName] = generictools.VariableInfo{
		Description:  variableInfo,
		Type:         variableType,
		DefaultValue: defaultValue,
	}
}

func appendEntitlementBlocksToRemove(dependencyAddresses *generictools.DependencyAddresses) {

	if toggles.IsEntitlementModuleGenerationDeactivated() {
		return
	}

	if len(dependencyAddresses.EntitlementAddress) == 0 {
		return
	}

	const blockIdentifier = "btp_subaccount_entitlement"
	for _, entitlementDependencyInfo := range dependencyAddresses.EntitlementAddress {
		entitlementImportToRemove := generictools.BlockSpecifier{
			BlockIdentifier: blockIdentifier,
			ResourceAddress: entitlementDependencyInfo.Address,
		}
		dependencyAddresses.BlocksToRemove = append(dependencyAddresses.BlocksToRemove, entitlementImportToRemove)
	}
}

func AppendImportBlocksForEntitlementModule(directory string, entitlementsToAdd map[generictools.EntitlementKey]generictools.EntitlementInfo, levelIds generictools.LevelIds) {
	if toggles.IsEntitlementModuleGenerationDeactivated() {
		return
	}

	if len(entitlementsToAdd) == 0 {
		return
	}

	filePath := filepath.Join(directory, "btp_subaccount_entitlement_import.tf")
	f := generictools.GetHclFile(filePath)
	body := f.Body()

	for key := range entitlementsToAdd {

		importToKey := key.ServiceName + "-" + key.PlanName
		importIdKey := levelIds.SubaccountId + "," + key.ServiceName + "," + key.PlanName

		body.AppendNewline()
		importBlock := body.AppendNewBlock("import", []string{})
		importBlock.Body().SetAttributeRaw("to", hclwrite.Tokens{
			{
				Type:  hclsyntax.TokenStringLit,
				Bytes: []byte(genericEntitlementModuleAddress + "." + subaccountEntitlementBlockIdentifier + ".entitlement[\"" + importToKey + "\"]"),
			},
		})

		importBlock.Body().SetAttributeRaw("id", hclwrite.Tokens{
			{
				Type:  hclsyntax.TokenOQuote,
				Bytes: []byte(`"`),
			},
			{
				Type:  hclsyntax.TokenStringLit,
				Bytes: []byte(importIdKey),
			},
			{
				Type:  hclsyntax.TokenOQuote,
				Bytes: []byte(`"`),
			},
		})
	}
	generictools.ProcessChanges(f, filePath)
}

func removeEntitlementConfigBlock(body *hclwrite.Body, entitlementAddress map[generictools.EntitlementKey]generictools.EntitlementInfo) {
	if toggles.IsEntitlementModuleGenerationDeactivated() {
		return
	}

	if len(entitlementAddress) == 0 {
		return
	}

	for _, entitlementInfo := range entitlementAddress {
		generictools.RemoveConfigBlock(body, entitlementInfo.Address)
	}
}

func handleGenericEntitlementModule(body *hclwrite.Body, subaccountId string, dependencyAddresses *generictools.DependencyAddresses, variables *generictools.VariableContent) {

	// Add module to main configuration
	addEntitlementModule(body, dependencyAddresses.SubaccountAddress, subaccountId, dependencyAddresses.EntitlementAddress)

	// Add variables to variables to be generated
	addEntitlementVariables(variables, dependencyAddresses.EntitlementAddress)

	// Remove the entitlement configurations
	removeEntitlementConfigBlock(body, dependencyAddresses.EntitlementAddress)

	// Add import blocks for entitlements to be removed
	appendEntitlementBlocksToRemove(dependencyAddresses)
}
