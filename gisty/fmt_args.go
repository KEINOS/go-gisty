package gisty

import "fmt"

// fmtArgs is a shorthand for fmt.Sprintf and fmt.Sprint arguments. It is a helper
// function to format the given arguments.
//
// If the inputs is empty, it returns an empty string.
// If the inputs has only one element, it returns the string representation of
// the element.
// If the inputs has more than one element, the first element is used as a
// format string and the rest are used as arguments.
func fmtArgs(inputs ...any) string {
	lenInput := len(inputs)

	if lenInput == 0 {
		return ""
	}

	if lenInput == 1 {
		return fmt.Sprint(inputs[0])
	}

	format, ok := inputs[0].(string)
	if !ok {
		return fmt.Sprint(inputs[:]...) //nolint:gocritic // false positive
	}

	return fmt.Sprintf(format, inputs[1:]...)
}
