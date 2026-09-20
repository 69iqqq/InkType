package ai

import (
	"context"

	"inktype-backend/internal/documentir"
)

type StubVisionProvider struct {
}

func NewStubVisionProvider() *StubVisionProvider {
	return &StubVisionProvider{}
}

func (p *StubVisionProvider) AnalyzePage(ctx context.Context, req PageAnalysisRequest) (PageAnalysisResult, error) {
	// 5. Validate AI JSON (simulated)
	// Return a stub PageIR
	return PageAnalysisResult{
		Page: documentir.Page{
			PageNumber: 1,
			Blocks: []documentir.Block{
				{
					ID:         "1",
					Type:       documentir.TypeHeading,
					Confidence: 0.99,
					Content: map[string]interface{}{
						"text": "Extracted Heading",
					},
				},
			},
		},
	}, nil
}
