package gisty

import (
	"github.com/pkg/errors"
)

// NewErr returns a new error object with the given message appending the file
// name and line number of the caller.
//
// It is a wrapper of errors.New() and errors.Errorf(). Which is the alternative
// of deprecated github.com/pkg/errors.
func NewErr(msgs ...any) error {
	lenMsgs := len(msgs)
	errPos := getErrorPos()

	if lenMsgs == 0 {
		return nil
	}

	fmtErr, ok := msgs[0].(string)
	if !ok {
		errMsg := fmtArgs(msgs[:]...) //nolint: gocritic // false positive

		return errors.New(errMsg + errPos)
	}

	fmtErr += errPos

	if lenMsgs == 1 {
		return errors.New(fmtErr)
	}

	return errors.Errorf(fmtErr, msgs[1:]...)
}
