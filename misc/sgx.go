package misc

import (
	"debug/elf"
	"fmt"

	"github.com/chzyer/logex"
)

func GetMrEnclave(file string) ([]byte, error) {
	f, err := elf.Open(file)
	if err != nil {
		return nil, logex.Trace(err)
	}
	defer f.Close()

	for _, sect := range f.Sections {
		if sect.Name == ".note.sgxmeta" {
			data, err := sect.Data()
			if err != nil {
				return nil, logex.Trace(err)
			}
			var hash [32]byte
			copy(hash[:], data[1049:1049+32])
			return hash[:], nil
		}
	}
	return nil, fmt.Errorf("sgxmeta not found")
}
