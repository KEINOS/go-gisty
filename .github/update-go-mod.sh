#!/bin/sh
# =============================================================================
#  This script updates Go modules to the latest version.
# =============================================================================
#  It updates dependencies for the Go version specified in `go_ver_min`, then
#  runs `go mod tidy` and the test suite.
#  If the tests fail, it discards all tracked changes with `git restore .` and
#  exits with status 1.
#  On success, it leaves any updates in the worktree for a later commit.
#
#  NOTE: This script is aimed to run in the container via docker-compose.
#    See "tidy" service: ./docker-compose.yml
# =============================================================================

go_ver_min="1.26.1"

set -eu

echo '* Updating modules ...'
go get -u ./...

echo '* Run go tidy ...'
go mod tidy -go "$go_ver_min"

echo '* Run tests ...'
if ! go test ./...; then
	echo '* Tests failed. Discarding all tracked changes ...' >&2
	git restore .
	exit 1
fi

echo 'Tests passed. Module updates are ready for review.'
