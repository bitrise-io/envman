#!/bin/bash

set -e

THIS_SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "${THIS_SCRIPT_DIR}/.."

set -v


# Build
docker build -t envman-img .

# Run
bash _scripts/docker_run.sh
