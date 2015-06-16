#!/bin/bash

set -e

THIS_SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "${THIS_SCRIPT_DIR}/.."

set -v


# Run
docker run --rm -it --volume `pwd`:/go/src/github.com/bitrise-io/envman --name envman-dev envman-img /bin/bash
# docker run --rm -it envman-img /bin/bash
