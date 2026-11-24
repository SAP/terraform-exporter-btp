package files

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func testRoot(t *testing.T, path string) *os.Root {
	t.Helper()
	root, err := os.OpenRoot(path)
	require.NoError(t, err)
	t.Cleanup(func() { _ = root.Close() })
	return root
}

func symlink(t *testing.T, target, name string) {
	t.Helper()
	if err := os.Symlink(target, name); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
}

func TestValidateConfigDir(t *testing.T) {
	for _, name := range []string{"", ".", "..", "a/..", "../outside", "a/../../outside", t.TempDir(), "bad\x00name"} {
		t.Run(name, func(t *testing.T) { require.Error(t, ValidateConfigDir(name)) })
	}
	for _, name := range []string{"export", "./export", "parent/export", "a/../export", "[export]", "..export"} {
		require.NoError(t, ValidateConfigDir(name), name)
	}
}

func TestValidateFilename(t *testing.T) {
	for _, name := range []string{"", ".", "..", "../outside", "nested/file", "nested\\file", "C:output", "/outside", "bad\x00name"} {
		require.Error(t, ValidateFilename(name), name)
	}
	require.NoError(t, ValidateFilename("btp_subaccount_import.tf"))
}

func TestImportFilesStayWithinRoots(t *testing.T) {
	for _, side := range []string{"source", "destination"} {
		t.Run(side, func(t *testing.T) {
			source := testRoot(t, t.TempDir())
			destination := testRoot(t, t.TempDir())
			outside := filepath.Join(t.TempDir(), "sentinel.tf")
			require.NoError(t, os.WriteFile(outside, []byte("untouched"), 0600))
			const name = "btp_subaccount_import.tf"
			if side == "source" {
				symlink(t, outside, filepath.Join(source.Name(), name))
			} else {
				require.NoError(t, source.WriteFile(name, []byte("import"), 0600))
				relative, err := filepath.Rel(destination.Name(), outside)
				require.NoError(t, err)
				symlink(t, relative, filepath.Join(destination.Name(), name))
			}
			require.Error(t, CopyImportFiles(source, destination))
			content, err := os.ReadFile(outside)
			require.NoError(t, err)
			require.Equal(t, "untouched", string(content))
		})
	}
}

func TestImportFileNamesAndRepeatedCopies(t *testing.T) {
	// Glob metacharacters in a directory name must be literal, not expand into
	// unrelated directories when selecting the source files.
	parent := t.TempDir()
	sourcePath := filepath.Join(parent, "[source]")
	require.NoError(t, os.Mkdir(sourcePath, 0700))
	source := testRoot(t, sourcePath)
	destination := testRoot(t, t.TempDir())
	require.NoError(t, WriteImportConfiguration(source, "btp_subaccount", "import\n"))
	require.NoError(t, source.WriteFile("ignore.tf", []byte("ignore"), 0600))
	require.NoError(t, CopyImportFiles(source, destination))
	require.NoError(t, CopyImportFiles(source, destination))
	content, err := destination.ReadFile("btp_subaccount_import.tf")
	require.NoError(t, err)
	require.Equal(t, "import\nimport\n", string(content))
	_, err = destination.Stat("ignore.tf")
	require.True(t, os.IsNotExist(err))
	for _, name := range []string{"../escape", "a/b", "a\\b", "", "."} {
		require.Error(t, WriteImportConfiguration(source, name, "bad"))
		require.Error(t, CopyFile(source, destination, name, os.O_TRUNC))
	}
}

func TestTempDirectoryOwnership(t *testing.T) {
	first, err := NewTempDirectory("btptf-test-")
	require.NoError(t, err)
	t.Cleanup(func() { _ = first.Cleanup() })
	second, err := NewTempDirectory("btptf-test-")
	require.NoError(t, err)
	t.Cleanup(func() { _ = second.Cleanup() })
	require.NotEqual(t, first.Root.Name(), second.Root.Name())
	outside := filepath.Join(t.TempDir(), "keep")
	require.NoError(t, os.WriteFile(outside, []byte("keep"), 0600))
	symlink(t, outside, filepath.Join(first.Root.Name(), "link"))
	require.NoError(t, first.Cleanup())
	require.NoDirExists(t, first.Root.Name())
	require.FileExists(t, outside)
	require.DirExists(t, second.Root.Name())
	require.NoError(t, first.Cleanup())
	require.NoError(t, (*TempDirectory)(nil).Cleanup())
}

func TestTempDirectoryRejectsReplacedPath(t *testing.T) {
	directory, err := NewTempDirectory("btptf-test-")
	require.NoError(t, err)
	original := directory.Root.Name()
	moved := original + "-moved"
	t.Cleanup(func() { _ = directory.Cleanup(); _ = os.Remove(original); _ = os.RemoveAll(moved) })
	require.NoError(t, os.Rename(original, moved))
	outside := t.TempDir()
	sentinel := filepath.Join(outside, "keep")
	require.NoError(t, os.WriteFile(sentinel, []byte("keep"), 0600))
	symlink(t, outside, original)
	require.ErrorContains(t, directory.Cleanup(), "replaced temporary directory")
	require.FileExists(t, sentinel)
	require.DirExists(t, moved)
}

func TestTerraformDirectoryPreflight(t *testing.T) {
	for _, name := range []string{"provider.tf", "nested", ".terraform/providers/cache", ".terraform/modules/module"} {
		t.Run(name, func(t *testing.T) {
			root := testRoot(t, t.TempDir())
			require.NoError(t, root.MkdirAll(filepath.Dir(name), 0700))
			symlink(t, t.TempDir(), filepath.Join(root.Name(), name))
			require.ErrorContains(t, ValidateTerraformDirectory(root), "unsafe Terraform path")
		})
	}
	t.Run("contained relative link", func(t *testing.T) {
		root := testRoot(t, t.TempDir())
		require.NoError(t, root.WriteFile("real.tf", []byte("# config"), 0600))
		symlink(t, "real.tf", filepath.Join(root.Name(), "linked.tf"))
		require.NoError(t, ValidateTerraformDirectory(root))
	})
	t.Run("replaced subprocess directory", func(t *testing.T) {
		parent := t.TempDir()
		name := filepath.Join(parent, "export")
		require.NoError(t, os.Mkdir(name, 0700))
		root := testRoot(t, name)
		require.NoError(t, os.Rename(name, name+"-moved"))
		require.NoError(t, os.Mkdir(name, 0700))
		require.ErrorContains(t, ValidateTerraformDirectory(root), "has been replaced")
		require.NoError(t, root.WriteFile("safe.tf", []byte("safe"), 0600))
		require.NoFileExists(t, filepath.Join(name, "safe.tf"))
		require.FileExists(t, filepath.Join(name+"-moved", "safe.tf"))
	})
}
