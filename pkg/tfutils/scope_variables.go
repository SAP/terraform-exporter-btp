package tfutils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// applyExecutionTargetVariables makes the exported configuration consume the
// selected target from terraform.tfvars. Existing resource references are left
// untouched so that cleanup-created dependencies remain in place.
func applyExecutionTargetVariables(directory string, target ExecutionTarget) error {
	terraformFiles, err := readTerraformFiles(directory)
	if err != nil {
		return err
	}

	for _, terraformFile := range terraformFiles {
		changed := rewriteScopeVariables(terraformFile.file.Body(), target, "")
		if target.OrgID != "" && rewriteLegacyOrgVariable(terraformFile.file.Body()) {
			changed = true
		}
		if changed {
			if err := os.WriteFile(terraformFile.path, terraformFile.file.Bytes(), terraformFile.mode.Perm()); err != nil {
				return fmt.Errorf("write Terraform configuration %q: %w", terraformFile.path, err)
			}
		}
	}

	return ensureExecutionTargetVariable(directory, target)
}

func rewriteScopeVariables(body *hclwrite.Body, target ExecutionTarget, parentBlockType string) bool {
	changed := false
	attributeName, variableName := scopeAttributeAndVariable(target)
	for name, attribute := range body.Attributes() {
		tokens := attribute.Expr().BuildTokens(nil)
		attributeChanged := false
		if parentBlockType == "resource" && name == attributeName && isExactString(tokens, target.ID()) {
			tokens = hclwrite.Tokens{{Type: hclsyntax.TokenIdent, Bytes: []byte("var." + variableName)}}
			attributeChanged = true
		}
		if parentBlockType == "import" {
			if name == "id" && isExactString(tokens, target.ID()) {
				tokens = hclwrite.Tokens{{Type: hclsyntax.TokenIdent, Bytes: []byte("var." + variableName)}}
				attributeChanged = true
			}
			if name == "identity" && replaceIdentityScopeValue(&tokens, attributeName, target.ID(), variableName) {
				attributeChanged = true
			}
		}
		if attributeChanged {
			body.SetAttributeRaw(name, tokens)
			changed = true
		}
	}

	for _, block := range body.Blocks() {
		if rewriteScopeVariables(block.Body(), target, block.Type()) {
			changed = true
		}
	}
	return changed
}

func scopeAttributeAndVariable(target ExecutionTarget) (string, string) {
	switch {
	case target.SubaccountID != "":
		return "subaccount_id", "subaccount_id"
	case target.DirectoryID != "":
		return "directory_id", "directory_id"
	default:
		return "org", "org_id"
	}
}

func isExactString(tokens hclwrite.Tokens, value string) bool {
	return len(tokens) == 3 && tokens[0].Type == hclsyntax.TokenOQuote && tokens[1].Type == hclsyntax.TokenQuotedLit && tokens[2].Type == hclsyntax.TokenCQuote && string(tokens[1].Bytes) == value
}

func replaceIdentityScopeValue(tokens *hclwrite.Tokens, scopeAttribute, scopeID, variableName string) bool {
	for i := 0; i+4 < len(*tokens); i++ {
		if (*tokens)[i].Type != hclsyntax.TokenIdent || string((*tokens)[i].Bytes) != scopeAttribute || (*tokens)[i+1].Type != hclsyntax.TokenEqual {
			continue
		}
		if !isExactString((*tokens)[i+2:i+5], scopeID) {
			continue
		}
		replacement := hclwrite.Tokens{{Type: hclsyntax.TokenIdent, Bytes: []byte("var." + variableName)}}
		*tokens = append((*tokens)[:i+2], append(replacement, (*tokens)[i+5:]...)...)
		return true
	}
	return false
}

func rewriteLegacyOrgVariable(body *hclwrite.Body) bool {
	changed := false
	for name, attribute := range body.Attributes() {
		tokens := attribute.Expr().BuildTokens(nil)
		attributeChanged := false
		for i := 0; i+2 < len(tokens); i++ {
			if tokens[i].Type == hclsyntax.TokenIdent && string(tokens[i].Bytes) == "var" && tokens[i+1].Type == hclsyntax.TokenDot && tokens[i+2].Type == hclsyntax.TokenIdent && string(tokens[i+2].Bytes) == "org" {
				tokens[i+2].Bytes = []byte("org_id")
				attributeChanged = true
			}
		}
		if attributeChanged {
			body.SetAttributeRaw(name, tokens)
			changed = true
		}
	}
	for _, block := range body.Blocks() {
		if rewriteLegacyOrgVariable(block.Body()) {
			changed = true
		}
	}
	return changed
}

func ensureExecutionTargetVariable(directory string, target ExecutionTarget) error {
	path := filepath.Join(directory, "variables.tf")
	f := hclwrite.NewEmptyFile()
	if source, err := os.ReadFile(path); err == nil {
		parsed, diagnostics := hclwrite.ParseConfig(source, path, hcl.InitialPos)
		if diagnostics.HasErrors() {
			return fmt.Errorf("parse Terraform configuration %q: %s", path, diagnostics.Error())
		}
		f = parsed
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("read Terraform configuration %q: %w", path, err)
	}

	body := f.Body()
	variableName := target.VariableName()
	var variableBlock *hclwrite.Block
	for _, block := range body.Blocks() {
		labels := block.Labels()
		if block.Type() == "variable" && len(labels) == 1 {
			if target.OrgID != "" && labels[0] == "org" {
				body.RemoveBlock(block)
				continue
			}
			if labels[0] == variableName {
				variableBlock = block
			}
		}
	}
	if variableBlock == nil {
		variableBlock = body.AppendNewBlock("variable", []string{variableName})
		body.AppendNewline()
	}
	variableBody := variableBlock.Body()
	variableBody.SetAttributeRaw("type", hclwrite.Tokens{{Type: hclsyntax.TokenIdent, Bytes: []byte("string")}})
	variableBody.SetAttributeValue("description", cty.StringVal(scopeVariableDescription(target)))
	variableBody.RemoveAttribute("default")

	return os.WriteFile(path, hclwrite.Format(f.Bytes()), 0600)
}

func scopeVariableDescription(target ExecutionTarget) string {
	switch {
	case target.SubaccountID != "":
		return "ID of the SAP BTP subaccount being imported"
	case target.DirectoryID != "":
		return "ID of the SAP BTP directory being imported"
	default:
		return "ID of the Cloud Foundry organization being imported"
	}
}
