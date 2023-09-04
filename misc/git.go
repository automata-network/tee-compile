package misc

import (
	"os"
	"path/filepath"

	"github.com/chzyer/logex"
)

type GitInfo struct {
	Commit string
}

func GetGitInfo(dir string) (*GitInfo, error) {
	fp := filepath.Join(dir, ".git", "HEAD")
	data, err := os.ReadFile(fp)
	if err != nil {
		return nil, logex.Trace(err)
	}
	return &GitInfo{Commit: string(data)}, nil
}
