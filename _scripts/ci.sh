#!/bin/bash

set -e

THIS_SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "${THIS_SCRIPT_DIR}/.."

#
# Script for Continuous Integration
#

# export PATH=$PATH:$GOPATH/bin
# export GOPATH="${THIS_SCRIPT_DIR}/../Godeps/_workspace"
# echo "GOPATH: $GOPATH"
# mkdir $GOPATH/bin
# export PATH=$PATH:$GOPATH/bin

set -v

go get github.com/tools/godep
godep restore

bash _scripts/install.sh

# Tests
go get github.com/kisielk/errcheck
go install github.com/kisielk/errcheck

errcheck ./

bash _test/test.sh

#
# ==> DONE - OK
#