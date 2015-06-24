#!/bin/bash

set -e

THIS_SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "${THIS_SCRIPT_DIR}/.."

#
# Script for Continuous Integration
#

set -v

export PATH=$PATH:$GOPATH/bin
export GOPATH="${THIS_SCRIPT_DIR}/../Godeps/_workspace"
export PATH=$PATH:$GOPATH/bin

bash _scripts/install.sh

# Tests
go install github.com/kisielk/errcheck

errcheck ./

bash _test/test.sh

#
# ==> DONE - OK
#