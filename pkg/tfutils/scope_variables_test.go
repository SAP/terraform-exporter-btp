package tfutils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestApplyExecutionTargetVariablesForSubaccount(t *testing.T) {
	directory := t.TempDir()
	configuration := "resource \"btp_subaccount_subscription\" \"subscription_0\" {\n" +
		"  subaccount_id = \"subaccount-id\"\n" +
		"}\n\n" +
		"resource \"btp_subaccount_role\" \"role_0\" {\n" +
		"  subaccount_id = btp_subaccount.subaccount_0.id\n" +
		"}\n\n" +
		"import {\n" +
		"  to = btp_subaccount_subscription.subscription_0\n" +
		"  identity = {\n" +
		"    subaccount_id = \"subaccount-id\"\n" +
		"    app_name      = \"application\"\n" +
		"  }\n" +
		"}\n\n" +
		"import {\n" +
		"  to = btp_subaccount.subaccount_0\n" +
		"  id = \"subaccount-id\"\n" +
		"}\n"
	if err := os.WriteFile(filepath.Join(directory, "resources.tf"), []byte(configuration), 0600); err != nil {
		t.Fatal(err)
	}

	if err := applyExecutionTargetVariables(directory, ExecutionTarget{SubaccountID: "subaccount-id"}); err != nil {
		t.Fatal(err)
	}

	resources, err := os.ReadFile(filepath.Join(directory, "resources.tf"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(resources)
	for _, expected := range []string{
		"subaccount_id = var.subaccount_id",
		"subaccount_id = btp_subaccount.subaccount_0.id",
		"id = var.subaccount_id",
	} {
		if !strings.Contains(text, expected) {
			t.Errorf("resources.tf does not contain %q:\n%s", expected, text)
		}
	}
	if strings.Count(text, "subaccount_id = var.subaccount_id") != 2 {
		t.Errorf("scope variable was not applied to both resource and import identity:\n%s", text)
	}
	variables, err := os.ReadFile(filepath.Join(directory, "variables.tf"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(variables), "variable \"subaccount_id\"") || strings.Contains(string(variables), "default") {
		t.Errorf("unexpected variables.tf:\n%s", variables)
	}
}

func TestApplyExecutionTargetVariablesMigratesCloudFoundryOrgVariable(t *testing.T) {
	directory := t.TempDir()
	configuration := "resource \"cloudfoundry_space\" \"space_0\" {\n" +
		"  org = var.org\n" +
		"}\n\n" +
		"import {\n" +
		"  to = cloudfoundry_space.space_0\n" +
		"  identity = {\n" +
		"    org  = \"org-id\"\n" +
		"    name = \"space\"\n" +
		"  }\n" +
		"}\n"
	variables := "variable \"org\" {\n  type = string\n  default = \"org-id\"\n}\n"
	if err := os.WriteFile(filepath.Join(directory, "resources.tf"), []byte(configuration), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(directory, "variables.tf"), []byte(variables), 0600); err != nil {
		t.Fatal(err)
	}

	if err := applyExecutionTargetVariables(directory, ExecutionTarget{OrgID: "org-id"}); err != nil {
		t.Fatal(err)
	}

	resources, err := os.ReadFile(filepath.Join(directory, "resources.tf"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(resources), "var.org_id") != 2 {
		t.Errorf("resources.tf was not migrated:\n%s", resources)
	}
	resultVariables, err := os.ReadFile(filepath.Join(directory, "variables.tf"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(resultVariables), "variable \"org\"") || !strings.Contains(string(resultVariables), "variable \"org_id\"") || strings.Contains(string(resultVariables), "default") {
		t.Errorf("unexpected variables.tf:\n%s", resultVariables)
	}
}

func TestApplyExecutionTargetVariablesForDirectory(t *testing.T) {
	directory := t.TempDir()
	configuration := "resource \"btp_directory_role\" \"role_0\" {\n  directory_id = \"directory-id\"\n}\n"
	if err := os.WriteFile(filepath.Join(directory, "resources.tf"), []byte(configuration), 0600); err != nil {
		t.Fatal(err)
	}

	if err := applyExecutionTargetVariables(directory, ExecutionTarget{DirectoryID: "directory-id"}); err != nil {
		t.Fatal(err)
	}

	resources, err := os.ReadFile(filepath.Join(directory, "resources.tf"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(resources), "directory_id = var.directory_id") {
		t.Errorf("resources.tf was not updated:\n%s", resources)
	}
	variables, err := os.ReadFile(filepath.Join(directory, "variables.tf"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(variables), "variable \"directory_id\"") {
		t.Errorf("variables.tf was not updated:\n%s", variables)
	}
}
