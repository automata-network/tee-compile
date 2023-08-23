package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/automata-network/attestation-build-tool/misc"
	"github.com/chzyer/logex"
)

type BuildToolBuild struct {
	Dir    string `default:"."`
	Vendor bool
	Nitro  string
	Output string
}

func (b *BuildToolBuild) FlaglyHandle() error {
	if err := os.Chdir(b.Dir); err != nil {
		return logex.Trace(err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		return logex.Trace(err)
	}

	var vendorTars [][2]string

	if b.Vendor {
		vendorDir, err := os.MkdirTemp("", "vendor*")
		if err != nil {
			return logex.Trace(err)
		}
		defer os.RemoveAll(vendorDir)

		if err := misc.Exec("docker", "run", "--rm",
			"-v", fmt.Sprintf("%v:/tmp/vendor", vendorDir),
			"-v", fmt.Sprintf("%v:/workspace/code", cwd),
			"ata-build-rust", "/workspace/attestation-build-tool", "vendor", "-dir", "/workspace/code",
		); err != nil {
			return logex.Trace(err)
		}
		dir, err := os.ReadDir(vendorDir)
		if err != nil {
			return logex.Trace(err)
		}
		for _, vendor := range dir {
			vendorTars = append(vendorTars, [2]string{
				filepath.Join(vendorDir, vendor.Name()),
				"/root",
			})
		}
	}

	logex.Infof("pkg codes: %v", b.Dir)
	if err := os.Chdir(b.Dir); err != nil {
		return logex.Trace(err)
	}
	tarFile, err := misc.Tar("sourcecode", []string{"."})
	if err != nil {
		return logex.Trace(err)
	}
	defer func() {
		os.Remove(tarFile)
	}()

	tarFd, err := os.Open(tarFile)
	if err != nil {
		return logex.Trace(err)
	}
	logex.Info("package tar file to:", tarFd.Name())
	defer tarFd.Close()

	targetFile, err := os.OpenFile(b.Output, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		return logex.Trace(err)
	}
	defer targetFile.Close()

	var cmd *exec.Cmd
	var client *http.Client
	var endpoint string
	if b.Nitro != "" {
		misc.Exec("nitro-cli", "terminate-enclave", "--all")
		cmd = misc.RunNitroEnclave(b.Nitro)
		if err := cmd.Start(); err != nil {
			return logex.Trace(err)
		}
		client = misc.NewVsockClient(nil)
		endpoint = "http://11:12345"
	} else {
		// local mode
		cmd = exec.Command(os.Args[0], "worker", "-listen", "tcp://localhost:12345")
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		if err := cmd.Start(); err != nil {
			return logex.Trace(err)
		}
		client = http.DefaultClient
		endpoint = "http://localhost:12345"
	}

	for {
		_, err := client.Get(endpoint + "/ping")
		if err != nil {
			logex.Error(err)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	for _, tf := range vendorTars {
		fd, err := os.Open(tf[0])
		if err != nil {
			return logex.Trace(err)
		}
		query := url.Values{"target": {tf[1]}}
		response, err := client.Post(endpoint+"/vendor?"+query.Encode(), "application/octet-stream", fd)
		fd.Close()
		if err != nil {
			return logex.Trace(err)
		}
		if err := checkResponseError(response); err != nil {
			return logex.Trace(err)
		}
	}

	response, err := client.Post(endpoint+"/build", "application/octet-stream", tarFd)
	tarFd.Close()
	if err != nil {
		return logex.Trace(err)
	}

	defer cmd.Wait()
	defer response.Body.Close()
	if err := checkResponseError(response); err != nil {
		return logex.Trace(err)
	}

	if _, err := io.Copy(targetFile, response.Body); err != nil {
		return logex.Trace(err)
	}

	fmt.Println("report", response.Header.Get("Report"))
	fmt.Println("save file to:", targetFile.Name())

	return nil
}

func checkResponseError(response *http.Response) error {
	if response.StatusCode != 200 {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return logex.Trace(err)
		}
		return logex.NewErrorf("%s", body)
	}
	return nil
}

type BuildMode string

var (
	NitroBuildMode  BuildMode = "nitro"
	DockerBuildMode BuildMode = "docker"
	LocalBuildMode  BuildMode = "local"
)
