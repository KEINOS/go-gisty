package gisty

import (
	"testing"

	"github.com/cli/cli/v2/pkg/cmd/gist/clone"
	"github.com/stretchr/testify/require"
)

func TestGisty_Clone_msg_on_error(t *testing.T) {
	t.Parallel()

	obj := NewGisty()

	// Execute the clone command.
	targetGist := "https://gist.github.com/7101f542be23e5048198e2a27c3cfda8.git"
	outDir := t.TempDir()
	args := []string{targetGist, outDir}

	obj.AltFunctions.Clone = func(*clone.CloneOptions) error {
		return NewErr("forced error for cloning")
	}

	err := obj.Clone(args)

	// Assert that the clone command failed.
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to execute gist clone")
	require.Contains(t, err.Error(), "forced error for cloning")
}
