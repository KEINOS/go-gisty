package gisty

import (
	"path/filepath"
	"testing"

	"github.com/cli/cli/v2/pkg/cmd/gist/create"
	"github.com/stretchr/testify/require"
)

func TestGisty_Create_msg_on_error(t *testing.T) {
	t.Parallel()

	obj := NewGisty()

	obj.AltFunctions.Create = func(*create.CreateOptions) error {
		return NewErr("forced error for creating")
	}

	// Execute the create command.
	argsCreate := CreateArgs{
		Description: "sample description",
		FilePaths: []string{
			filepath.Join("testdata", "foo.md"),
			filepath.Join("testdata", "bar.md"),
		},
		AsPublic: true,
	}

	gistURL, err := obj.Create(argsCreate)

	// Assert that the create command failed.
	require.Error(t, err)
	require.Nil(t, gistURL, "returned gistURL should be nil on error")
	require.Contains(t, err.Error(), "failed to execute create command")
	require.Contains(t, err.Error(), "forced error for creating")
}

func TestGisty_Create_unexpected_response(t *testing.T) {
	t.Parallel()

	obj := NewGisty()

	obj.AltFunctions.Create = func(createOpt *create.CreateOptions) error {
		// Set a dummy response.
		_, err := createOpt.IO.Out.Write([]byte("unexpected \nresponse"))
		require.NoError(t, err)

		return nil
	}

	// Execute the create command.
	argsCreate := CreateArgs{
		Description: "sample description",
		FilePaths: []string{
			filepath.Join("testdata", "foo.md"),
			filepath.Join("testdata", "bar.md"),
		},
		AsPublic: true,
	}

	gistURL, err := obj.Create(argsCreate)

	// Assert that the create command failed.
	require.Error(t, err)
	require.Nil(t, gistURL, "returned gistURL should be nil on error")
	require.Contains(t, err.Error(), "failed to parse gist URL")
}
