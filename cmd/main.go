package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/david-tobi-peter/Prism/internal/config"
	"github.com/david-tobi-peter/Prism/internal/handlers"
	"github.com/david-tobi-peter/Prism/internal/middlewares"
	"github.com/david-tobi-peter/Prism/internal/proxy"
	"github.com/labstack/echo/v5"
)

func main() {
	configFlag := flag.String("config", "", "Path to the required configuration TOML file")
	flag.Parse()

	if *configFlag == "" {
		log.Fatalf("Missing required routing configuration file. Usage: ./prism -config=path/to/routes.toml")
	}

	if _, err := os.Stat(*configFlag); os.IsNotExist(err) {
		log.Fatalf("Specified configuration file does not exist at path: %s", *configFlag)
	}

	cfg, err := config.LoadConfig(*configFlag)
	if err != nil {
		log.Fatalf("Failed to parse configuration file cleanly: %v", err)
	}

	router := proxy.NewRouter(cfg)
	proxyEngine := handlers.NewProxyEngine(router)

	e := echo.New()
	e.Use(middlewares.RequestTracer())
	e.Any("/*", proxyEngine.HandleForward())

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	sc := echo.StartConfig{
		Address:         fmt.Sprintf(":%d", cfg.Port),
		GracefulTimeout: 5 * time.Second,
	}

	if err := sc.Start(ctx, e); err != nil {
		e.Logger.Error("Prism failed to start", "error", err)
	}
}
