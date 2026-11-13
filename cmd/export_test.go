package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExportImportsOnlyFlag(t *testing.T) {
	for _, command := range []*cobra.Command{exportByResourceCmd, exportByJsonCmd} {
		t.Run(command.Name(), func(t *testing.T) {
			flag := command.Flags().Lookup("imports-only")
			require.NotNil(t, flag)
			assert.Equal(t, "false", flag.DefValue)
			require.NoError(t, command.Flags().Parse([]string{"--imports-only"}))
			enabled, err := command.Flags().GetBool("imports-only")
			require.NoError(t, err)
			assert.True(t, enabled)
			t.Cleanup(func() {
				_ = command.Flags().Set("imports-only", "false")
				flag.Changed = false
				command.SetOut(nil)
			})
			var help bytes.Buffer
			command.SetOut(&help)
			require.NoError(t, command.Help())
			assert.Contains(t, help.String(), "--imports-only")
			assert.Contains(t, help.String(), "without resource configuration or NextSteps.md")
		})
	}
}
