package di_container

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"
)

const serverShutdownTimeout = 10 * time.Second

func (c *cont) GetServer(ctx context.Context) *http.Server {
	return makeSingleton(&c.server, func() *http.Server {
		mux := http.NewServeMux()
		mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		mux.HandleFunc("POST /training/create", c.GetCreateTrainingHandler(ctx).Handle)

		server := &http.Server{
			Addr:              ":" + c.cfg.Ports.Http,
			Handler:           mux,
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       15 * time.Second,
			WriteTimeout:      30 * time.Second,
			IdleTimeout:       60 * time.Second,
		}

		c.addCancel(func() {
			shutdownContext, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
			defer cancel()

			if err := server.Shutdown(shutdownContext); err != nil {
				log.Printf("graceful HTTP shutdown failed: %v", err)
				if closeErr := server.Close(); closeErr != nil {
					log.Printf("close HTTP server: %v", closeErr)
				}
			}
		})

		return server
	})
}

func (c *cont) StartServer(ctx context.Context) <-chan error {
	serverErrors := make(chan error, 1)
	server := c.GetServer(ctx)

	go func() {
		serverErrors <- runServer(server)
	}()

	return serverErrors
}

func runServer(server *http.Server) error {
	certFile := os.Getenv("TLS_CERT_FILE")
	keyFile := os.Getenv("TLS_KEY_FILE")

	if certFile != "" || keyFile != "" {
		if certFile == "" || keyFile == "" {
			return errors.New("both TLS_CERT_FILE and TLS_KEY_FILE must be set")
		}
		log.Printf("HTTPS service is listening on %s", server.Addr)
		return server.ListenAndServeTLS(certFile, keyFile)
	}

	log.Printf("HTTP service is listening on %s", server.Addr)
	return server.ListenAndServe()
}
