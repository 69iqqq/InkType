package api

import (
	"log/slog"
	"net/http"
	"time"

	"inktype-backend/internal/auth"
	"inktype-backend/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	router chi.Router
	logger *slog.Logger
	repo   repository.Querier
}

func NewServer(logger *slog.Logger, repo repository.Querier) *Server {
	s := &Server{
		router: chi.NewRouter(),
		logger: logger,
		repo:   repo,
	}

	s.mountMiddleware()
	s.mountRoutes()

	return s
}

func (s *Server) mountMiddleware() {
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.ClientIPFromRemoteAddr)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(60 * time.Second))
}

func (s *Server) mountRoutes() {
	s.router.Get("/health", s.handleHealthCheck)

	s.router.Route("/v1", func(r chi.Router) {
		r.Use(auth.RequireAuth(s.logger))
		
		r.Get("/me", s.handleMe)
	})
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
