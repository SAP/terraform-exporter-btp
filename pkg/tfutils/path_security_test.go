package tfutils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/SAP/terraform-exporter-btp/pkg/files"
	"github.com/stretchr/testify/require"
)

func linkForTest(t *testing.T, target, name string) {
	t.Helper()
	if err := os.Symlink(target, name); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
}

func stubProvider(t *testing.T) {
	t.Helper()
	previous := providerDirectory
	directory, err := files.NewTempDirectory("provider-test-")
	require.NoError(t, err)
	require.NoError(t, directory.Root.WriteFile("provider.tf", []byte("# provider"), 0600))
	providerDirectory = directory
	t.Cleanup(func() { _ = directory.Cleanup(); providerDirectory = previous })
}

func TestSetupRejectsPathsBeforeWriting(t *testing.T) {
	work := t.TempDir()
	t.Chdir(work)
	outside := t.TempDir()
	// Reproduce the old absolute-input discrepancy: even with the mirrored
	// setup parents present, absolute input must fail before either write.
	if filepath.VolumeName(outside) == "" {
		require.NoError(t, os.MkdirAll(filepath.Dir(filepath.Join(work, outside)), 0700))
	}
	for _, name := range []string{"", ".", "..", "a/..", "../outside", outside} {
		root, err := SetupConfigDir(name, false, SubaccountLevel)
		require.Error(t, err, name)
		require.Nil(t, root)
	}
	require.NoFileExists(t, filepath.Join(outside, "terraform.tfvars"))
	require.NoFileExists(t, filepath.Join(work, outside, "provider.tf"))
}

func TestConfigDirectoryRejectsSymlinkEscapes(t *testing.T) {
	work := t.TempDir()
	t.Chdir(work)
	root := openTestRoot(t, work)
	outside := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(outside, "victim"), 0700))
	sentinel := filepath.Join(outside, "victim", "keep.txt")
	require.NoError(t, os.WriteFile(sentinel, []byte("keep"), 0600))
	relativeOutside, err := filepath.Rel(work, outside)
	require.NoError(t, err)
	linkForTest(t, relativeOutside, filepath.Join(work, "link"))
	linkForTest(t, outside, filepath.Join(work, "absolute-link"))
	linkForTest(t, ".", filepath.Join(work, "self"))
	for _, name := range []string{"link", "link/victim", "link/new", "absolute-link/victim", "self"} {
		directory, err := SetupConfigDir(name, false, SubaccountLevel)
		require.Error(t, err, name)
		require.Nil(t, directory)
		require.Error(t, recreateExistingConfigDir(root, name), name)
	}
	require.FileExists(t, sentinel)
	require.NoDirExists(t, filepath.Join(outside, "new"))
}

func TestConfigDirectoryCreateOverwriteAndResume(t *testing.T) {
	t.Chdir(t.TempDir())
	stubProvider(t)
	require.NoError(t, os.Mkdir("parent", 0700))
	root, err := SetupConfigDir("parent/export", false, SubaccountLevel)
	require.NoError(t, err)
	require.NoError(t, root.WriteFile("keep", []byte("keep"), 0600))
	require.NoError(t, root.Close())
	choose := func(choice string) func(bool, *os.Root, string) (string, error) {
		return func(_ bool, _ *os.Root, _ string) (string, error) { return choice, nil }
	}
	root, err = setupConfigDir("parent/export", true, SubaccountLevel, choose(selectionResume))
	require.NoError(t, err)
	require.FileExists(t, filepath.Join("parent", "export", "keep"))
	require.NoError(t, root.Close())
	root, err = setupConfigDir("parent/export", true, SubaccountLevel, choose(selectionOverwrite))
	require.NoError(t, err)
	defer func() { _ = root.Close() }()
	require.NoFileExists(t, filepath.Join("parent", "export", "keep"))
	require.FileExists(t, filepath.Join("parent", "export", "provider.tf"))
	require.NoError(t, WriteExecutionTargetTFVars(root, ExecutionTarget{SubaccountID: "test"}))
	require.FileExists(t, filepath.Join("parent", "export", "terraform.tfvars"))
}

func TestResumeRejectsOutputSymlinks(t *testing.T) {
	t.Chdir(t.TempDir())
	stubProvider(t)
	require.NoError(t, os.Mkdir("export", 0700))
	outside := filepath.Join(t.TempDir(), "keep")
	require.NoError(t, os.WriteFile(outside, []byte("keep"), 0600))
	linkForTest(t, outside, filepath.Join("export", "provider.tf"))
	root, err := setupConfigDir("export", true, SubaccountLevel, func(bool, *os.Root, string) (string, error) { return selectionResume, nil })
	require.Error(t, err)
	require.Nil(t, root)
	content, err := os.ReadFile(outside)
	require.NoError(t, err)
	require.Equal(t, "keep", string(content))
}

func TestOverwriteRejectsDangerousDeletionTargets(t *testing.T) {
	root := openTestRoot(t, t.TempDir())
	require.NoError(t, root.WriteFile("keep", []byte("keep"), 0600))
	for _, path := range []string{"", ".", "..", "nested/..", root.Name(), "../outside"} {
		require.Error(t, recreateExistingConfigDir(root, path), path)
	}
	require.FileExists(t, filepath.Join(root.Name(), "keep"))
}

func TestMergeRejectsTraversalAndSymlinks(t *testing.T) {
	source := openTestRoot(t, t.TempDir())
	target := openTestRoot(t, t.TempDir())
	const name = "btp_resources.tf"
	require.NoError(t, source.WriteFile(name, []byte("generated"), 0600))
	for _, path := range []string{"", "../escape", "a/file", "a\\file", filepath.Join(t.TempDir(), "escape")} {
		require.Error(t, mergeTfConfig(target, path, source, "resource"))
	}
	outside := filepath.Join(t.TempDir(), "keep")
	require.NoError(t, os.WriteFile(outside, []byte("keep"), 0600))
	linkForTest(t, outside, filepath.Join(target.Name(), name))
	require.Error(t, mergeTfConfig(target, name, source, "resource"))
	content, err := os.ReadFile(outside)
	require.NoError(t, err)
	require.Equal(t, "keep", string(content))
}

func TestQueryAndTFVarsWritesRejectEscapes(t *testing.T) {
	directory := openTestRoot(t, t.TempDir())
	outside := filepath.Join(t.TempDir(), "keep")
	require.NoError(t, os.WriteFile(outside, []byte("keep"), 0600))
	linkForTest(t, outside, filepath.Join(directory.Name(), "resources.tf"))
	result := []QueryResult{{Config: "# config", ImportConfig: "# import"}}
	require.Error(t, WriteQueryResults(directory, "resources.tf", "btp_subaccount", result))
	require.Error(t, WriteQueryResults(directory, "../escape", "btp_subaccount", result))
	require.Error(t, WriteQueryResults(directory, "safe.tf", "../escape", result))
	require.NoFileExists(t, filepath.Join(directory.Name(), "safe.tf"))
	// A dangling external symlink must not turn a missing tfvars into a write
	// outside the output root.
	missing := filepath.Join(t.TempDir(), "terraform.tfvars")
	linkForTest(t, missing, filepath.Join(directory.Name(), "terraform.tfvars"))
	require.Error(t, WriteExecutionTargetTFVars(directory, ExecutionTarget{SubaccountID: "test"}))
	require.NoFileExists(t, missing)
	content, err := os.ReadFile(outside)
	require.NoError(t, err)
	require.Equal(t, "keep", string(content))
}

func TestTerraformPreflightRunsBeforeExecutable(t *testing.T) {
	root := openTestRoot(t, t.TempDir())
	linkForTest(t, t.TempDir(), filepath.Join(root.Name(), "escape"))
	previous := getIaCTool
	getIaCTool = func() (string, error) { t.Fatal("unsafe directory reached Terraform invocation"); return "", nil }
	t.Cleanup(func() { getIaCTool = previous })
	_, err := runTfCmdInDirectory(root, "fmt", "-recursive")
	require.ErrorContains(t, err, "unsafe Terraform path")
}

func TestContainedDirectorySymlinkAndOwnedQueryCleanup(t *testing.T) {
	t.Chdir(t.TempDir())
	stubProvider(t)
	require.NoError(t, os.Mkdir("parent", 0700))
	linkForTest(t, "parent", "alias")
	root, err := SetupConfigDir("alias/export", false, SubaccountLevel)
	require.NoError(t, err)
	defer func() { _ = root.Close() }()
	require.FileExists(t, filepath.Join("parent", "export", "provider.tf"))
	require.NoError(t, os.Mkdir("subaccount-config", 0700))
	require.NoError(t, os.WriteFile(filepath.Join("subaccount-config", "keep"), []byte("keep"), 0600))
	query, err := NewQueryDirectory(SubaccountLevel)
	require.NoError(t, err)
	t.Cleanup(func() { _ = query.Cleanup() })
	require.FileExists(t, filepath.Join(query.Root.Name(), "provider.tf"))
	CleanupTempFiles(query)
	require.NoDirExists(t, query.Root.Name())
	require.FileExists(t, filepath.Join("subaccount-config", "keep"))
}
