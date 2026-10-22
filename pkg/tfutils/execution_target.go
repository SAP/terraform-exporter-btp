package tfutils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

const terraformVarsFileName = "terraform.tfvars"

type ExecutionTarget struct {
	SubaccountID string
	DirectoryID  string
	OrgID        string
	TFVars       []byte
	fromFlag     bool
}

func (t ExecutionTarget) LevelAndID() (string, string) {
	return GetExecutionLevelAndId(t.SubaccountID, t.DirectoryID, t.OrgID, "")
}

func (t ExecutionTarget) VariableName() string {
	switch {
	case t.SubaccountID != "":
		return "subaccount_id"
	case t.DirectoryID != "":
		return "directory_id"
	default:
		return "org_id"
	}
}

func (t ExecutionTarget) ID() string {
	_, id := t.LevelAndID()
	return id
}

// ResolveExecutionTarget resolves one export target from explicit command-line
// values and the terraform.tfvars file in the current working directory.
func ResolveExecutionTarget(subaccountID, directoryID, orgID string) (ExecutionTarget, error) {
	fileValues, tfvars, err := readExecutionTargetTFVars()
	if err != nil {
		return ExecutionTarget{}, err
	}

	if scopeCount(fileValues.SubaccountID, fileValues.DirectoryID, fileValues.OrgID) > 1 {
		return ExecutionTarget{}, fmt.Errorf("%s must define exactly one of subaccount_id, directory_id, or org_id", terraformVarsFileName)
	}
	if scopeCount(subaccountID, directoryID, orgID) > 1 {
		return ExecutionTarget{}, fmt.Errorf("only one of subaccount, directory, or org can be specified")
	}

	cliScope := scopeName(subaccountID, directoryID, orgID)
	fileScope := scopeName(fileValues.SubaccountID, fileValues.DirectoryID, fileValues.OrgID)
	if cliScope != "" && fileScope != "" && cliScope != fileScope {
		return ExecutionTarget{}, fmt.Errorf("command-line %s conflicts with %s target in %s", cliScope, fileScope, terraformVarsFileName)
	}

	result := fileValues
	if subaccountID != "" {
		result = ExecutionTarget{SubaccountID: subaccountID, fromFlag: true}
	}
	if directoryID != "" {
		result = ExecutionTarget{DirectoryID: directoryID, fromFlag: true}
	}
	if orgID != "" {
		result = ExecutionTarget{OrgID: orgID, fromFlag: true}
	}
	if result.ID() == "" {
		return ExecutionTarget{}, fmt.Errorf("specify one of --subaccount, --directory, or --org, or define one target in %s", terraformVarsFileName)
	}
	result.TFVars = tfvars
	return result, nil
}

func readExecutionTargetTFVars() (ExecutionTarget, []byte, error) {
	path := filepath.Join(".", terraformVarsFileName)
	source, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return ExecutionTarget{}, nil, nil
	}
	if err != nil {
		return ExecutionTarget{}, nil, fmt.Errorf("read %s: %w", terraformVarsFileName, err)
	}

	file, diagnostics := hclsyntax.ParseConfig(source, path, hcl.InitialPos)
	if diagnostics.HasErrors() {
		return ExecutionTarget{}, nil, fmt.Errorf("parse %s: %s", terraformVarsFileName, diagnostics.Error())
	}
	values := ExecutionTarget{}
	for name, attribute := range file.Body.(*hclsyntax.Body).Attributes {
		switch name {
		case "subaccount_id", "directory_id", "org_id":
			value, diagnostics := attribute.Expr.Value(nil)
			if diagnostics.HasErrors() || value.Type() != cty.String || value.IsNull() || !value.IsKnown() {
				return ExecutionTarget{}, nil, fmt.Errorf("%s in %s must be a known non-null string", name, terraformVarsFileName)
			}
			stringValue := strings.TrimSpace(value.AsString())
			if stringValue == "" {
				return ExecutionTarget{}, nil, fmt.Errorf("%s in %s must not be empty", name, terraformVarsFileName)
			}
			switch name {
			case "subaccount_id":
				values.SubaccountID = stringValue
			case "directory_id":
				values.DirectoryID = stringValue
			case "org_id":
				values.OrgID = stringValue
			}
		}
	}
	return values, source, nil
}

func WriteExecutionTargetTFVars(configDir string, target ExecutionTarget) error {
	path := filepath.Join(configDir, terraformVarsFileName)
	if existing, err := os.ReadFile(path); err == nil {
		if err := validateExistingExecutionTargetTFVars(path, existing, target); err != nil {
			return err
		}
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("read %s: %w", path, err)
	}

	source := target.TFVars
	if source == nil {
		source = []byte(fmt.Sprintf("%s = %q\n", target.VariableName(), target.ID()))
	} else if target.fromFlag {
		file, diagnostics := hclwrite.ParseConfig(source, terraformVarsFileName, hcl.InitialPos)
		if diagnostics.HasErrors() {
			return fmt.Errorf("parse %s: %s", terraformVarsFileName, diagnostics.Error())
		}
		file.Body().SetAttributeValue(target.VariableName(), cty.StringVal(target.ID()))
		source = hclwrite.Format(file.Bytes())
	}
	return os.WriteFile(path, source, 0600)
}

func validateExistingExecutionTargetTFVars(path string, source []byte, target ExecutionTarget) error {
	file, diagnostics := hclsyntax.ParseConfig(source, path, hcl.InitialPos)
	if diagnostics.HasErrors() {
		return fmt.Errorf("parse existing %s: %s", terraformVarsFileName, diagnostics.Error())
	}
	attribute, exists := file.Body.(*hclsyntax.Body).Attributes[target.VariableName()]
	if !exists {
		return fmt.Errorf("existing %s does not define %s; refusing to overwrite it", terraformVarsFileName, target.VariableName())
	}
	value, diagnostics := attribute.Expr.Value(nil)
	if diagnostics.HasErrors() || value.Type() != cty.String || value.IsNull() || !value.IsKnown() {
		return fmt.Errorf("%s in existing %s must be a known non-null string", target.VariableName(), terraformVarsFileName)
	}
	if value.AsString() != target.ID() {
		return fmt.Errorf("existing %s defines %s for a different target; refusing to overwrite it", terraformVarsFileName, target.VariableName())
	}
	return nil
}

func scopeCount(values ...string) int {
	count := 0
	for _, value := range values {
		if value != "" {
			count++
		}
	}
	return count
}

func scopeName(subaccountID, directoryID, orgID string) string {
	switch {
	case subaccountID != "":
		return "subaccount"
	case directoryID != "":
		return "directory"
	case orgID != "":
		return "org"
	default:
		return ""
	}
}
