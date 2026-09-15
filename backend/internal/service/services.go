package service

import (
	"context"
	"time"

	"inktype-backend/internal/lib/job"
	"inktype-backend/internal/repository"
	"inktype-backend/internal/server"
	"inktype-backend/internal/storage"
)

type Services struct {
	Auth         *AuthService
	Job          *job.JobService
	DocumentRepo repository.Querier
	Storage      *storage.Client
}

func NewServices(s *server.Server, repos *repository.Repositories) (*Services, error) {
	authService := NewAuthService(s)

	// Initialize repository
	var repo repository.Querier
	if s.DB != nil {
		repo = repository.New(s.DB.Pool)
	}

	// Initialize storage
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var store *storage.Client
	var err error
	if s.Config.Storage.R2AccountID != "" {
		store, err = storage.NewClient(
			ctx,
			s.Config.Storage.R2AccountID,
			s.Config.Storage.R2AccessKeyID,
			s.Config.Storage.R2SecretAccessKey,
			s.Config.Storage.R2BucketName,
		)
		if err != nil {
			s.Logger.Warn().Err(err).Msg("failed to initialize R2 storage client")
		}
	}

	return &Services{
		Job:          s.Job,
		Auth:         authService,
		DocumentRepo: repo,
		Storage:      store,
	}, nil
}
