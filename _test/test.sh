#!/bin/bash

set -e

THIS_SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "${THIS_SCRIPT_DIR}/.."

set -v

# **************************
# Init in default envman dir
# **************************
envman init
envman print
echo "bitrise from stdin" | envman add --key BITRISE_FROM_STDIN
envman add --key BITRISE --value bitrise
envman run bash "$THIS_SCRIPT_DIR/runCmd_test.sh"
envman clear
envman print

# **********************
# Init in specified path
# **********************
envman --path /Users/godrei/dev/.envstore.yml init
envman --path /Users/godrei/dev/.envstore.yml print
echo "bitrise from stdin" | envman --path /Users/godrei/dev/.envstore.yml add --key BITRISE_FROM_STDIN
envman --path /Users/godrei/dev/.envstore.yml add --key BITRISE --value "bitrise with path flag"
envman --path /Users/godrei/dev/.envstore.yml run bash "$THIS_SCRIPT_DIR/runCmd_test.sh"
envman --path /Users/godrei/dev/.envstore.yml clear
envman --path /Users/godrei/dev/.envstore.yml print

# ****************************************************
# Init in current path, if .envstore.yml exist (exist)
# ****************************************************
cd /Users/godrei/dev/
envman init
envman print
echo "bitrise from stdin" | envman add --key BITRISE_FROM_STDIN
envman add --key BITRISE --value "bitrise from current path"
envman run bash "$THIS_SCRIPT_DIR/runCmd_test.sh"
envman clear
envman print

# ********************************************************
# Init in current path, if .envstore.yml exist (not exist)
# ********************************************************
cd /Users/godrei
envman init
envman print
echo "bitrise from stdin" | envman add --key BITRISE_FROM_STDIN
envman add --key BITRISE --value "bitrise"
envman run bash "$THIS_SCRIPT_DIR/runCmd_test.sh"
envman clear
envman print
