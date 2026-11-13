package tfutils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveExecutionTarget(t *testing.T) {
	t.Run("uses terraform.tfvars when no command-line scope is provided", func(t *testing.T) {
		withTargetWorkingDirectory(t, "subaccount_id = \"subaccount-from-file\"\n")
		target, err := ResolveExecutionTarget("", "", "")
		if err != nil {
			t.Fatal(err)
		}
		if target.SubaccountID != "subaccount-from-file" || target.VariableName() != "subaccount_id" {
			t.Fatalf("unexpected target: %#v", target)
		}
	})

	t.Run("command-line value overrides the same file scope", func(t *testing.T) {
		withTargetWorkingDirectory(t, "directory_id = \"directory-from-file\"\n")
		target, err := ResolveExecutionTarget("", "directory-from-flag", "")
		if err != nil {
			t.Fatal(err)
		}
		if target.DirectoryID != "directory-from-flag" {
			t.Fatalf("directory ID = %q", target.DirectoryID)
		}
	})

	t.Run("different scopes conflict", func(t *testing.T) {
		withTargetWorkingDirectory(t, "org_id = \"org-from-file\"\n")
		_, err := ResolveExecutionTarget("subaccount-from-flag", "", "")
		if err == nil || !strings.Contains(err.Error(), "conflicts") {
			t.Fatalf("error = %v, want conflict", err)
		}
	})

	t.Run("multiple file scopes are rejected", func(t *testing.T) {
		withTargetWorkingDirectory(t, "subaccount_id = \"subaccount\"\norg_id = \"org\"\n")
		_, err := ResolveExecutionTarget("", "", "")
		if err == nil || !strings.Contains(err.Error(), "exactly one") {
			t.Fatalf("error = %v, want multiple-scope error", err)
		}
	})
}

func TestWriteExecutionTargetTFVars(t *testing.T) {
	directory := t.TempDir()
	target := ExecutionTarget{OrgID: "org-id", TFVars: []byte("org_id = \"org-id\"\nregion = \"eu10\"\n")}
	if err := WriteExecutionTargetTFVars(directory, target); err != nil {
		t.Fatal(err)
	}
	actual, err := os.ReadFile(filepath.Join(directory, terraformVarsFileName))
	if err != nil {
		t.Fatal(err)
	}
	if string(actual) != string(target.TFVars) {
		t.Fatalf("terraform.tfvars = %q, want %q", actual, target.TFVars)
	}
}

func TestWriteExecutionTargetTFVarsAppliesCommandLinePrecedence(t *testing.T) {
	directory := t.TempDir()
	target := ExecutionTarget{
		DirectoryID: "directory-from-flag",
		TFVars:      []byte("directory_id = \"directory-from-file\"\nregion = \"eu10\"\n"),
		fromFlag:    true,
	}
	if err := WriteExecutionTargetTFVars(directory, target); err != nil {
		t.Fatal(err)
	}
	actual, err := os.ReadFile(filepath.Join(directory, terraformVarsFileName))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(actual), "directory_id = \"directory-from-flag\"") || !strings.Contains(string(actual), "eu10") {
		t.Fatalf("unexpected terraform.tfvars:\n%s", actual)
	}
}

func TestWriteExecutionTargetTFVarsPreservesExistingFile(t *testing.T) {
	directory := t.TempDir()
	path := filepath.Join(directory, terraformVarsFileName)
	existing := "directory_id = \"directory-id\"\nregion = \"existing-region\"\n"
	if err := os.WriteFile(path, []byte(existing), 0600); err != nil {
		t.Fatal(err)
	}
	target := ExecutionTarget{DirectoryID: "directory-id", TFVars: []byte("directory_id = \"directory-id\"\nregion = \"source-region\"\n")}
	if err := WriteExecutionTargetTFVars(directory, target); err != nil {
		t.Fatal(err)
	}
	actual, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(actual) != existing {
		t.Fatalf("terraform.tfvars was overwritten:\n%s", actual)
	}
}

func TestWriteExecutionTargetTFVarsRejectsExistingDifferentTarget(t *testing.T) {
	directory := t.TempDir()
	path := filepath.Join(directory, terraformVarsFileName)
	if err := os.WriteFile(path, []byte("org_id = \"other-org\"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	err := WriteExecutionTargetTFVars(directory, ExecutionTarget{OrgID: "requested-org"})
	if err == nil || !strings.Contains(err.Error(), "refusing to overwrite") {
		t.Fatalf("error = %v, want overwrite-protection error", err)
	}
}

func withTargetWorkingDirectory(t *testing.T, tfvars string) {
	t.Helper()
	directory := t.TempDir()
	if err := os.WriteFile(filepath.Join(directory, terraformVarsFileName), []byte(tfvars), 0600); err != nil {
		t.Fatal(err)
	}
	original, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(directory); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(original) })
}
