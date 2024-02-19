package misc

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/chzyer/logex"
)

type LogOutput struct {
	Stdout io.Writer
	Stderr io.Writer
}

func Exec(out *LogOutput, name string, args ...string) error {
	if out == nil {
		out = &LogOutput{}
	}
	if out.Stdout == nil {
		out.Stdout = os.Stdout
	}
	if out.Stderr == nil {
		out.Stderr = os.Stderr
	}
	fmt.Fprintf(out.Stdout, "exec %q\n", strings.Join(append([]string{name}, args...), " "))

	cmd := exec.Command(name, args...)
	cmd.Stderr = out.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = out.Stdout
	if err := cmd.Run(); err != nil {
		return logex.Trace(err)
	}
	exitCode := cmd.ProcessState.ExitCode()
	if exitCode != 0 {
		return logex.NewErrorf("exit by code: %v", exitCode)
	}
	return nil
}

func Tar(out *LogOutput, prefix string, filelist []string) (string, error) {
	return TarTo(out, "", prefix, filelist)
}

func TarTo(out *LogOutput, dir, tag string, filelist []string) (string, error) {
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
	if err := Exec(out, "tar", args...); err != nil {
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

func RunNitroEnclave(path string, mem string, cpu uint, cid uint32, debug bool) *exec.Cmd {
	args := []string{
		"run-enclave",
		"--cpu-count", fmt.Sprint(cpu),
		"--memory", mem,
		"--enclave-cid", fmt.Sprint(cid),
		"--eif-path", path,
	}
	if debug {
		args = append(args,
			"--debug-mode",
			"--attach-console",
		)
	}
	cmd := exec.Command("nitro-cli", args...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd
}
