package gisty_test

import (
	"fmt"
	"testing"

	"github.com/KEINOS/go-gisty/gisty"
	"github.com/stretchr/testify/require"
)

// ============================================================================
//  Examples
// ============================================================================

func ExampleNewErr() {
	if err := gisty.NewErr(); err == nil {
		fmt.Println("empty args returns nil")
	}

	// Note the output contains the file name and line number of the caller.
	if err := gisty.NewErr("simple error message"); err != nil {
		fmt.Println(err)
	}

	if err := gisty.NewErr("%v error message", "formatted"); err != nil {
		fmt.Println(err)
	}

	if err := gisty.NewErr("%v error message(s)", 3); err != nil {
		fmt.Println(err)
	}

	if err := gisty.NewErr(1, 2, 3); err != nil {
		fmt.Println(err)
	}

	// Output:
	// empty args returns nil
	// simple error message (file: new_err_test.go, line: 21)
	// formatted error message (file: new_err_test.go, line: 25)
	// 3 error message(s) (file: new_err_test.go, line: 29)
	// 1 2 3 (file: new_err_test.go, line: 33)
}

// ============================================================================
//  Tests
// ============================================================================

func TestNewErr(t *testing.T) {
	t.Parallel()

	require.NoError(t, gisty.NewErr(), "empty args should return nil")
}
