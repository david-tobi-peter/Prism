package handlers

import (
	"io"
	"log"
	"net"
	"net/http"

	"github.com/david-tobi-peter/Prism/internal/validators"
	"github.com/labstack/echo/v5"
)

func HandleForward(downStreamServer string) echo.HandlerFunc {
	return func(c *echo.Context) error {
		req := c.Request()
		res := c.Response()

		if err := validators.ValidateMessageBodyLength(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		backendURL := downStreamServer + req.URL.Path

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

		clientIP := c.RealIP()
		host, _, err := net.SplitHostPort(req.RemoteAddr)
		if err == nil && clientIP == "" {
			clientIP = host
		}

		priorXFF := req.Header.Get("X-Forwarded-For")
		if priorXFF != "" {
			proxyReq.Header.Set("X-Forwarded-For", priorXFF+", "+clientIP)
		} else {
			proxyReq.Header.Set("X-Forwarded-For", clientIP)
		}

		proxyReq.Header.Set("Test-Header", "Test")

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
	}
}
