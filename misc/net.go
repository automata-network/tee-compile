package misc

import (
	"net"
	"net/url"
	"strconv"

	"github.com/chzyer/logex"
	"github.com/mdlayher/vsock"
)

func Listen(uri *url.URL) (net.Listener, error) {
	switch uri.Scheme {
	case "unix", "tcp":
		ln, err := net.Listen(uri.Scheme, uri.Host)
		if err != nil {
			return nil, logex.Trace(err)
		}
		return ln, nil
	case "vsock":
		port, err := strconv.Atoi(uri.Port())
		if err != nil {
			return nil, logex.Trace(err)
		}
		return vsock.Listen(uint32(port), &vsock.Config{})
	default:
		return nil, logex.NewErrorf("unsupport uri: %v", uri)
	}
}
