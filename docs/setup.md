# Setup

## Setup Nitro-cli

Check the steps [here](https://github.com/aws/aws-nitro-enclaves-cli/blob/main/docs/ubuntu_20.04_how_to_install_nitro_cli_from_github_sources.md)

## Setup Golang

Download the golang [here](https://go.dev/dl/)

## Build Image

```
# clone this repo
> cd image
> LANG=phala ./build-image.sh
* Build a docker image with the tag named: `ata-build-phala`
* Build a nitro enclave eif to `~/ata-build-phala-latest.eif`
* Deploy new bin
> 
```

## Configure the allocator

```
> sudo vim /etc/nitro_enclaves/allocator.yaml
> sudo systemctl restart nitro-enclaves-allocator
```

## Test

```
> cd ~/$REPO
> attestable-build-tool build -mem 10240
```