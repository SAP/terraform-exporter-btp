package resourceprocessor

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// Cloud Foundry exports include runtime status in the provisioning parameters.
func processEnvironmentInstanceParameters(body *hclwrite.Body) {
	environmentType, ok := literalStringAttribute(body, "environment_type")
	if !ok || environmentType != "cloudfoundry" {
		return
	}
	parameters := body.GetAttribute("parameters")
	if parameters == nil {
		return
	}
	tokens := parameters.Expr().BuildTokens(nil)
	expression, diagnostics := hclsyntax.ParseExpression(tokens.Bytes(), "", hcl.InitialPos)
	if diagnostics.HasErrors() {
		return
	}
	call, ok := expression.(*hclsyntax.FunctionCallExpr)
	if !ok || call.Name != "jsonencode" || len(call.Args) != 1 {
		return
	}
	object, ok := call.Args[0].(*hclsyntax.ObjectConsExpr)
	if !ok {
		return
	}

	attributes := make([]hclwrite.ObjectAttrTokens, 0, len(object.Items))
	removed := false
	for _, item := range object.Items {
		key, diagnostics := item.KeyExpr.Value(nil)
		if !diagnostics.HasErrors() && key.IsKnown() && !key.IsNull() && key.Type() == cty.String && key.AsString() == "status" {
			removed = true
			continue
		}
		// Keep expressions intact: other parameters may contain references or
		// nested objects with their own status fields.
		attributes = append(attributes, hclwrite.ObjectAttrTokens{
			Name:  environmentParameterTokens(tokens, item.KeyExpr.Range()),
			Value: environmentParameterTokens(tokens, item.ValueExpr.Range()),
		})
	}
	if !removed {
		return
	}
	if len(attributes) == 0 {
		body.RemoveAttribute("parameters")
		return
	}
	body.SetAttributeRaw("parameters", hclwrite.TokensForFunctionCall("jsonencode", hclwrite.TokensForObject(attributes)))
}

// Expression ranges refer to the bytes emitted by tokens.Bytes(), including
// each token's leading spaces. Select the original tokens to preserve values.
func environmentParameterTokens(tokens hclwrite.Tokens, expressionRange hcl.Range) hclwrite.Tokens {
	var selected hclwrite.Tokens
	offset := 0
	for _, token := range tokens {
		offset += token.SpacesBefore
		end := offset + len(token.Bytes)
		if offset >= expressionRange.Start.Byte && end <= expressionRange.End.Byte {
			selected = append(selected, token)
		}
		offset = end
	}
	return selected
}
