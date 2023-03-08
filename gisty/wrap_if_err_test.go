package gisty_test

import (
	"fmt"

	"github.com/KEINOS/go-gisty/gisty"
)

// ============================================================================
//  Examples
// ============================================================================

//nolint:lll
func ExampleWrapIfErr() {
	var err error

	// WrapIfErr returns nil if err is nil
	fmt.Println("err is nil:", gisty.WrapIfErr(err, "error at line 18"))

	// Cause err to be non-nil
	err = gisty.NewErr("error occurred at line 21")
	// Wrap with no additional message
	fmt.Println("err is non-nil:\n", gisty.WrapIfErr(err))
	// Wrap with additional message
	fmt.Println("err is non-nil:\n", gisty.WrapIfErr(err, "wrapped at line 25"))
	// Output:
	// err is nil: <nil>
	// err is non-nil:
	//  error occurred at line 21 (file: wrap_if_err_test.go, line: 21)
	// err is non-nil:
	//  wrapped at line 25 (file: wrap_if_err_test.go, line: 25): error occurred at line 21 (file: wrap_if_err_test.go, line: 21)
}

//nolint:lll
func ExampleWrapIfErr_disable_error_position() {
	// Backup and defer restore the original value of AppendErrPos
	oldAppendErrPos := gisty.AppendErrPos
	defer func() {
		gisty.AppendErrPos = oldAppendErrPos
	}()

	{
		gisty.AppendErrPos = false // Disable appending the error position

		err := gisty.NewErr("error occurred at line 45")
		fmt.Println(gisty.WrapIfErr(err, "wrapped at line 46"))
	}
	{
		gisty.AppendErrPos = true // Enable appending the error position (default)

		err := gisty.NewErr("error occurred at line 51")
		fmt.Println(gisty.WrapIfErr(err, "wrapped at line 52"))
	}
	// Output:
	// wrapped at line 46: error occurred at line 45
	// wrapped at line 52 (file: wrap_if_err_test.go, line: 52): error occurred at line 51 (file: wrap_if_err_test.go, line: 51)
}
