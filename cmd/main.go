package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-monitor/internal/config"
	"go-monitor/internal/logger"
	"go-monitor/internal/rest"
)

func main() {
	// Initialize configuration
	c, err := config.NewConfig()
	if err != nil {
		log.Fatalf("error during config initialization: %s", err)
	}

	// Initialize logger
	l, err := logger.NewLogger(c.LokiClientURL)
	if err != nil {
		log.Fatalf("error during logger initialization: %s", err)
	}
	defer l.Sync() // Ensure log messages are flushed before exiting
	l.Info("Logger initialized")

	// Initialize REST controller
	r := rest.NewREST(l, c.Version)
	l.Info("REST controller initialized")

	// Start REST server in a goroutine
	go r.Listen()
	l.Info("REST server started and accepting connections on port :8080")

	// Handle OS signals to gracefully shutdown
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signalCh)

	// Wait for signals
	sig := <-signalCh
	l.Sugar().Infof("Received signal %s, exiting...", sig)
}
