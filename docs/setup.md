# Setup

## Setup Nitro-cli

Check the steps [here](https://github.com/aws/aws-nitro-enclaves-cli/blob/main/docs/ubuntu_20.04_how_to_install_nitro_cli_from_github_sources.md)

## Setup Golang

Download the golang [here](https://go.dev/dl/)

## Build Image

Please ensure sufficient memory is available before building (recommended: 4GB).

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
memory_mib: 12288
cpu_count: 2
> sudo systemctl restart nitro-enclaves-allocator
```

## Test

```
> cd ~/$REPO
> attestable-build-tool build -mem 10240
```