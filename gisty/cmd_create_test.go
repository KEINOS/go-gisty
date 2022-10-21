package gisty

import (
	"path/filepath"
	"testing"

	"github.com/cli/cli/v2/pkg/cmd/gist/create"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

// Disable the below test for now because the test actually creates a gist. (E2E test)
//
// func ExampleGisty_Create() {
// 	obj := NewGisty()
// 	argsCreate := CreateArgs{
// 		Description: "sample description",
// 		FilePaths: []string{
// 			filepath.Join("testdata", "foo.md"),
// 			filepath.Join("testdata", "bar.md"),
// 		},
// 		AsPublic: true,
// 	}

// 	gistURL, err := obj.Create(argsCreate)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println(gistURL.Path)
// 	fmt.Println("OK")

// 	// Output: OK
// }

func TestGisty_Create_msg_on_error(t *testing.T) {
	t.Parallel()

	obj := NewGisty()

	obj.AltFunctions.Create = func(*create.CreateOptions) error {
		return errors.New("forced error for creating")
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
	require.Contains(t, err.Error(), "failed to create gist")
	require.Contains(t, err.Error(), "forced error for creating")
}

func TestGisty_Create_unexpected_response(t *testing.T) {
	t.Parallel()

	obj := NewGisty()

	obj.AltFunctions.Create = func(createOpt *create.CreateOptions) error {
		// Set a dummy response.
		createOpt.IO.Out.Write([]byte("unexpected \nresponse"))

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
