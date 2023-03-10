# =============================================================================
#  Test Container for Vaious Go Versions (ver. 2023.03.10-15:00)
# =============================================================================

# Default version
ARG VARIANT="alpine"

# -----------------------------------------------------------------------------
#  Main Stage
# -----------------------------------------------------------------------------
FROM golang:${VARIANT}

ENV GO111MODULE=on

RUN apk add --no-cache \
	build-base \
    alpine-sdk \
    git \
	make

WORKDIR /workspaces

ENTRYPOINT go mod download && go test -race ./...
