package gisty

import (
	"fmt"
	"path/filepath"
	"runtime"
)

// AppendErrPos is a flag to disable the file name and line number of the caller
// from the error message. If set to false, it will not be appended.
var AppendErrPos = true

// getErrorPos returns a string containing the file name and line number of the
// caller.
func getErrorPos() string {
	grandparent := 2 // 0 = self, 1 = parent, 2 = grandparent

	_, file, line, ok := runtime.Caller(grandparent)
	if !ok || !AppendErrPos {
		return ""
	}

	return fmt.Sprintf(" (file: %s, line: %d)", filepath.Base(file), line)
}
