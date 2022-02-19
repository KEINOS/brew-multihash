#!/bin/sh
# =============================================================================
#  Run basic tests.
# =============================================================================

# CONST
SUCCESS=0
FAILURE=1

# Unit Tests
printf "Running tests (unit test) ... "
LOG=$(go test -race -cover ./... 2>&1) || {
    printf '\033[31m%s\033[m\n' 'fail'
    echo >&2 "$LOG"
    exit $FAILURE
}
echo "$LOG"

# Formatting Tests
printf "Running golint (format check) ... "
LOG=$(golint ./... 2>&1) || {
    printf '\033[31m%s\033[m\n' 'fail'
    echo >&2 "$LOG"
    exit $FAILURE
}
echo "ok"

# Coverage Tests
printf "Running golangci-lint (static analysis) ... "
LOG=$(golangci-lint run --fix --color always 2>&1) || {
    printf '\033[31m%s\033[m\n' 'fail'
    echo >&2 "$LOG"
    exit $FAILURE
}
echo "ok"

exit $SUCCESS
