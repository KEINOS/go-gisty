package gisty

import (
	"testing"

	"github.com/stretchr/testify/require"
)

//nolint:paralleltest // do not parallelize due to global variable change
func Test_getErrorPos(t *testing.T) {
	// Backup and defer restore the original value of AppendErrPos
	oldAppendErrPos := AppendErrPos
	defer func() {
		AppendErrPos = oldAppendErrPos
	}()

	//nolint:gocritic
	fn1 := func() string {
		return getErrorPos() // capture line number of the caller
	}

	fn2 := func() string {
		return fn1() // should return line 22
	}

	t.Run("disable append", func(t *testing.T) {
		// Disable append the error position info
		AppendErrPos = false

		result := fn2()

		require.Empty(t, result, "setting AppendErrPos to false should return empty string")
	})

	t.Run("enable append (default)", func(t *testing.T) {
		// Enable append the error position info
		AppendErrPos = true

		result := fn2()

		require.Contains(t, result, "file: get_error_pos_test.go",
			"on AppendErrPos true, the caller's file name should be included")
		require.Contains(t, result, "line:",
			"on AppendErrPos true, the caller's line number should be included")
	})
}
