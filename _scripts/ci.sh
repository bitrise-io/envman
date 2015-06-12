#!/bin/bash

set -e

THIS_SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "${THIS_SCRIPT_DIR}/.."

#
# Script for Continuous Integration
#

set -v

# Go get
go get ./...

bash _scripts/install.sh

# Tests
bash _test/test.sh



#
# ==> DONE - OK
#