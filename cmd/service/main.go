package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"training_planner/internal/di_container"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	c := di_container.New()
	defer c.Close()

	serverErrors := c.StartServer(context.Background())

	shutdownSignals := make(chan os.Signal, 1)
	signal.Notify(shutdownSignals, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(shutdownSignals)

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("run HTTP server: %w", err)
		}
	case <-shutdownSignals:
	}

	return nil
}
