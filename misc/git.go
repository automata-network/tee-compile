package misc

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/chzyer/logex"
)

type GitInfo struct {
	Commit string
}

func GetGitInfo(dir string) (*GitInfo, error) {
	var gitInfo GitInfo
	fp := filepath.Join(dir, ".git", "packed-refs")
	data, err := os.ReadFile(fp)
	if err != nil {
		fp := filepath.Join(dir, ".git", "HEAD")
		data, err := os.ReadFile(fp)
		if err != nil {
			return nil, logex.Trace(err)
		}
		gitInfo.Commit = string(data)
	} else {
		lines := strings.Split(string(data), "\n")
		lines = strings.Split(lines[1], " ")
		gitInfo.Commit = lines[0]
	}
	gitInfo.Commit = strings.TrimSpace(gitInfo.Commit)
	return &gitInfo, nil
}
