package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/automata-network/attestation-build-tool/build"
	"github.com/automata-network/attestation-build-tool/misc"
	"github.com/chzyer/logex"
)

type BuildToolWorker struct {
	Listen string `desc:"unix://"`
	Dir    string `default:"."`

	Server *http.Server `flagly:"-"`
}

func (b *BuildToolWorker) FlaglyHandle() error {
	if err := os.Chdir(b.Dir); err != nil {
		return logex.Trace(err)
	}
	uri, err := url.Parse(b.Listen)
	if err != nil {
		return logex.Trace(err)
	}
	listener, err := misc.Listen(uri)
	if err != nil {
		return logex.Trace(err)
	}
	defer listener.Close()

	b.Server = &http.Server{
		Addr:    uri.Host,
		Handler: b,
	}
	if err := b.Server.Serve(listener); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return logex.Trace(err)
		}
	}
	return nil
}

type BuildResult struct {
	Report string
	Body   io.ReadCloser
}

func (b *BuildToolWorker) Vendor(target string, data io.Reader) error {
	if err := os.MkdirAll(target, 0755); err != nil {
		return logex.Trace(err)
	}
	fd, err := os.CreateTemp("", "worker-vendor-tar-*")
	if err != nil {
		return logex.Trace(err)
	}
	defer fd.Close()
	defer os.Remove(fd.Name())

	if _, err := io.Copy(fd, data); err != nil {
		return logex.Trace(err)
	}

	if err := misc.InDir(target, func() error {
		if err := misc.Exec("tar", "xvf", fd.Name()); err != nil {
			return logex.Trace(err)
		}
		return nil
	}); err != nil {
		return logex.Trace(err)
	}
	return nil
}

func (b *BuildToolWorker) TestSpace() error {

	buf := make([]byte, 1<<20)
	for i := 0; i < 2; i++ {
		fd, err := os.Create("/workspace/test")
		if err != nil {
			return logex.Trace(err)
		}
		defer fd.Close()

		total := 0
		for {
			n, err := fd.Write(buf)
			if err != nil {
				logex.Infof("[%v] total write: %v MB", i, total/1024/1024)
				logex.Errorf("write fail: %v", err)
				break
			}
			total += n
		}
		fd.Close()
		os.RemoveAll(fd.Name())
	}
	return nil
}

func (b *BuildToolWorker) Build(data io.Reader) (*BuildResult, error) {
	fd, err := os.CreateTemp("", "worker-builder-tar-*")
	if err != nil {
		return nil, logex.Trace(err)
	}
	defer fd.Close()
	defer os.Remove(fd.Name())

	if _, err := io.Copy(fd, data); err != nil {
		return nil, logex.Trace(err)
	}

	if err := misc.Exec("tar", "xvf", fd.Name()); err != nil {
		return nil, logex.Trace(err)
	}

	os.Remove(fd.Name())

	manifest, err := build.NewManifest("build.json")
	if err != nil {
		return nil, logex.Trace(err)
	}

	builder := build.NewBuilder(manifest)
	if err := builder.Build(); err != nil {
		return nil, logex.Trace(err)
	}

	output, err := builder.Tar()
	if err != nil {
		return nil, logex.Trace(err)
	}
	outputFd, err := os.Open(output)
	if err != nil {
		return nil, logex.Trace(err)
	}

	logex.Infof("hash: %x", builder.Result.Root)

	return &BuildResult{
		Report: "hello",
		Body:   outputFd,
	}, nil
}

func (b *BuildToolWorker) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	switch req.URL.Path {
	case "/build":
		misc.Exec("du", "-d1", "-h", "/root/.cargo")
		report, err := b.Build(req.Body)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprint(w, err.Error())
		} else {
			w.WriteHeader(200)
			w.Header().Set("Report", report.Report)
			io.Copy(w, report.Body)
			report.Body.Close()
		}

		go func() {
			// TODO: remove file
			time.Sleep(time.Second)
			b.Server.Shutdown(context.TODO())
		}()
	case "/testspace":
		b.TestSpace()
	case "/vendor":
		target := req.URL.Query().Get("target")
		if err := b.Vendor(target, req.Body); err != nil {
			w.WriteHeader(400)
			fmt.Fprint(w, err.Error())
		} else {
			w.WriteHeader(200)
		}
	}

}
