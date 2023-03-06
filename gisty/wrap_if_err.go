package gisty

import (
	"fmt"

	"github.com/pkg/errors"
)

// WrapIfErr returns nil if err is nil.
//
// Otherwise, it returns an error annotating err with a stack trace at the point
// WrapIfErr is called. The supplied message contains the file name and line
// number of the caller.
//
// Note that if the "msgs" arg is more than one, the first arg is used as a
// format string and the rest are used as arguments.
//
// E.g.
//
//	WrapIfErr(nil, "it wil do nothing")
//	WrapIfErr(err)                                 // returns err as is
//	WrapIfErr(err, "failed to do something")       // eq to errors.Wrap
//	WrapIfErr(err, "failed to do %s", "something") // eq to errors.Wrapf
//
// It is a wrapper of errors.Wrap() and errors.Wrapf(). Which is the alternative
// of deprecated github.com/pkg/errors.
func WrapIfErr(err error, msgs ...any) error {
	if err == nil {
		return nil
	}

	if len(msgs) == 0 {
		return fmt.Errorf("%w", err)
	}

	errMsg := fmtArgs(msgs...)
	errPos := getErrorPos()

	return errors.Wrap(err, errMsg+errPos)
}
