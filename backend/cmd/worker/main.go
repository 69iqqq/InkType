package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"inktype-backend/internal/ai"
	"inktype-backend/internal/config"
	"inktype-backend/internal/logger"
	"inktype-backend/internal/ocr"
	"inktype-backend/internal/pdf"
	"inktype-backend/internal/repository"
	"inktype-backend/internal/server"
	"inktype-backend/internal/storage"
	"inktype-backend/internal/worker"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	loggerService := logger.NewLoggerService(cfg.Observability)
	defer loggerService.Shutdown()

	log := logger.NewLoggerWithService(cfg.Observability, loggerService)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv, err := server.New(cfg, &log, loggerService)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize server")
	}

	repo := repository.New(srv.DB.Pool)

	store, err := storage.NewClient(
		ctx,
		cfg.Storage.R2AccountID,
		cfg.Storage.R2AccessKeyID,
		cfg.Storage.R2SecretAccessKey,
		cfg.Storage.R2BucketName,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize storage client")
	}

	pdfRender := pdf.NewSystemRenderer(store)
	visionProvider, err := ai.NewGenAIVisionProvider(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize GenAI Vision Provider")
	}
	ocrProvider := ocr.NewStubOCRProvider()

	w := worker.NewWorker(cfg, srv.DB.Pool, repo, store, pdfRender, visionProvider, ocrProvider, log)

	if err := w.Start(ctx); err != nil {
		log.Fatal().Err(err).Msg("Failed to start worker")
	}

	// Wait for interrupt signal to gracefully shutdown the worker
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down worker...")
	cancel() // Cancel the context to stop workers gracefully
	w.Stop()

	// Shutdown resources
	if srv.DB != nil {
		srv.DB.Close()
	}
	if srv.Redis != nil {
		srv.Redis.Close()
	}

	log.Info().Msg("Worker stopped")
}
