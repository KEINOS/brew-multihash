#!/bin/sh
# =============================================================================
#  Builds binary under ./dist/ directory.
# =============================================================================
#  Note that this script will not upload the binary to GitHub Releases.

goreleaser release --snapshot --skip-publish --rm-dist
