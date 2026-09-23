package handler

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"inktype-backend/internal/server"

	"github.com/labstack/echo/v4"
)

type OpenAPIHandler struct {
	Handler
	templateOnce sync.Once
	templateHTML string
}

func NewOpenAPIHandler(s *server.Server) *OpenAPIHandler {
	return &OpenAPIHandler{
		Handler: NewHandler(s),
	}
}

func (h *OpenAPIHandler) ServeOpenAPIUI(c echo.Context) error {
	h.templateOnce.Do(func() {
		data, err := os.ReadFile("static/openapi.html")
		if err != nil {
			h.server.Logger.Error().Err(err).Msg("failed to load openapi template")
			return
		}
		h.templateHTML = string(data)
	})

	if h.templateHTML == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, "OpenAPI template not available")
	}

	c.Response().Header().Set("Cache-Control", "no-cache")

	err := c.HTML(http.StatusOK, h.templateHTML)
	if err != nil {
		return fmt.Errorf("failed to write HTML response: %w", err)
	}

	return nil
}
