FROM golang:alpine

RUN \
    apk add --no-cache \
    git \
    curl \
    #alpine-sdk \
    build-base \
    \
    && \
    # Install latest golangci-lint, the meta-linter for go
    wget -O- https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "${GOPATH:?Undefined}/bin" \
    \
    && \
    # Install latest golint for Codevoc.io compatibility
    go install "golang.org/x/lint/golint@latest" \
    \
    && \
    # Install latest goreleaser
    go install "github.com/goreleaser/goreleaser@latest"

WORKDIR /workspaces

COPY ./go.mod /workspaces/go.mod

# Initial cache
RUN go mod download
