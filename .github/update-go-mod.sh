#!/bin/sh
# =============================================================================
#  This script updates Go modules to the latest version.
# =============================================================================
#  It will remove the go.mod file and run `go mod tidy` to get the latest moule
#  versions.
#  Then it will run the tests to make sure the code is still working, and fails
#  if any errors are found during the process.
#
#  NOTE: This script is aimed to run in the container via docker-compose.
#    See "tidy" service: ./docker-compose.yml
# =============================================================================

go_ver_min="1.18"

set -eu

echo '* Backup module files ...'
cp go.mod go.mod.bak
cp go.sum go.sum.bak

# name_package="github.com/KEINOS/go-gisty"
#
# echo '* Create new blank go.mod ...'
# go mod init "${name_package}"
#
# echo '* Add the package to the go.mod ...'
# go get \
# 	github.com/stretchr/testify \
# 	github.com/cli/cli/v2 \
# 	github.com/alessio/shellescape \
# 	github.com/pkg/errors

echo '* Updating modules ...'
go get -u ./...

echo '* Run go tidy ...'
go mod tidy -go "$go_ver_min"

echo '* Run tests ...'
go test ./... && {
	echo '* Testing passed. Removing old go.mod file ...'
	rm -f go.mod.bak
	rm -f go.sum.bak
	echo 'Successfully updated modules!'
}
