package resourceprocessor

import (
	"github.com/SAP/terraform-exporter-btp/pkg/tfutils"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// Resolve collection references before scope literals are replaced. Collect all
// addresses first because a curated inventory can export assignments before
// their collections, or omit the collections entirely.
func linkRoleCollectionAssignments(body *hclwrite.Body) {
	type collectionKey struct {
		subaccount string
		name       string
	}
	collections := make(map[collectionKey]string)
	for _, block := range body.Blocks() {
		labels := block.Labels()
		if block.Type() != "resource" || len(labels) != 2 || labels[0] != subaccountRoleCollectionBlockIdentifier {
			continue
		}
		subaccount, hasSubaccount := literalStringAttribute(block.Body(), subaccountIdentifier)
		name, hasName := literalStringAttribute(block.Body(), "name")
		if hasSubaccount && hasName {
			collections[collectionKey{subaccount, name}] = labels[1]
		}
	}
	for _, block := range body.Blocks() {
		labels := block.Labels()
		if block.Type() != "resource" || len(labels) != 2 || labels[0] != tfutils.SubaccountRoleCollectionAssignmentType {
			continue
		}
		subaccount, hasSubaccount := literalStringAttribute(block.Body(), subaccountIdentifier)
		name, hasName := literalStringAttribute(block.Body(), "role_collection_name")
		if !hasSubaccount || !hasName {
			continue
		}
		if label, exists := collections[collectionKey{subaccount, name}]; exists {
			block.Body().SetAttributeTraversal("role_collection_name", hcl.Traversal{
				hcl.TraverseRoot{Name: subaccountRoleCollectionBlockIdentifier},
				hcl.TraverseAttr{Name: label},
				hcl.TraverseAttr{Name: "name"},
			})
		}
	}
}

func literalStringAttribute(body *hclwrite.Body, name string) (string, bool) {
	attribute := body.GetAttribute(name)
	if attribute == nil {
		return "", false
	}
	expression, diagnostics := hclsyntax.ParseExpression(attribute.Expr().BuildTokens(nil).Bytes(), "", hcl.InitialPos)
	if diagnostics.HasErrors() {
		return "", false
	}
	value, diagnostics := expression.Value(nil)
	if diagnostics.HasErrors() || !value.IsKnown() || value.IsNull() || value.Type() != cty.String {
		return "", false
	}
	return value.AsString(), true
}
