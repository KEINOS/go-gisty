package gisty

import "os"

var (
	// osGetwd is a copy of os.Getwd to ease testing.
	osGetwd = os.Getwd
	// osChdir is a copy of os.Chdir to ease testing.
	osChdir = os.Chdir
)

// ChDir changes the current working directory to the given path and returns the
// previous working directory.
//
// It is the callers choice to change the working directory back to the previous
// working directory.
func ChDir(path string) (string, error) {
	returnPath, err := osGetwd()
	if err != nil {
		return "", WrapIfErr(err, "failed to get current working directory")
	}

	err = osChdir(path)
	if err != nil {
		return "", WrapIfErr(err, "failed to change working directory to %s", path)
	}

	return returnPath, nil
}
