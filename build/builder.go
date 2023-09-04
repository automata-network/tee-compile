package build

import (
	"os"

	"github.com/automata-network/attestable-build-tool/misc"
	"github.com/chzyer/logex"
)

type Builder struct {
	Manifest     *Manifest
	Nonce        string
	GitInfo      *misc.GitInfo
	OutputResult *misc.MerkleTreeResult
	InputResult  *misc.MerkleTreeResult
	logOutput    *misc.LogOutput
}

func NewBuilder(manifest *Manifest, nonce string, logOutput *misc.LogOutput) *Builder {
	return &Builder{Manifest: manifest, Nonce: nonce, logOutput: logOutput}
}

func (b *Builder) Vendor() error {
	vendorExecutor, ok := VendorList[b.Manifest.Language]
	if !ok {
		return logex.NewErrorf("unknown language=%q for vendor", b.Manifest.Language)
	}

	if err := b.build(b.Manifest.Input.Vendor); err != nil {
		return logex.Trace(err)
	}

	if err := vendorExecutor.Vendor(b.logOutput); err != nil {
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
	if err := misc.Exec(b.logOutput, "bash", "-c", cmd); err != nil {
		return logex.Trace(err)
	}
	return nil
}

func (b *Builder) Build() error {
	gitInfo, err := misc.GetGitInfo(".")
	if err != nil {
		return logex.Trace(err)
	}
	inputResult, err := misc.FilesMerkleTree([]string{"."}, 10, nil)
	if err != nil {
		return logex.Trace(err)
	}

	if err := b.build(""); err != nil {
		return logex.Trace(err)
	}

	outputResult, err := misc.FilesMerkleTree(b.Manifest.Output.Files, 10, nil)
	if err != nil {
		return logex.Trace(err)
	}

	b.GitInfo = gitInfo
	b.InputResult = inputResult
	b.OutputResult = outputResult
	return nil
}

func (b *Builder) Tar() (string, error) {
	tarFile, err := misc.Tar(b.logOutput, "output", b.OutputResult.FileList)
	if err != nil {
		return "", logex.Trace(err)
	}
	return tarFile, nil
}
