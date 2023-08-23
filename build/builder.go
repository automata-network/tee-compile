package build

import (
	"os"

	"github.com/automata-network/attestation-build-tool/misc"
	"github.com/chzyer/logex"
)

type Builder struct {
	Manifest *Manifest
	Result   *misc.MerkleTreeResult
}

func NewBuilder(manifest *Manifest) *Builder {
	return &Builder{Manifest: manifest}
}

func (b *Builder) Vendor() error {
	vendorExecutor, ok := VendorList[b.Manifest.Language]
	if !ok {
		return logex.NewErrorf("unknown language=%q for vendor", b.Manifest.Language)
	}

	if err := b.build(b.Manifest.Input.Vendor); err != nil {
		return logex.Trace(err)
	}

	if err := vendorExecutor.Vendor(); err != nil {
		return logex.Trace(err)
	}
	return nil
}

func (b *Builder) build(vendor string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return logex.Trace(err)
	}
	logex.Infof("cwd: %s", cwd)
	cmd := b.Manifest.Input.Cmd
	if vendor != "" {
		cmd = vendor
	}
	logex.Infof("running cmd: %q", cmd)
	if err := misc.Exec("bash", "-c", cmd); err != nil {
		return logex.Trace(err)
	}
	return nil
}

func (b *Builder) Build() error {
	if err := b.build(""); err != nil {
		return logex.Trace(err)
	}

	result, err := misc.FilesMerkleTree(b.Manifest.Output.Files, 10, nil)
	if err != nil {
		return logex.Trace(err)
	}

	b.Result = result
	return nil
}

func (b *Builder) Tar() (string, error) {
	tarFile, err := misc.Tar("output", b.Result.FileList)
	if err != nil {
		return "", logex.Trace(err)
	}
	return tarFile, nil
}
