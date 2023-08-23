package misc

import (
	"io"
	"os"
	"sync"

	"github.com/chzyer/logex"
	"github.com/wealdtech/go-merkletree"
	"github.com/wealdtech/go-merkletree/keccak256"
	"golang.org/x/crypto/sha3"
)

func GetFileHash(fp string) ([]byte, error) {
	fd, err := os.Open(fp)
	if err != nil {
		return nil, logex.Trace(err)
	}
	defer fd.Close()

	fi, err := fd.Stat()
	if err != nil {
		return nil, logex.Trace(err)
	}

	hash := sha3.NewLegacyKeccak256()
	if _, err := hash.Write([]byte(fp)); err != nil {
		return nil, logex.Trace(err)
	}

	if !fi.IsDir() {
		n, err := io.Copy(hash, fd)
		if err != nil {
			return nil, logex.Trace(err, fp)
		}
		if n != fi.Size() {
			return nil, logex.NewErrorf("size mismatch")
		}
	}
	data := hash.Sum(nil)
	return data, nil
}

type MerkleTreeResult struct {
	Tree     *merkletree.MerkleTree
	Root     []byte
	FileList []string
}

func FilesMerkleTree(patterns []string, workers int, salt []byte) (*MerkleTreeResult, error) {
	type Task struct {
		FilePath string
		Idx      int
	}
	fileList, err := GlobSortList(patterns)
	if err != nil {
		return nil, logex.Trace(err)
	}
	output := make([][]byte, len(fileList))
	ch := make(chan *Task, workers)
	errs := make(chan error, 1)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var err error
			for task := range ch {
				output[task.Idx], err = GetFileHash(task.FilePath)
				if err != nil {
					errs <- logex.Trace(err)
					break
				}
			}
		}()
	}
	go func() {
		for idx, item := range fileList {
			ch <- &Task{Idx: idx, FilePath: item}
		}
		close(ch)
	}()
	wg.Wait()

	select {
	case err := <-errs:
		return nil, logex.Trace(err)
	default:
	}

	tree, err := merkletree.NewUsing(output, keccak256.New(), salt)
	if err != nil {
		return nil, logex.Trace(err)
	}
	return &MerkleTreeResult{
		Tree:     tree,
		Root:     tree.Root(),
		FileList: fileList,
	}, nil
}
