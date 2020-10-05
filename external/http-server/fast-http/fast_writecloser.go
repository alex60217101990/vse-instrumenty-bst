package fast_http

import "github.com/valyala/fasthttp"

type RequestWriteCloser struct {
	*fasthttp.RequestCtx
}

func (rwc *RequestWriteCloser) Close() error {
	// Noop
	return rwc.Close()
}
