package build

import (
	"github.com/automata-network/attestation-build-tool/misc"
	"github.com/chzyer/logex"
)

type VendorExecutor interface {
	Vendor() error
}

var VendorList = map[string]VendorExecutor{
	"rust": &RustVendor{},
}

type RustVendor struct{}

func (r *RustVendor) Vendor() error {
	if err := misc.InDir("/root", func() error {
		fp, err := misc.TarTo("/tmp/vendor/", "vendor", []string{".cargo/registry", ".cargo/git"})
		if err != nil {
			return logex.Trace(err)
		}
		logex.Infof("vendor to %v", fp)
		return nil
	}); err != nil {
		return logex.Trace(err)
	}
	return nil
}
