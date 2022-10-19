package build

import "runtime/debug"

// Version is dynamically set by the toolchain or overridden by the Makefile.
var Version = ""

// Date is dynamically set at build time in the Makefile.
var Date = "" // YYYY-MM-DD

//nolint:gochecknoinits // We initialize the global variable here.
func init() {
	Version = getVersion()
}

// debugReadBuildInfo is a copy of debug.ReadBuildInfo to ease testing.
var debugReadBuildInfo = debug.ReadBuildInfo

func getVersion() string {
	if Version != "" {
		return Version
	}

	if info, ok := debugReadBuildInfo(); ok {
		version := info.Main.Version
		if version != "" {
			return version
		}
	}

	return "(devel)"
}
