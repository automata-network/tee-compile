#!/bin/bash -e

if [[ "$LANG" == "" ]]; then
    export LANG=rust
fi
if [[ "$NITRO_DIR" == "" ]]; then
    export NITRO_DIR=$HOME
fi

function build_tool() {
    `cd ../ && CGO_ENABLED=0 go build -o image/ .`
}

function build_docker() {
    build_tool
    docker build --tag ata-build-$LANG -f $LANG/Dockerfile .
    echo "* Build a docker image with the tag named: ata-build-$LANG"
}

function build_enclave() {
    nitro-cli build-enclave --docker-uri ata-build-$LANG:latest --output-file $NITRO_DIR/ata-build-$LANG-latest.eif
    echo "* Build a nitro enclave eif to $NITRO_DIR/ata-build-$LANG-latest.eif"
}

function deploy() {
    bin="attestable-build-tool"
    if [[ -f "/usr/bin/$bin" ]]; then
        old=$(cat /usr/bin/$bin | md5sum)
        new=$(cat $bin | md5sum)
        if [[ "$old" != "$new" ]]; then
            sudo cp $bin /usr/bin/$bin
            echo "* Deploy new bin"
        fi
    else
        sudo cp $bin /usr/bin/$bin
        echo "* Deploy new bin"
    fi
}

if [[ "$@" == "" ]]; then
	build_docker
	build_enclave
    deploy
else
	$@
fi
