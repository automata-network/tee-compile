package build

import (
	"encoding/json"
	"os"
	"strings"

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
	SgxSignedSo string     `json:"sgx_signed_so"`
	ExtraHash   *ExtraHash `json:"extra_hash"`
	Files       []string   `json:"files"`
}

type ExtraHash struct {
	JsonFile  string `json:"json_file"`
	JsonField string `json:"json_field"`
}

func (eh *ExtraHash) Get() (string, error) {
	if eh == nil {
		return "", logex.NewErrorf("no extra hash")
	}
	data, err := os.ReadFile(eh.JsonFile)
	if err != nil {
		return "", logex.Trace(err)
	}
	var userData map[string]interface{}
	if err := json.Unmarshal(data, &userData); err != nil {
		return "", logex.Trace(err)
	}
	fieldPath := strings.Split(eh.JsonField, ".")
	for idx, field := range fieldPath {
		if idx == len(fieldPath)-1 {
			return userData[field].(string), nil
		}
		userData = userData[field].(map[string]interface{})
	}
	return "", logex.NewErrorf("hash not found")
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
