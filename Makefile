.PHONY: all test lint vuln
.PHONY: race fixwhere download update install

all: test lint vuln

test: download
	go test ./...
race: download
	go test -race ./...
# Run static analysis and linters. See .golangci.yml for the lint configuration.
# Requires: https://golangci-lint.run/
lint:
	golangci-lint run
# Run vulnerability scanner.
# Requires: https://go.dev/security/vuln/
vuln:
	govulncheck ./...
# Displays the location of the un-covered code.
# Requires: https://github.com/msoap/go-carpet
fixwhere:
	go-carpet -mincov 99.9
# Download all go module dependencies.
download:
	go mod download
# Update all go module dependencies.
update:
	./.github/update-go-mod.sh
# Install govulncheck if not already installed.
install:
	type govulncheck 2>/dev/null 1>/dev/null || go install golang.org/x/vuln/cmd/govulncheck@latest
