#!/bin/sh
# =============================================================================
#  Updates ./multihash/version.json, go.mod, and go.sum.
# =============================================================================

set -eu

echo '- Removing old go.mod and go.sum ...'
rm -f go.mod
rm -f go.sum

echo '- Regenerating go.mod ...'
go mod init "github.com/KEINOS/multihash"

echo '- Updating go.mod and go.sum ...'
go get "github.com/KEINOS/go-utiles"
go get "github.com/multiformats/go-multihash@master"
go get "github.com/pkg/errors"
go get "github.com/stretchr/testify"
go get "github.com/zenizh/go-capturer"
go mod tidy

echo '- Downloading latest version.json from upstream ...'
go generate ./...

echo '- Unit Testing ...'
go test -race ./...
