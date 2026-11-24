package files

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// ValidateTerraformDirectory checks the path-based view a Terraform subprocess
// will use. Terraform does not inherit os.Root containment: the caller must
// trust the directory against concurrent modification during its execution.
func ValidateTerraformDirectory(directory *os.Root) error {
	want, err := directory.Stat(".")
	if err != nil {
		return err
	}
	got, err := os.Stat(directory.Name())
	if err != nil {
		return err
	}
	if !os.SameFile(want, got) {
		return fmt.Errorf("terraform directory %q has been replaced", directory.Name())
	}
	return fs.WalkDir(directory.FS(), ".", func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.Type()&os.ModeSymlink != 0 {
			if _, err := directory.Stat(filepath.FromSlash(path)); err != nil {
				return fmt.Errorf("unsafe Terraform path %q: %w", path, err)
			}
		} else if !entry.IsDir() && !entry.Type().IsRegular() {
			return fmt.Errorf("terraform path %q is not a regular file", path)
		}
		return nil
	})
}
