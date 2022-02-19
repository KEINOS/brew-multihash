#!/bin/sh
# Prints the version of the current project.
# It appends the current commit hash to the Upstream version.
pathFileJSON="./cmd/multihash/version.json"
printf "%s" "$(cat "$pathFileJSON" | jq -r .version)-$(git rev-parse --short HEAD)"
