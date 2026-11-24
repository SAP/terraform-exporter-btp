package files

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ValidateConfigDir checks the CLI path before it is cleaned or joined. All
// filesystem operations must additionally use a Root opened on the working
// directory; lexical validation alone cannot prevent symlink traversal.
func ValidateConfigDir(name string) error {
	if !filepath.IsLocal(name) || filepath.Clean(name) == "." || strings.ContainsRune(name, 0) {
		return fmt.Errorf("configuration directory %q must be a relative subdirectory of the working directory", name)
	}
	return nil
}

// ValidateFilename accepts one portable filename component, never a path.
func ValidateFilename(name string) error {
	if !filepath.IsLocal(name) || name == "." || strings.ContainsAny(name, "/\\:\x00") {
		return fmt.Errorf("invalid filename %q: expected a single filename component", name)
	}
	return nil
}

// TempDirectory owns a randomly allocated directory. Only this constructor
// grants recursive cleanup authority; callers cannot supply a deletion path.
type TempDirectory struct {
	Root   *os.Root
	parent *os.Root
	name   string
}

func NewTempDirectory(prefix string) (*TempDirectory, error) {
	if err := ValidateFilename(prefix); err != nil {
		return nil, err
	}
	path, err := os.MkdirTemp("", prefix)
	if err != nil {
		return nil, err
	}
	parent, err := os.OpenRoot(filepath.Dir(path))
	if err != nil {
		_ = os.Remove(path) // Still empty and allocated by this call.
		return nil, err
	}
	name := filepath.Base(path)
	root, err := parent.OpenRoot(name)
	if err != nil {
		_ = parent.Remove(name)
		_ = parent.Close()
		return nil, err
	}
	return &TempDirectory{Root: root, parent: parent, name: name}, nil
}

func (d *TempDirectory) Cleanup() error {
	if d == nil || d.parent == nil {
		return nil
	}
	defer func() {
		_ = d.Root.Close()
		_ = d.parent.Close()
		d.parent = nil
	}()
	want, err := d.Root.Stat(".")
	if err != nil {
		return err
	}
	got, err := d.parent.Lstat(d.name)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if !os.SameFile(want, got) {
		return fmt.Errorf("refusing to remove replaced temporary directory %q", d.name)
	}
	// Close first so Windows can remove the directory. RemoveAll stays rooted
	// in the allocation parent and does not follow a replacement final symlink.
	return errors.Join(d.Root.Close(), d.parent.RemoveAll(d.name))
}
