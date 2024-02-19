package build

import (
	"github.com/automata-network/attestable-build-tool/misc"
	"github.com/chzyer/logex"
)

type VendorExecutor interface {
	Vendor(log *misc.LogOutput) error
	VendorPath() string
}

var VendorList = map[string]VendorExecutor{
	"rust":  &RustVendor{},
	"phala": &PhalaVendor{},
}

type RustVendor struct{}

func (r *RustVendor) Vendor(log *misc.LogOutput) error {
	if err := misc.InDir("/root", func() error {
		fp, err := misc.TarTo(log, "/tmp/vendor/", "vendor", []string{".cargo/registry", ".cargo/git", ".cargo/bin"})
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

func (r *RustVendor) VendorPath() string {
	return "/root"
}

type PhalaVendor struct{}

func (r *PhalaVendor) Vendor(log *misc.LogOutput) error {
	if err := misc.InDir("/usr/local/", func() error {
		fp, err := misc.TarTo(log, "/tmp/vendor/", "vendor", []string{"cargo", "rustup"})
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

func (r *PhalaVendor) VendorPath() string {
	return "/usr/local/"
}
