package rest_app

import (
	"bufio"
	"context"
	"net"
	"net/http"
)

type Server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type ResponseWriter interface {
	Header() http.Header
	Write([]byte) (int, error)
	WriteHeader(statusCode int)
}

type Flusher interface {
	Flush()
}

type Hijacker interface {
	Hijack() (net.Conn, *bufio.ReadWriter, error)
}

type CloseNotifier interface {
	CloseNotify() <-chan bool
}

type closeWriter interface {
	CloseWrite() error
}
