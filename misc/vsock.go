package misc

import (
	"bytes"
	"context"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/chzyer/logex"
	"github.com/mdlayher/vsock"
)

type VsockLogWriter struct {
	client *http.Client
	url    string
}

func NewVsockLogWriter(url string) *VsockLogWriter {
	client := NewVsockClient(nil)
	return &VsockLogWriter{url: url, client: client}
}

func (w *VsockLogWriter) Write(data []byte) (int, error) {
	resp, err := w.client.Post(w.url, "application/octet-stream", bytes.NewReader(data))
	if err != nil {
		return 0, logex.Trace(err)
	}
	resp.Body.Close()
	return len(data), nil
}

func NewVsockClient(proxy func(*http.Request) (*url.URL, error)) *http.Client {
	return &http.Client{
		Transport: NewVsockTransport(proxy),
	}
}

func NewVsockTransport(proxy func(*http.Request) (*url.URL, error)) *http.Transport {
	if proxy == nil {
		proxy = http.ProxyFromEnvironment
	}
	return &http.Transport{
		Proxy: proxy,
		DialContext: defaultTransportDialContext(&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}),
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

func defaultTransportDialContext(dialer *net.Dialer) func(context.Context, string, string) (net.Conn, error) {
	return func(ctx context.Context, network string, addr string) (net.Conn, error) {
		sp := strings.Split(addr, ":")
		contextId, err := strconv.Atoi(sp[0])
		if err != nil {
			return nil, logex.Trace(err)
		}
		port, err := strconv.Atoi(sp[1])
		if err != nil {
			return nil, logex.Trace(err)
		}

		conn, err := vsock.Dial(uint32(contextId), uint32(port), nil)
		if err != nil {
			return nil, logex.Trace(err)
		}
		return conn, nil
	}
}
