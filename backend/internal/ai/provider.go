package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"inktype-backend/internal/documentir"

	"google.golang.org/genai"
)

type GenAIVisionProvider struct {
	client *genai.Client
}

func NewGenAIVisionProvider(ctx context.Context) (*GenAIVisionProvider, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{})
	if err != nil {
		return nil, err
	}
	return &GenAIVisionProvider{client: client}, nil
}

func (p *GenAIVisionProvider) AnalyzePage(ctx context.Context, req PageAnalysisRequest) (PageAnalysisResult, error) {
	prompt := `Analyze this image and the provided OCR text. Return a JSON array of blocks. Each block must have:
- "type" (one of: heading, paragraph, equation, code, list)
- "content" (a JSON object with a "text" string field containing the extracted text/math)
- "confidence" (number between 0.0 and 1.0)
Do not include markdown blocks like ` + "```json" + `. Return ONLY the JSON array.`

	ocrText := ""
	if req.OCRResult != nil {
		ocrText = req.OCRResult.ExtractedText
	}

	// Create request content
	var parts []*genai.Part
	parts = append(parts, genai.NewPartFromText(prompt))
	if ocrText != "" {
		parts = append(parts, genai.NewPartFromText("OCR Text: "+ocrText))
	}
	if len(req.ImageBytes) > 0 {
		parts = append(parts, genai.NewPartFromBytes(req.ImageBytes, "image/png"))
	}

	models := []string{"gemini-2.5-flash-lite", "gemini-3.5-flash", "gemini-3.6-flash", "gemini-3.7-flash"}
	var resp *genai.GenerateContentResponse
	var lastErr error
	for _, modelName := range models {
		resp, lastErr = p.client.Models.GenerateContent(ctx, modelName, []*genai.Content{
			{Role: "user", Parts: parts},
		}, nil)
		if lastErr == nil && resp.Candidates != nil && len(resp.Candidates) > 0 {
			break
		}
	}
	if lastErr != nil {
		return PageAnalysisResult{}, fmt.Errorf("gemini API error: %w", lastErr)
	}

	if resp.Candidates == nil || len(resp.Candidates) == 0 {
		return PageAnalysisResult{}, fmt.Errorf("no candidates returned")
	}

	respText := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		if part.Text != "" {
			respText += part.Text
		}
	}

	// Clean up markdown just in case
	cleanJSON := strings.TrimSpace(respText)
	if strings.HasPrefix(cleanJSON, "```json") {
		cleanJSON = strings.TrimPrefix(cleanJSON, "```json")
		cleanJSON = strings.TrimSuffix(cleanJSON, "```")
	}

	var blocks []documentir.Block
	if err := json.Unmarshal([]byte(cleanJSON), &blocks); err != nil {
		return PageAnalysisResult{}, fmt.Errorf("failed to parse JSON: %w (raw: %s)", err, respText)
	}

	return PageAnalysisResult{
		Page: documentir.Page{
			PageNumber: req.PageNumber, // Defaulting to 1, could be mapped if provided
			Blocks:     blocks,
		},
		Raw: respText,
	}, nil
}
