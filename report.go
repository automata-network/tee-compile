package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/chzyer/logex"
	"github.com/hf/nitrite"
)

type BuildToolReport struct {
	File string `type:"[0]"`
}

func (r *BuildToolReport) FlaglyHandle() error {
	reportBytes, err := os.ReadFile(r.File)
	if err != nil {
		return logex.Trace(err, r.File)
	}
	report, err := nitrite.Verify(reportBytes, nitrite.VerifyOptions{
		CurrentTime: time.Now(),
	})
	if err != nil {
		return logex.Trace(err)
	}
	fmt.Printf("PCR0: 0x%v\n", hex.EncodeToString(report.Document.PCRs[0]))
	fmt.Printf("%s\n", report.Document.UserData)
	return nil
}
