#!/bin/bash

set -e

THIS_SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "${THIS_SCRIPT_DIR}/.."

set -v

mkdir -p "$HOME/.envman/test"

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
envman --path "$HOME/.envman/test/.envstore.yml" init
envman --path "$HOME/.envman/test/.envstore.yml" print
echo "bitrise from stdin" | envman --path "$HOME/.envman/test/.envstore.yml" add --key BITRISE_FROM_STDIN
envman --path "$HOME/.envman/test/.envstore.yml" add --key BITRISE --value "bitrise with path flag"
envman --path "$HOME/.envman/test/.envstore.yml" run bash "$THIS_SCRIPT_DIR/runCmd_test.sh"
envman --path "$HOME/.envman/test/.envstore.yml" clear
envman --path "$HOME/.envman/test/.envstore.yml" print


# ****************************************************
# Init in current path, if .envstore.yml exist (exist)
# ****************************************************
cd "$HOME/.envman/test/"
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
set +e
rm -rf "$HOME/.envman/test-emtpy"
set -e
mkdir -p "$HOME/.envman/test-emtpy"

cd "$HOME/.envman/test-emtpy"
envman init
envman print
echo "bitrise from stdin" | envman add --key BITRISE_FROM_STDIN
envman add --key BITRISE --value "bitrise"
envman run bash "$THIS_SCRIPT_DIR/runCmd_test.sh"
envman clear
envman print
