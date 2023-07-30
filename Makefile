test: download
	go test ./...
race: download
	go test -race ./...
lint:
	# Run static analysis and linters. See .golangci.yml for the lint configuration.
	# Requires: https://golangci-lint.run/
	golangci-lint run
vuln:
	# Run vulnerability scanner.
	# Requires: https://go.dev/security/vuln/
	govulncheck ./...
fixwhere:
	# Displays the location of the un-covered code.
	# Requires: https://github.com/msoap/go-carpet
	go-carpet -mincov 99.9
download:
	# Download all go module dependencies.
	go mod download
update:
	# Update all go module dependencies.
	./.github/update-go-mod.sh
install:
	# Install govulncheck if not already installed.
	type govulncheck 2>/dev/null 1>/dev/null || go install golang.org/x/vuln/cmd/govulncheck@latest
all: test lint vuln
