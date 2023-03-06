package buildinfo

import (
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/require"
)

//nolint:paralleltest // This test is not parallel because it changes the global variable.
func Test_getVersion_ver_is_set(t *testing.T) {
	oldVersion := Version
	defer func() {
		Version = oldVersion
	}()

	Version = "v1.2.3"

	require.Equal(t, "v1.2.3", getVersion())
}

//nolint:paralleltest // This test is not parallel because it changes the global variable.
func Test_getVersion_ver_is_empty(t *testing.T) {
	oldVersion := Version
	defer func() {
		Version = oldVersion
	}()

	Version = ""

	expect := "(devel)"
	actual := getVersion()

	require.Equal(t, expect, actual)
}

//nolint:paralleltest // This test is not parallel because it changes the global variable.
func Test_getVersion_from_build_info(t *testing.T) {
	oldVersion := Version
	oldDebugReadBuildInfo := debugReadBuildInfo

	defer func() {
		Version = oldVersion
		debugReadBuildInfo = oldDebugReadBuildInfo
	}()

	// mock debugReadBuildInfo to force return a version
	debugReadBuildInfo = func() (*debug.BuildInfo, bool) {
		//nolint:exhaustruct // this is a test
		return &debug.BuildInfo{
			Main: debug.Module{
				Version: "v1.2.3",
			},
		}, true
	}

	Version = ""

	expect := "v1.2.3"
	actual := getVersion()

	require.Equal(t, expect, actual)
}
