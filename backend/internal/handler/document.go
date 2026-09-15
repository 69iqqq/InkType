package handler

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	
	"inktype-backend/internal/errs"
	"inktype-backend/internal/middleware"
	"inktype-backend/internal/repository"
	"inktype-backend/internal/server"
	"inktype-backend/internal/service"
)

type DocumentHandler struct {
	Handler
	services *service.Services
}

func NewDocumentHandler(s *server.Server, services *service.Services) *DocumentHandler {
	return &DocumentHandler{
		Handler:  NewHandler(s),
		services: services,
	}
}

type CreateDocumentRequest struct {
	Filename string `json:"filename"`
}

func (h *DocumentHandler) Create(c echo.Context) error {
	userID, ok := c.Get(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		return errs.NewUnauthorizedError("Unauthorized", false)
	}

	var req CreateDocumentRequest
	if err := c.Bind(&req); err != nil {
		return errs.NewBadRequestError("invalid request body", false, nil, nil, nil)
	}

	if req.Filename == "" {
		return errs.NewBadRequestError("filename is required", false, nil, nil, nil)
	}

	docID := uuid.New()
	objectKey := "uploads/" + userID + "/" + docID.String() + "/" + req.Filename

	presignedURL, err := h.services.Storage.GeneratePresignedUploadURL(c.Request().Context(), objectKey, 15*time.Minute)
	if err != nil {
		h.server.Logger.Error().Err(err).Msg("failed to generate presigned url")
		return errs.NewInternalServerError()
	}

	doc, err := h.services.DocumentRepo.CreateDocument(c.Request().Context(), repository.CreateDocumentParams{
		ClerkUserID:      userID,
		OriginalFilename: req.Filename,
		InputObjectKey:   objectKey,
	})
	if err != nil {
		h.server.Logger.Error().Err(err).Msg("failed to create document in db")
		return errs.NewInternalServerError()
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"document":   doc,
		"upload_url": presignedURL,
	})
}

func (h *DocumentHandler) List(c echo.Context) error {
	userID, ok := c.Get(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		return errs.NewUnauthorizedError("Unauthorized", false)
	}

	docs, err := h.services.DocumentRepo.ListDocuments(c.Request().Context(), userID)
	if err != nil {
		h.server.Logger.Error().Err(err).Msg("failed to list documents")
		return errs.NewInternalServerError()
	}

	return c.JSON(http.StatusOK, docs)
}

func (h *DocumentHandler) Get(c echo.Context) error {
	userID, ok := c.Get(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		return errs.NewUnauthorizedError("Unauthorized", false)
	}

	idParam := c.Param("id")
	docID, err := uuid.Parse(idParam)
	if err != nil {
		return errs.NewBadRequestError("invalid document id", false, nil, nil, nil)
	}

	pgDocID := pgtype.UUID{Bytes: docID, Valid: true}

	doc, err := h.services.DocumentRepo.GetDocument(c.Request().Context(), repository.GetDocumentParams{
		ID:          pgDocID,
		ClerkUserID: userID,
	})
	if err != nil {
		h.server.Logger.Error().Err(err).Msg("failed to get document")
		return errs.NewNotFoundError("document not found", false, nil)
	}

	return c.JSON(http.StatusOK, doc)
}

func (h *DocumentHandler) Delete(c echo.Context) error {
	userID, ok := c.Get(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		return errs.NewUnauthorizedError("Unauthorized", false)
	}

	idParam := c.Param("id")
	docID, err := uuid.Parse(idParam)
	if err != nil {
		return errs.NewBadRequestError("invalid document id", false, nil, nil, nil)
	}

	pgDocID := pgtype.UUID{Bytes: docID, Valid: true}

	err = h.services.DocumentRepo.DeleteDocument(c.Request().Context(), repository.DeleteDocumentParams{
		ID:          pgDocID,
		ClerkUserID: userID,
	})
	if err != nil {
		h.server.Logger.Error().Err(err).Msg("failed to delete document")
		return errs.NewInternalServerError()
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *DocumentHandler) Me(c echo.Context) error {
	userID, ok := c.Get(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		return errs.NewUnauthorizedError("Unauthorized", false)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"user_id": userID,
	})
}
