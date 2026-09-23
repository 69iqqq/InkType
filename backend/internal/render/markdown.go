package render

import (
	"bytes"
	"context"
	"fmt"

	"inktype-backend/internal/documentir"
)

type MarkdownRenderer struct{}

func NewMarkdownRenderer() *MarkdownRenderer {
	return &MarkdownRenderer{}
}

func (r *MarkdownRenderer) CompileDocument(ctx context.Context, doc *documentir.Document) ([]byte, error) {
	var md bytes.Buffer

	for _, page := range doc.Pages {
		for _, block := range page.Blocks {
			text := ""
			if t, ok := block.Content["text"].(string); ok {
				text = t
			} else if t, ok := block.Content["text"].(map[string]interface{}); ok {
				// Fallback
				_ = t
			}

			switch block.Type {
			case documentir.TypeHeading:
				md.WriteString(fmt.Sprintf("# %s\n\n", text))
			case documentir.TypeParagraph:
				md.WriteString(fmt.Sprintf("%s\n\n", text))
			case documentir.TypeEquation:
				md.WriteString(fmt.Sprintf("$$ %s $$\n\n", text))
			case documentir.TypeCode:
				md.WriteString(fmt.Sprintf("```\n%s\n```\n\n", text))
			case documentir.TypeList:
				md.WriteString(fmt.Sprintf("- %s\n\n", text))
			default:
				md.WriteString(fmt.Sprintf("%s\n\n", text))
			}
		}
		md.WriteString("---\n\n")
	}

	return md.Bytes(), nil
}
