package ocr

import (
	"context"
)

type OCRResult struct {
	ExtractedText string `json:"extracted_text"`
	// Additional bounding boxes, confidence scores can be added here
}

type OCRProvider interface {
	ExtractText(ctx context.Context, image []byte) (OCRResult, error)
}
