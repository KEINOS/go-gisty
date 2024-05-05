/*
Package buildinfo provides version information about the current build.

The values "Version" and "Date" should be set via `-ldflags` option during build.
Sample command to build:

	$ VER_APP="$(git describe --tag)"
	$ go build -ldflags="-X 'main.Version=${VER_APP}'" ./path/to/main.go
*/
package buildinfo

import "runtime/debug"

// Version is dynamically set by the toolchain or overridden by the Makefile.
var Version = ""

// Date is dynamically set at build time in the Makefile.
var Date = "" // YYYY-MM-DD

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
