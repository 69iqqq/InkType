package ai

import (
	"context"

	"inktype-backend/internal/documentir"
	"inktype-backend/internal/ocr"
)

type PageAnalysisRequest struct {
	ImageObjectKey string
	ImageBytes     []byte
	OCRResult      *ocr.OCRResult
	EmbeddedText   string
	PageNumber     int
}

type PageAnalysisResult struct {
	Page documentir.Page
	Raw  string
}

type VisionProvider interface {
	AnalyzePage(ctx context.Context, req PageAnalysisRequest) (PageAnalysisResult, error)
}
