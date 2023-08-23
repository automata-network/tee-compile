package misc

import (
	"crypto/rand"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/chzyer/logex"
)

func Exec(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return logex.Trace(err)
	}
	exitCode := cmd.ProcessState.ExitCode()
	if exitCode != 0 {
		return logex.NewErrorf("exit by code: %v", exitCode)
	}
	return nil
}

func Tar(prefix string, filelist []string) (string, error) {
	return TarTo("", prefix, filelist)
}

func TarTo(dir, tag string, filelist []string) (string, error) {
	var buf [16]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return "", logex.Trace(err)
	}
	if dir == "" {
		dir = os.TempDir()
	}
	fp := filepath.Join(dir, fmt.Sprintf("%v-%x.tar", tag, buf))
	args := []string{"cf", fp}
	args = append(args, filelist...)
	if err := Exec("tar", args...); err != nil {
		os.RemoveAll(fp)
		return "", logex.Trace(err)
	}
	return fp, nil
}

func InDir(dir string, run func() error) error {
	cwd, err := os.Getwd()
	if err != nil {
		return logex.Trace(err)
	}
	if err := os.Chdir(dir); err != nil {
		return logex.Trace(err)
	}
	if err := run(); err != nil {
		return logex.Trace(err)
	}
	os.Chdir(cwd)
	return nil
}

func RunNitroEnclave(path string) *exec.Cmd {
	args := []string{
		"run-enclave",
		"--cpu-count", "2",
		"--memory", "12288",
		"--enclave-cid", "11",
		"--eif-path", path,
		"--debug-mode",
		"--attach-console",
	}
	cmd := exec.Command("nitro-cli", args...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd
}
