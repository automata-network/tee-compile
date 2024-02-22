package misc

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

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

func RunNitroEnclave(path string, mem string, cpu uint, cid uint32, debug bool) (*Process, error) {
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
	proc := NewProcess(context.TODO(), "nitro-cli", args...)
	if err := proc.Start(); err != nil {
		return nil, logex.Trace(err)
	}
	return proc, nil
}

type Process struct {
	cmd *exec.Cmd
	ctx context.Context

	errorMsg *FixedBuffer

	mutex     sync.Mutex
	waitError error
	waitCh    chan struct{}
}

func NewProcess(ctx context.Context, name string, args ...string) *Process {
	cmd := exec.CommandContext(ctx, name, args...)
	errorMsg := NewFixedBuffer(1024)

	cmd.Stderr = io.MultiWriter(os.Stderr, errorMsg)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return &Process{
		ctx:       ctx,
		cmd:       cmd,
		errorMsg:  errorMsg,
		waitError: nil,
		waitCh:    make(chan struct{}),
	}
}

func (p *Process) Start() error {
	if err := p.cmd.Start(); err != nil {
		return logex.Trace(err)
	}

	go func() {
		if err := p.cmd.Wait(); err != nil {
			p.mutex.Lock()
			p.waitError = err
			p.mutex.Unlock()
		}
		close(p.waitCh)
	}()
	return nil
}

func (p *Process) Done() <-chan struct{} {
	return p.waitCh
}

func (p *Process) ErrorMsg() string {
	return p.errorMsg.String()
}

func (p *Process) Exited() bool {
	if p.cmd.ProcessState == nil {
		return false
	}
	return p.cmd.ProcessState.Exited()
}

func (p *Process) ExitCode() int {
	if p.Exited() {
		return p.cmd.ProcessState.ExitCode()
	}
	return 0
}

func (p *Process) Wait() error {
	p.mutex.Lock()
	if p.waitError != nil {
		err := p.waitError
		p.mutex.Unlock()
		return logex.Trace(err)
	}
	p.mutex.Unlock()

	select {
	case <-p.waitCh:
		p.mutex.Lock()
		defer p.mutex.Unlock()
		if p.waitError == nil {
			return nil
		}
		return logex.Trace(p.waitError)
	case <-p.ctx.Done():
		return logex.Trace(p.ctx.Err())
	}
}

type FixedBuffer struct {
	data []byte
}

func NewFixedBuffer(cap int) *FixedBuffer {
	return &FixedBuffer{data: make([]byte, 0, cap)}
}

func (f *FixedBuffer) Available() int {
	return cap(f.data) - len(f.data)
}

func (f *FixedBuffer) Write(data []byte) (int, error) {
	n := len(data)
	data = data[:f.Available()]
	f.data = append(f.data, data...)
	return n, nil
}

func (f *FixedBuffer) String() string {
	return string(f.data)
}
