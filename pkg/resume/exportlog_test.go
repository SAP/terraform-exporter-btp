package resume

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExportLogRootedReadWriteAndRemoval(t *testing.T) {
	root, err := os.OpenRoot(t.TempDir())
	require.NoError(t, err)
	defer func() { _ = root.Close() }()
	require.NoError(t, WriteExportLog(root, "subaccount", "btp_subaccount", 1))
	require.NoError(t, WriteExportLog(root, "roles", "btp_subaccount_role", 2))
	resources, err := GetExistingExportLog(root)
	require.NoError(t, err)
	require.Equal(t, []string{"subaccount", "roles"}, resources)
	complete, err := GetExistingExportLogComplete(root)
	require.NoError(t, err)
	require.Len(t, complete.Resources, 2)
	require.NoError(t, RemoveExportLog(root))
	require.NoError(t, RemoveExportLog(root))
}

func TestExportLogRejectsExternalSymlink(t *testing.T) {
	root, err := os.OpenRoot(t.TempDir())
	require.NoError(t, err)
	defer func() { _ = root.Close() }()
	outside := filepath.Join(t.TempDir(), "keep.json")
	const content = "{\"resources\":[]}"
	require.NoError(t, os.WriteFile(outside, []byte(content), 0600))
	if err := os.Symlink(outside, filepath.Join(root.Name(), "importlog.json")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	require.Error(t, WriteExportLog(root, "subaccount", "btp_subaccount", 1))
	_, err = GetExistingExportLog(root)
	require.Error(t, err)
	_, err = GetExistingExportLogComplete(root)
	require.Error(t, err)
	require.Error(t, RemoveExportLog(root))
	actual, err := os.ReadFile(outside)
	require.NoError(t, err)
	require.Equal(t, content, string(actual))
}
