package gisty

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGisty_golden(t *testing.T) {
	t.Parallel()

	pathDirOrig := chDirCleanUp(t)
	pathDirTmp := t.TempDir()

	// Change the working directory to the temporary directory.
	returnPath, err := ChDir(pathDirTmp)
	require.NoError(t, err, "failed to change working directory to %s", pathDirTmp)

	// Test the return path.
	require.Equal(t, pathDirOrig, returnPath,
		"return path is not the original working directory")

	// Test the current working directory.
	pathDirCurr, err := os.Getwd()
	require.NoError(t, err, "failed to get current working directory")

	// Compare by substring.
	// On macOS, the actual temporary directory is under "/private" directory
	// and os.Getwd() returns the path without the "/private" prefix.
	require.Contains(t, pathDirCurr, pathDirTmp,
		"current working directory is not the temporary directory")
}

// ----------------------------------------------------------------------------
//  Helper functions
// ----------------------------------------------------------------------------

// chDirCleanUp is a helper function that ensures the working directory is changed
// back to the original working directory.
//
// The returned string is the original working directory to change back to.
func chDirCleanUp(t *testing.T) string {
	t.Helper()

	pathDirOrig, err := os.Getwd()
	require.NoError(t, err, "failed to get current working directory during test setup")

	t.Cleanup(func() {
		// Change the working directory back to the original working directory.
		require.NoError(t, os.Chdir(pathDirOrig), "failed to change working directory back to %s", pathDirOrig)
	})

	return pathDirOrig
}
