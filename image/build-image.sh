#!/bin/bash -e

if [[ "$LANG" == "" ]]; then
    export LANG=rust
fi

function build_docker() {
    `cd ../ && CGO_ENABLED=0 go build -o image/ .`
    docker build --tag ata-build-$LANG -f $LANG/Dockerfile .
}

function build_enclave() {
    nitro-cli build-enclave --docker-uri ata-build-$LANG:latest --output-file ~/ata-build-$LANG-latest.eif
}

if [[ "$@" == "" ]]; then
	build_docker
	build_enclave
else
	$@
fi
