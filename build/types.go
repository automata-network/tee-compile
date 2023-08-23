package build

import (
	"encoding/json"
	"os"

	"github.com/chzyer/logex"
)

type Manifest struct {
	Language string          `json:"language"`
	Input    *ManifestInput  `json:"input"`
	Output   *ManifestOutput `json:"output"`
}

type ManifestInput struct {
	Cmd    string   `json:"cmd"`
	Vendor string   `json:"vendor"`
	Env    []string `json:"env"`
}

type ManifestOutput struct {
	Files []string `json:"files"`
}

func NewManifest(file string) (*Manifest, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, logex.Trace(err, file)
	}
	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, logex.Trace(err)
	}
	return &manifest, nil
}
