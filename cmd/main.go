package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/david-tobi-peter/Prism/internal/handlers"
	"github.com/david-tobi-peter/Prism/internal/middlewares"
	"github.com/labstack/echo/v5"
)

func main() {
	e := echo.New()

	e.Use(middlewares.RequestTracer())

	proxyEngine := handlers.NewProxyEngine("http://localhost:8081")

	e.Any("/*", proxyEngine.HandleForward())

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
