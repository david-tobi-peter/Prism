package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v5"
)

func main() {
	e := echo.New()

	e.Any("/*", func(c *echo.Context) error {
		req := c.Request()
		res := c.Response()

		backendURL := "http://localhost:8081" + req.URL.Path

		if req.URL.RawQuery != "" {
			backendURL += "?" + req.URL.RawQuery
		}

		proxyReq, err := http.NewRequestWithContext(req.Context(), req.Method, backendURL, req.Body)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal Proxy Error")
		}

		for key, values := range req.Header {
			for _, value := range values {
				proxyReq.Header.Add(key, value)
			}
		}

		resp, err := http.DefaultClient.Do(proxyReq)
		if err != nil {
			log.Printf("Backend connection faiure: %v", err)
			return echo.NewHTTPError(http.StatusBadGateway, "Bad Gateway Error")
		}
		defer resp.Body.Close()

		for key, values := range resp.Header {
			for _, value := range values {
				res.Header().Add(key, value)
			}
		}

		res.WriteHeader(resp.StatusCode)

		_, err = io.Copy(res, resp.Body)
		if err != nil {
			log.Printf("Error copying response body: %v", err)
		}

		return nil
	})

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	sc := echo.StartConfig{
		Address:         ":8080",
		GracefulTimeout: 5 * time.Second,
	}

	if err := sc.Start(ctx, e); err != nil {
		e.Logger.Error("Prism failed to start", "error", err)
	}
}
