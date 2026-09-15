package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"inktype-backend/internal/api"
	"inktype-backend/internal/config"
	"inktype-backend/internal/database"
	"inktype-backend/internal/repository"
	"inktype-backend/internal/storage"

	"github.com/clerk/clerk-sdk-go/v2"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.Load()

	clerk.SetKey(cfg.ClerkSecretKey)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Warn("failed to connect to database (is DATABASE_URL set?)", "error", err)
	} else {
		defer db.Pool.Close()
		logger.Info("successfully connected to Neon Postgres")
	}

	var repo repository.Querier
	if db != nil {
		repo = repository.New(db.Pool)
	}

	store, err := storage.NewClient(ctx, cfg.R2AccountID, cfg.R2AccessKeyID, cfg.R2SecretAccessKey, cfg.R2BucketName)
	if err != nil {
		logger.Warn("failed to initialize R2 storage client", "error", err)
	}

	server := api.NewServer(logger, repo, store)

	serverAddr := ":" + cfg.Port
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: server,
	}

	go func() {
		logger.Info("starting InkType API server", "addr", serverAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
	}

	logger.Info("server exiting gracefully")
}
