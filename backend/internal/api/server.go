package api

import (
	"log/slog"
	"net/http"
	"time"

	"inktype-backend/internal/auth"
	"inktype-backend/internal/repository"
	"inktype-backend/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Server struct {
	router chi.Router
	logger *slog.Logger
	repo   repository.Querier
	store  *storage.Client
}

func NewServer(logger *slog.Logger, repo repository.Querier, store *storage.Client) *Server {
	s := &Server{
		router: chi.NewRouter(),
		logger: logger,
		repo:   repo,
		store:  store,
	}

	s.mountMiddleware()
	s.mountRoutes()

	return s
}

func (s *Server) mountMiddleware() {
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
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

		r.Post("/documents", s.handleCreateDocument)
		r.Get("/documents", s.handleListDocuments)
		r.Get("/documents/{id}", s.handleGetDocument)
		r.Delete("/documents/{id}", s.handleDeleteDocument)
	})
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
