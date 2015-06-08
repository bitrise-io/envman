#!/bin/bash

set -e

THIS_SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "${THIS_SCRIPT_DIR}/.."

set -v

# Format
go fmt

# Test
go test ./...

# Build
go build

# DONE - OK
