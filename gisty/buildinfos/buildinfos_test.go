package buildinfos

import (
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/require"
)

const testVersion123 = "v1.2.3"

//nolint:paralleltest // This test is not parallel because it changes the global variable.
func Test_getVersion_ver_is_set(t *testing.T) {
	oldVersion := Version

	defer func() {
		Version = oldVersion
	}()

	Version = testVersion123

	require.Equal(t, testVersion123, getVersion())
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
				Version: testVersion123,
			},
		}, true
	}

	Version = ""

	expect := testVersion123
	actual := getVersion()

	require.Equal(t, expect, actual)
}

//nolint:paralleltest // This test is not parallel because it changes global variables.
func Test_getVersion_from_build_info_with_empty_version(t *testing.T) {
	oldVersion := Version
	oldDebugReadBuildInfo := debugReadBuildInfo

	defer func() {
		Version = oldVersion
		debugReadBuildInfo = oldDebugReadBuildInfo
	}()

	debugReadBuildInfo = func() (*debug.BuildInfo, bool) {
		//nolint:exhaustruct // this is a test
		return &debug.BuildInfo{
			Main: debug.Module{
				Version: "",
			},
		}, true
	}

	Version = ""

	require.Equal(t, "(devel)", getVersion())
}
