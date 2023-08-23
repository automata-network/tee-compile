package misc

import (
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
			logex.Error("vsock dial fail:", err)
			return nil, logex.Trace(err)
		}
		return conn, nil
	}
}
