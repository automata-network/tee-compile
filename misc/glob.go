package misc

import (
	"path/filepath"
	"sort"

	"github.com/chzyer/logex"
)

func GlobSortList(patterns []string) ([]string, error) {
	var allList []string
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, logex.Trace(err, pattern)
		}
		allList = append(allList, matches...)
	}
	sort.Slice(allList, func(i, j int) bool {
		return allList[i] < allList[j]
	})
	return allList, nil
}
