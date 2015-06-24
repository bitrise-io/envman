#!/bin/bash

set -e

THIS_SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "${THIS_SCRIPT_DIR}/.."

#
# Script for Continuous Integration
#

set -v


export PATH=$PATH:$GOPATH/bin

# Go get
go get ./...

bash _scripts/install.sh

# Tests
go get github.com/kisielk/errcheck
go install github.com/kisielk/errcheck

errcheck ./

bash _test/test.sh

#
# ==> DONE - OK
#
