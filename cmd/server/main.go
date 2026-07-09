package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ryanrmg/backend-alpha/internal/config"
	"github.com/ryanrmg/backend-alpha/internal/server"
)

func main() {

	ctx := context.Background()

	srv, err := server.New(ctx, config.Load())

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cancel()

		srv.Shutdown(shutdownCtx)
	}()

	if err := srv.Start(); err != nil &&
		err != http.ErrServerClosed {

		log.Fatal(err)
	}
}
