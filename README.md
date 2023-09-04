# Attestable Build Tool

## Github Runner
1. Apply for a nitro enclave machine on AWS.
2. Configure the GitHub runner:  
    2.1. Settings → Actions → Runners → New self-hosted runner  
    2.2. Follow the instructions to configure the GitHub Runner  
3. Download the Software Build Attestation Image.
4. Download the Attestation Build Tool.

## Github Action

Create `build_attestation.yml` under the project's `.github/workflow` directory

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
        attestation-build-tool build -output release.tar -nitro ~/ata-build-rust-latest.eif
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

## Enclave Images

* [rust](https://attestation-build-image.s3.ap-southeast-1.amazonaws.com/ata-build-rust-latest.eif)
* [go](https://attestation-build-image.s3.ap-southeast-1.amazonaws.com/ata-build-go-latest.eif)