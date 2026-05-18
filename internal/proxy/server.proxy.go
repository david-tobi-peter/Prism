package proxy

import (
	"net/http"
	"time"
)

func ServerConfiguration(handler http.Handler, addr string) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: handler,

		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       60 * time.Second,

		MaxHeaderBytes: 1024 * 1024,
	}
}
