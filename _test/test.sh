#!/bin/bash

set -e

THIS_SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "${THIS_SCRIPT_DIR}/.."

set -v

# ADD_CMD TESTS
# Add new environment variable 
# (list should contains: BITRISE:bitrise)
envman add --key BITRISE --value bitrise

# Add new environment variable from pipe 
# (list should contains: BITRISE_FROM_STDIN=bitrise from stdin)
echo "bitrise from stdin" | envman add --key BITRISE_FROM_STDIN

# PRINT_CMD TESTS
# Print out environment variables
# (should print added environment variables)
envman print

# RUN_CMD TESTS
# Run specified command 
# (should print "Bitrise value: bitrise")
envman run bash "$THIS_SCRIPT_DIR/runCmd_test.sh"