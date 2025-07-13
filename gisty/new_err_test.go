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
	err := gisty.NewErr()
	if err == nil {
		fmt.Println("empty args should return nil")
	}

	// Note the output contains the file name and line number of the caller.
	err = gisty.NewErr("simple error message")
	if err != nil {
		fmt.Println(err)
	}

	err = gisty.NewErr("%v error message", "formatted")
	if err != nil {
		fmt.Println(err)
	}

	err = gisty.NewErr("%v error message(s)", 3)
	if err != nil {
		fmt.Println(err)
	}

	err = gisty.NewErr(1, 2, 3)
	if err != nil {
		fmt.Println(err)
	}

	// Output:
	// empty args should return nil
	// simple error message (file: new_err_test.go, line: 22)
	// formatted error message (file: new_err_test.go, line: 27)
	// 3 error message(s) (file: new_err_test.go, line: 32)
	// 1 2 3 (file: new_err_test.go, line: 37)
}

// ============================================================================
//  Tests
// ============================================================================

func TestNewErr(t *testing.T) {
	t.Parallel()

	require.NoError(t, gisty.NewErr(), "empty args should return nil")
}
