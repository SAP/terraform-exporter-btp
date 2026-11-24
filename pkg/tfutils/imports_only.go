package tfutils

import (
	"fmt"
	"os"
)

// removeResourceConfiguration runs after import transformations and address
// normalization, which still need the raw generated resource configuration.
// Removing only exporter-owned files also handles resumed exports.
func removeResourceConfiguration(directory *os.Root) error {
	for _, name := range []string{"btp_resources.tf", "NextSteps.md"} {
		if err := directory.Remove(name); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("remove %s: %w", name, err)
		}
	}
	return nil
}
