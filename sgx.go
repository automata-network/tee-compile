package main

import (
	"fmt"

	"github.com/automata-network/tee-compile/misc"
	"github.com/chzyer/logex"
)

type BuildToolSGX struct {
	MrEnclave *BuildToolSGXMrEnclave `flagly:"handler"`
}

type BuildToolSGXMrEnclave struct {
	File string `type:"[0]"`
}

func (h *BuildToolSGXMrEnclave) FlaglyHandle() error {
	mrenclave, err := misc.GetMrEnclave(h.File)
	if err != nil {
		return logex.Trace(err)
	}
	fmt.Printf("0x%x\n", mrenclave)
	return nil
}
