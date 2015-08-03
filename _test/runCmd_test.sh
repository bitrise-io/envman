#!/bin/bash

echo "Bitrise value: $BITRISE"
if [[ "$BITRISE" != "test bitrise value" ]] ; then
    exit 1
fi
