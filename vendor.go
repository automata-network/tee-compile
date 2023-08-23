package main

import (
	"os"

	"github.com/automata-network/attestation-build-tool/build"
	"github.com/chzyer/logex"
)

type BuildToolVendor struct {
	Dir string `default:"."`
}

func (b *BuildToolVendor) FlaglyHandle() error {
	if err := os.Chdir(b.Dir); err != nil {
		return logex.Trace(err)
	}

	manifest, err := build.NewManifest("build.json")
	if err != nil {
		return logex.Trace(err)
	}
	builder := build.NewBuilder(manifest)
	if err := builder.Vendor(); err != nil {
		return logex.Trace(err)
	}
	return nil
}
