package tfutils

import (
	"fmt"
	"os"
	"path/filepath"
)

// removeResourceConfiguration runs after import transformations and address
// normalization, which still need the raw generated resource configuration.
// Removing only exporter-owned files also handles resumed exports.
func removeResourceConfiguration(directory string) error {
	for _, name := range []string{"btp_resources.tf", "NextSteps.md"} {
		if err := os.Remove(filepath.Join(directory, name)); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("remove %s: %w", name, err)
		}
	}
	return nil
}
