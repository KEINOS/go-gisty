package gisty

import (
	"testing"

	"github.com/cli/cli/v2/pkg/cmd/gist/delete"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestGisty_Delete_msg_on_error(t *testing.T) {
	t.Parallel()

	obj := NewGisty()

	// Execute the delete command.
	targetGist := "https://gist.github.com/unknown.git"

	obj.AltFunctions.Delete = func(*delete.DeleteOptions) error {
		return errors.New("forced error for deleting")
	}

	err := obj.Delete(targetGist)

	// Assert that the delete command failed.
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to delete gist")
	require.Contains(t, err.Error(), "forced error for deleting")
}
