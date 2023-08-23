package misc

import (
	"encoding/json"
	"errors"

	"github.com/chzyer/logex"
	"github.com/hf/nsm"
	"github.com/hf/nsm/request"
)

type AttestationReport struct {
	GitCommit  string
	Repo       string
	InputHash  string
	Image      string
	OutputHash string
	Nonce      string
}

func Attestation(report *AttestationReport) ([]byte, error) {
	sess, err := nsm.OpenDefaultSession()
	if err != nil {
		return nil, err
	}
	defer sess.Close()

	data, err := json.Marshal(report)
	if err != nil {
		return nil, logex.Trace(err)
	}

	res, err := sess.Send(&request.Attestation{
		UserData: data,
		Nonce:    []byte(report.Nonce),
	})
	if err != nil {
		return nil, err
	}

	if res.Error != "" {
		return nil, errors.New(string(res.Error))
	}

	if res.Attestation == nil || res.Attestation.Document == nil {
		return nil, errors.New("NSM device did not return an attestation")
	}

	return res.Attestation.Document, nil
}
