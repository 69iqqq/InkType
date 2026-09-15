package api

import (
	"encoding/json"
	"net/http"
	"time"

	"inktype-backend/internal/auth"
	"inktype-backend/internal/httperr"
	"inktype-backend/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateDocumentRequest struct {
	Filename string `json:"filename"`
}

func (s *Server) handleCreateDocument(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		httperr.InternalError(w)
		return
	}

	var req CreateDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperr.BadRequest(w, "invalid request body")
		return
	}

	if req.Filename == "" {
		httperr.BadRequest(w, "filename is required")
		return
	}

	docID := uuid.New()
	objectKey := "uploads/" + userID + "/" + docID.String() + "/" + req.Filename

	presignedURL, err := s.store.GeneratePresignedUploadURL(r.Context(), objectKey, 15*time.Minute)
	if err != nil {
		s.logger.ErrorContext(r.Context(), "failed to generate presigned url", "error", err)
		httperr.InternalError(w)
		return
	}

	doc, err := s.repo.CreateDocument(r.Context(), repository.CreateDocumentParams{
		ClerkUserID:      userID,
		OriginalFilename: req.Filename,
		InputObjectKey:   objectKey,
	})
	if err != nil {
		s.logger.ErrorContext(r.Context(), "failed to create document in db", "error", err)
		httperr.InternalError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"document":   doc,
		"upload_url": presignedURL,
	})
}

func (s *Server) handleListDocuments(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		httperr.InternalError(w)
		return
	}

	docs, err := s.repo.ListDocuments(r.Context(), userID)
	if err != nil {
		s.logger.ErrorContext(r.Context(), "failed to list documents", "error", err)
		httperr.InternalError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(docs)
}

func (s *Server) handleGetDocument(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		httperr.InternalError(w)
		return
	}

	idParam := chi.URLParam(r, "id")
	docID, err := uuid.Parse(idParam)
	if err != nil {
		httperr.BadRequest(w, "invalid document id")
		return
	}
	
	pgDocID := pgtype.UUID{Bytes: docID, Valid: true}

	doc, err := s.repo.GetDocument(r.Context(), repository.GetDocumentParams{
		ID:          pgDocID,
		ClerkUserID: userID,
	})
	if err != nil {
		s.logger.ErrorContext(r.Context(), "failed to get document", "error", err)
		httperr.JSON(w, "document not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doc)
}

func (s *Server) handleDeleteDocument(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		httperr.InternalError(w)
		return
	}

	idParam := chi.URLParam(r, "id")
	docID, err := uuid.Parse(idParam)
	if err != nil {
		httperr.BadRequest(w, "invalid document id")
		return
	}

	pgDocID := pgtype.UUID{Bytes: docID, Valid: true}

	err = s.repo.DeleteDocument(r.Context(), repository.DeleteDocumentParams{
		ID:          pgDocID,
		ClerkUserID: userID,
	})
	if err != nil {
		s.logger.ErrorContext(r.Context(), "failed to delete document", "error", err)
		httperr.InternalError(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
