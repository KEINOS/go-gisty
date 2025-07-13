package gisty

import (
	"os"
	"testing"

	"github.com/cli/cli/v2/pkg/cmd/repo/sync"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
//  Success cases
// ----------------------------------------------------------------------------

//nolint:paralleltest // do not parallelize due to temporary directory change
func TestGisty_Update_golden(t *testing.T) {
	chDirCleanUp(t) // defer change working dir back to original. see chdir_test.go

	// Instantiate the object.
	obj := NewGisty()
	args := NewUpdateArgs(t.TempDir())

	// Mock the update function.
	obj.AltFunctions.Update = func(*sync.SyncOptions) error {
		// On success, the output should include the following:
		const successMsgPfx = "✓ Synced"

		_, err := obj.Stdout.WriteString(successMsgPfx)
		require.NoError(t, err, "failed to write to stdout during mock")

		return nil
	}

	// Test
	_, err := obj.Update(args)

	require.NoError(t, err, "failed to update the gist")
}

//nolint:paralleltest // do not parallelize due to temporary directory change
func TestGisty_Update_golden_with_flags(t *testing.T) {
	chDirCleanUp(t) // defer change working dir back to original. see chdir_test.go

	// Instantiate the object.
	obj := NewGisty()

	// Set all the flags.

	branch := "main"
	remoteFork := "github.com/unknown/to.git"
	source := "github.com/unknown/from.git"
	force := true

	args := UpdateArgs{
		PathDirRepo: t.TempDir(),
		Branch:      branch,
		Destination: remoteFork,
		Source:      source,
		Force:       force,
	}

	// Mock the update function.
	obj.AltFunctions.Update = func(opt *sync.SyncOptions) error {
		// opt contains the given flags.
		require.Equal(t, branch, opt.Branch, "branch should be set")
		require.Equal(t, remoteFork, opt.DestArg, "remote fork should be set")
		require.Equal(t, source, opt.SrcArg, "source should be set")
		require.Equal(t, force, opt.Force, "force should be set")

		// On success, the output should include the following:
		const successMsgPfx = "✓ Synced"

		_, err := obj.Stdout.WriteString(successMsgPfx)
		require.NoError(t, err, "failed to write to stdout during mock")

		return nil
	}

	// Test
	_, err := obj.Update(args)

	require.NoError(t, err, "failed to update the gist")
}

// ----------------------------------------------------------------------------
//  Failure cases
// ----------------------------------------------------------------------------

//nolint:paralleltest // do not parallelize due to temporary directory change
func TestGisty_Update_execute_success_but_wrong_output(t *testing.T) {
	chDirCleanUp(t) // defer change working dir back to original. see chdir_test.go

	// Instantiate the object.
	obj := NewGisty()
	args := NewUpdateArgs(t.TempDir())

	// Mock the update function.
	obj.AltFunctions.Update = func(*sync.SyncOptions) error {
		const successMsgPfx = "success (no error) but unexpected output"

		_, err := obj.Stdout.WriteString(successMsgPfx)
		require.NoError(t, err, "failed to write to stdout during mock")

		return nil
	}

	// Test
	_, err := obj.Update(args)

	require.Error(t, err, "if update command fails to execute, it should return an error")
	assert.Contains(t, err.Error(), "failed to sync gist",
		"it should contain the error reason")
	assert.Contains(t, err.Error(), "success (no error) but unexpected output",
		"it should contain the original error")
}

//nolint:paralleltest // do not parallelize due to temporary directory change
func TestGisty_Update_fails_execute(t *testing.T) {
	chDirCleanUp(t) // defer change working dir back to original. see chdir_test.go

	// Instantiate the object.
	obj := NewGisty()
	args := NewUpdateArgs(t.TempDir())

	// Mock the update function.
	obj.AltFunctions.Update = func(*sync.SyncOptions) error {
		return NewErr("forced error")
	}

	// Test
	_, err := obj.Update(args)

	require.Error(t, err, "if update command fails to execute, it should return an error")
	assert.Contains(t, err.Error(), "failed to execute update/sync command",
		"it should contain the error reason")
	assert.Contains(t, err.Error(), "forced error",
		"it should contain the original error")
}

//nolint:paralleltest // do not parallelize due to temporary directory change
func TestGisty_Update_fail_to_change_dir(t *testing.T) {
	chDirCleanUp(t) // defer change working dir back to original. see chdir_test.go

	// Backup and defer restore the original osChdir function.
	oldOSChdir := osChdir

	defer func() {
		osChdir = oldOSChdir
	}()

	// Mock osChdir to return an error.
	osChdir = func(_ string) error {
		return NewErr("forced error")
	}

	// Instantiate the object.
	obj := NewGisty()
	args := NewUpdateArgs(t.TempDir())

	// Test
	_, err := obj.Update(args)

	require.Error(t, err, "expected error when osChdir returns an error")
	assert.Contains(t, err.Error(), "failed to change working directory to",
		"it should contain the error reason")
	assert.Contains(t, err.Error(), "forced error",
		"it should contain the original error")
}

//nolint:paralleltest // do not parallelize due to temporary directory change
func TestGisty_Update_fail_to_get_wd(t *testing.T) {
	chDirCleanUp(t) // defer change working dir back to original. see chdir_test.go

	// Backup and defer restore the original osGetwd function.
	oldOSGetwd := osGetwd

	defer func() {
		osGetwd = oldOSGetwd
	}()

	// Mock osGetwd to return an error.
	osGetwd = func() (string, error) {
		return "", NewErr("forced error")
	}

	// Instantiate the object.
	obj := NewGisty()
	args := NewUpdateArgs(t.TempDir())

	// Test
	_, err := obj.Update(args)

	require.Error(t, err, "expected error when osGetwd returns an error")
	assert.Contains(t, err.Error(), "failed to get current working directory",
		"it should contain the error reason")
	assert.Contains(t, err.Error(), "forced error",
		"it should contain the original error")
}

//nolint:paralleltest // do not parallelize due to temporary directory change
func TestGisty_Update_target_dir_is_empty(t *testing.T) {
	chDirCleanUp(t) // defer change working dir back to original. see chdir_test.go

	// Instantiate the object.
	obj := NewGisty()
	args := NewUpdateArgs("") // pathDirRepo is empty

	// Test
	_, err := obj.Update(args)

	require.Error(t, err, "expected error when target dir is empty")
	assert.Contains(t, err.Error(), "path to local repository is required",
		"it should contain the error reason")
}

// ----------------------------------------------------------------------------
//  Helper functions
// ----------------------------------------------------------------------------

// chDirCleanUp is a helper function that ensures the working directory is changed
// back to the original working directory.
func chDirCleanUp(t *testing.T) {
	t.Helper()

	pathDirOrig, err := os.Getwd()
	require.NoError(t, err, "failed to get current working directory during test setup")

	// Ensure we return to the original directory regardless of test outcome
	t.Cleanup(func() {
		// Only change back if the original directory still exists to avoid
		// conflicts when temporary directories are cleaned up
		_, err := os.Stat(pathDirOrig)
		if err == nil {
			// Change the working directory back to the original working directory.
			//nolint:usetesting // t.Chdir() has issues in Go 1.24, stick with os.Chdir()
			require.NoError(t, os.Chdir(pathDirOrig), "failed to change working directory back to %s", pathDirOrig)
		}
	})
}
