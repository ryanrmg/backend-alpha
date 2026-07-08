package main

func main() {

	ctx := context.Background()

	srv, err := server.New(ctx, config.Load())

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		<-interruptSignal()

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
