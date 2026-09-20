package ocr

import (
	"context"
)

type StubOCRProvider struct{}

func NewStubOCRProvider() *StubOCRProvider {
	return &StubOCRProvider{}
}

func (p *StubOCRProvider) ExtractText(ctx context.Context, image []byte) (OCRResult, error) {
	return OCRResult{
		ExtractedText: "stub OCR output",
	}, nil
}
