#!/bin/bash

set -e

THIS_SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "${THIS_SCRIPT_DIR}/.."

set -v


# Build
docker build -t envman-img .


# Run
docker run --rm -it -v `pwd`:/go/src/github.com/bitrise-io/envman --name envman-dev envman-img /bin/bash
