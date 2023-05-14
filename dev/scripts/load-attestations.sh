#!/bin/bash

set -e

ARCHIVISTASERVER="https://archivista.testifysec.localhost"

script_dir=$(dirname "$0")
cd "../$script_dir"
root_dir=$(pwd)


ATTESATION_DIR=$root_dir/.attestations

loopnum=1
total=$(ls -1 $ATTESATION_DIR | wc -l)

# Loop through attestations
for attestation in $ATTESATION_DIR/*; do
    echo "Loading attestation $loopnum of $total"
    set +e
	"$root_dir"/bin/archivistactl -u "$ARCHIVISTASERVER" store "$attestation"
    set -e
    loopnum=$((loopnum+1))
done