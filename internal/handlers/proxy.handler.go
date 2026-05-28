package handlers

import (
	"io"
	"log"
	"net"
	"net/http"

	"github.com/david-tobi-peter/Prism/internal/proxy"
	"github.com/david-tobi-peter/Prism/internal/validators"
	"github.com/labstack/echo/v5"
)

type ProxyEngine struct {
	Client *http.Client
	Router *proxy.Router
}

func NewProxyEngine(router *proxy.Router) *ProxyEngine {
	return &ProxyEngine{
		Client: proxy.PooledClient(),
		Router: router,
	}
}

func (pe *ProxyEngine) HandleForward() echo.HandlerFunc {
	return func(c *echo.Context) error {
		req := c.Request()
		res := c.Response()

		if err := validators.ValidateMessageBodyLength(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		validators.StripHopByHopHeaders(req)

		backend, rewrittenPath, found := pe.Router.Resolve(req.URL.Path)

		if !found {
			return echo.NewHTTPError(http.StatusNotFound, "No route found for "+req.URL.Path)
		}

		targetURL := backend + rewrittenPath
		if req.URL.RawQuery != "" {
			targetURL += "?" + req.URL.RawQuery
		}

		proxyReq, err := http.NewRequestWithContext(req.Context(), req.Method, targetURL, req.Body)
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

		resp, err := pe.Client.Do(proxyReq)
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
