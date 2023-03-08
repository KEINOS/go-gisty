test:
	go test ./...
race:
	go test -race ./...
lint:
	# Run static analysis and linters. See .golangci.yml for the lint configuration.
	# Requires: https://golangci-lint.run/
	golangci-lint run
fixwhere:
	# Displays the location of the un-covered code.
	# Requires: https://github.com/msoap/go-carpet
	go-carpet -mincov 99.9
