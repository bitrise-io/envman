#!/bin/bash

set -e

THIS_SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "${THIS_SCRIPT_DIR}/.."

set +e

echo "-> Fmt"
go fmt

echo "-> Testing"
go test ./...
if [[ $? -ne 0 ]]; then
	echo " [!] Tests Failed"
	exit 1
fi

echo "-> Building"
go build
if [[ $? -ne 0 ]]; then
	echo " [!] Build Failed!"
	exit 1
fi

echo " (i) DONE - OK"
