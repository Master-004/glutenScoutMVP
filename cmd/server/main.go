// Command server runs the GlutenScout web app.
package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Master-004/glutenscout/internal/handlers"
	"github.com/Master-004/glutenscout/internal/store"
	"github.com/Master-004/glutenscout/web"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	if err := run(logger); err != nil {
		logger.Error("server stopped", "err", err)
		os.Exit(1)
	}
}

// run holds the real startup logic. Returning an error (instead of calling
// log.Fatal everywhere) lets deferred cleanup like st.Close() actually run.
func run(logger *slog.Logger) error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return errors.New("DATABASE_URL is not set")
	}

	// ctx is cancelled when you press Ctrl+C or Docker stops the container,
	// which triggers the graceful shutdown below.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	st, err := store.New(ctx, dbURL)
	if err != nil {
		return err
	}
	defer st.Close()

	tmpl, err := handlers.LoadTemplates(web.FS)
	if err != nil {
		return err
	}
	static, err := fs.Sub(web.FS, "static")
	if err != nil {
		return err
	}

	h := handlers.New(st, tmpl, static, logger)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: h.Routes(),
		// Timeouts stop slow or stuck clients from tying up the server forever.
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Serve in the background so this goroutine can wait for a shutdown signal.
	errCh := make(chan error, 1)
	go func() {
		logger.Info("listening", "addr", "http://localhost:"+port)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("listen: %w", err)
	case <-ctx.Done():
	}

	// Give in-flight requests a few seconds to finish before exiting.
	logger.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
