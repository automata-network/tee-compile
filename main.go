package main

import (
	"os"

	"github.com/chzyer/flagly"
	"github.com/chzyer/logex"
)

type BuildTool struct {
	Build  *BuildToolBuild  `flagly:"handler"`
	Worker *BuildToolWorker `flagly:"handler"`
	Vendor *BuildToolVendor `flagly:"handler"`
	SGX    *BuildToolSGX    `flagly:"handler"`
	Report *BuildToolReport `flagly:"handler"`
}

func main() {
	fset := flagly.New(os.Args[0])
	app := BuildTool{}
	if err := fset.Compile(app); err != nil {
		logex.Fatal(err)
	}
	if err := fset.Run(os.Args[1:]); err != nil {
		logex.Fatal(err)
	}
}
