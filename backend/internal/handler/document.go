package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"

	"inktype-backend/internal/documentir"
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

	docs, err := h.services.DocumentRepo.ListDocuments(c.Request().Context(), repository.ListDocumentsParams{
		ClerkUserID: userID,
		Limit:       50,
		Offset:      0,
	})
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

	pages, err := h.services.DocumentRepo.GetDocumentPages(c.Request().Context(), pgDocID)
	if err != nil {
		h.server.Logger.Error().Err(err).Msg("failed to get document pages")
		return errs.NewInternalServerError()
	}

	type PageResponse struct {
		PageNumber int              `json:"page_number"`
		Status     string           `json:"status"`
		Blocks     []documentir.Block `json:"blocks"`
	}
	type DocumentResponse struct {
		repository.Document
		Pages  []PageResponse `json:"pages"`
		PDFUrl string         `json:"pdf_url,omitempty"`
		MDUrl  string         `json:"markdown_url,omitempty"`
	}

	var pageResponses []PageResponse
	for _, p := range pages {
		var pageIR documentir.Page
		if len(p.ExtractedJson) > 0 {
			_ = json.Unmarshal(p.ExtractedJson, &pageIR)
		}
		pageResponses = append(pageResponses, PageResponse{
			PageNumber: int(p.PageNumber),
			Status:     p.Status,
			Blocks:     pageIR.Blocks,
		})
	}

	resp := DocumentResponse{
		Document: doc,
		Pages:    pageResponses,
	}

	if doc.OutputObjectKey.Valid && doc.OutputObjectKey.String != "" {
		resp.PDFUrl = "/api/v1/documents/" + doc.ID.String() + "/download?format=pdf"
		resp.MDUrl = "/api/v1/documents/" + doc.ID.String() + "/download?format=md"
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *DocumentHandler) Download(c echo.Context) error {
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
        ID: pgDocID, ClerkUserID: userID,
    })
    if err != nil {
        return errs.NewNotFoundError("document not found", false, nil)
    }
    format := c.QueryParam("format")
    var objectKey string
    if format == "md" {
        objectKey = "output/" + docID.String() + "/result.md"
    } else {
        if doc.OutputObjectKey.Valid {
            objectKey = doc.OutputObjectKey.String
        } else {
            return errs.NewNotFoundError("output not available yet", false, nil)
        }
    }
    reader, err := h.services.Storage.GetFile(c.Request().Context(), objectKey)
    if err != nil {
        return errs.NewNotFoundError("file not found", false, nil)
    }
    defer reader.Close()
    contentType := "application/pdf"
    if format == "md" {
        contentType = "text/markdown"
    }
    c.Response().Header().Set("Content-Type", contentType)
    c.Response().Header().Set("Content-Disposition", "attachment")
    _, err = io.Copy(c.Response().Writer, reader)
    return err
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

func (h *DocumentHandler) Convert(c echo.Context) error {
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

	// Verify document exists and belongs to user
	_, err = h.services.DocumentRepo.GetDocument(c.Request().Context(), repository.GetDocumentParams{
		ID:          pgDocID,
		ClerkUserID: userID,
	})
	if err != nil {
		return errs.NewNotFoundError("document not found", false, nil)
	}

	// Enqueue job
	job, err := h.services.DocumentRepo.EnqueueJob(c.Request().Context(), repository.EnqueueJobParams{
		DocumentID:     pgDocID,
		DocumentPageID: pgtype.UUID{Valid: false},
		Type:           "DOCUMENT_PROCESS",
	})
	if err != nil {
		h.server.Logger.Error().Err(err).Msg("failed to enqueue job")
		return errs.NewInternalServerError()
	}

	return c.JSON(http.StatusAccepted, job)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *DocumentHandler) WS(c echo.Context) error {
	userID, ok := c.Get(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		return errs.NewUnauthorizedError("Unauthorized", false)
	}

	idParam := c.Param("id")
	docID, err := uuid.Parse(idParam)
	if err != nil {
		return errs.NewBadRequestError("invalid document id", false, nil, nil, nil)
	}

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	pgDocID := pgtype.UUID{Bytes: docID, Valid: true}
	ctx := c.Request().Context()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			doc, err := h.services.DocumentRepo.GetDocument(ctx, repository.GetDocumentParams{
				ID:          pgDocID,
				ClerkUserID: userID,
			})
			if err != nil {
				_ = ws.WriteJSON(map[string]interface{}{"error": "document not found"})
				return nil
			}

			pages, err := h.services.DocumentRepo.GetDocumentPages(ctx, pgDocID)
			completed := 0
			if err == nil {
				for _, p := range pages {
					if p.Status == "completed" {
						completed++
					}
				}
			}

			pagesTotal := 0
			if doc.PageCount.Valid {
				pagesTotal = int(doc.PageCount.Int32)
			} else {
				pagesTotal = len(pages)
			}

			msg := map[string]interface{}{
				"status":          doc.Status,
				"pages_total":     pagesTotal,
				"pages_completed": completed,
			}

			if err := ws.WriteJSON(msg); err != nil {
				return nil
			}

			if doc.Status == "completed" || doc.Status == "failed" {
				return nil
			}
		}
	}
}

func (h *DocumentHandler) LocalUpload(c echo.Context) error {
	objectKey := c.QueryParam("key")
	if objectKey == "" {
		return errs.NewBadRequestError("missing key", false, nil, nil, nil)
	}

	objectKey = filepath.Clean(objectKey)
	if strings.HasPrefix(objectKey, "..") || strings.HasPrefix(objectKey, "/") || strings.HasPrefix(objectKey, "\\") {
		return errs.NewBadRequestError("invalid key", false, nil, nil, nil)
	}

	c.Request().Body = http.MaxBytesReader(c.Response(), c.Request().Body, 100*1024*1024)

	bodyBytes, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return errs.NewInternalServerError()
	}

	if err := h.services.Storage.Upload(c.Request().Context(), objectKey, bodyBytes); err != nil {
		h.server.Logger.Error().Err(err).Msg("failed to upload local file")
		return errs.NewInternalServerError()
	}

	return c.NoContent(http.StatusOK)
}
