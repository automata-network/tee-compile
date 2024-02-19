# Attestable Build Tool

[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## About

In the current technological landscape, there's a glaring absence of mechanisms to verify that an executable file has been compiled from a specific source code. This gap poses significant challenges in terms of security, transparency, and trust, as there is no definitive way to prove the authenticity of the compiled software.

To address this issue, we are introducing a method that involves standardizing the compilation process within an enclave environment. An enclave is a protected, isolated area of execution, where code can be run in confidentiality and integrity, safeguarded from potential tampering or unauthorized access.

## Architecture

![architecture](https://1440033567-files.gitbook.io/~/files/v0/b/gitbook-x-prod.appspot.com/o/spaces%2FtYKuUrKWPlgYjy0suCeT%2Fuploads%2F5z6W0W8hyBeYLHmfKMB1%2Fimage.png?alt=media&token=ac96957d-78ef-4717-991d-ce4093ce912e)


## Usage

### Github Runner
1. Apply for a nitro enclave machine on AWS.
2. Configure the GitHub runner:  
    2.1. Settings → Actions → Runners → New self-hosted runner  
    2.2. Follow the instructions to configure the GitHub Runner  
3. Download the Software Build Attestation Image.
4. Download the Attestation Build Tool.

### Github Action

Create `reproducible-build.yml` under the project's `.github/workflow` directory

```
name: Software Build Attestation

on:
  release:
    types: [published]

jobs:
  build:
    permissions: write-all
    runs-on: [self-hosted]
    steps:
    - name: Checkout
      uses: actions/checkout@v2
    - name: Build
      run: |
        attestable-build-tool build -output release.tar -nitro ~/ata-build-rust-latest.eif
    - name: Release
      uses: softprops/action-gh-release@v1
      with:
        files: release.tar
```

Create the `build.json` file in the project.
```
{
	"language": "rust",
	"input": {
		"cmd": "./scripts/build.sh",
		"vendor": "./scripts/vendor.sh"
	},
	"output": {
		"files": [
			"target/release/binary",
		]
	}
}
```

### Enclave Images

* [rust](https://attestation-build-image.s3.ap-southeast-1.amazonaws.com/ata-build-rust-latest.eif)
* [go](https://attestation-build-image.s3.ap-southeast-1.amazonaws.com/ata-build-go-latest.eif)

## Build Image

### Requirement

* Go
  * download golang [here](https://go.dev/dl/)
* Nitro Cli

### Usage

```
> cd image
> LANG=phala ./build-image.sh
* Build a docker image with the tag named: `ata-build-phala`
* Build a nitro enclave eif to `~/ata-build-phala-latest.eif`

```


## See also

* [Software Build Attestation](https://docs.ata.network/attestation-module/software-build-attestation)
* [Reproducible Build](https://docs.ata.network/research/reproducible-build)

## Contributing

**Before You Contribute**:
* **Raise an Issue**: If you find a bug or wish to suggest a feature, please open an issue first to discuss it. Detail the bug or feature so we understand your intention.  
* **Pull Requests (PR)**: Before submitting a PR, ensure:  
    * Your contribution successfully builds.
    * It includes tests, if applicable.

## License

MIT