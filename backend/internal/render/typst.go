package render

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"inktype-backend/internal/documentir"
)

type TypstRenderer struct{}

func NewTypstRenderer() *TypstRenderer {
	return &TypstRenderer{}
}

func (r *TypstRenderer) CompileDocument(ctx context.Context, doc *documentir.Document, templateName string) ([]byte, error) {
	var typstCode bytes.Buffer

	if templateName == "academic" {
		typstCode.WriteString("#set page(margin: 1.5in)\n\n")
	} else {
		typstCode.WriteString("#set page(margin: 1in)\n\n")
	}

	for _, page := range doc.Pages {
		for _, block := range page.Blocks {
			text := ""
			if t, ok := block.Content["text"].(string); ok {
				text = t
			} else if t, ok := block.Content["text"].(map[string]interface{}); ok {
				// Fallback if structured weirdly
				_ = t
			}

			// Escape Typst syntax characters in plain text
			escapedText := strings.ReplaceAll(text, "@", "\\@")

			switch block.Type {
			case documentir.TypeHeading:
				typstCode.WriteString(fmt.Sprintf("= %s\n\n", escapedText))
			case documentir.TypeParagraph:
				typstCode.WriteString(fmt.Sprintf("%s\n\n", escapedText))
			case documentir.TypeEquation:
				typstCode.WriteString(fmt.Sprintf("$ %s $\n\n", text))
			case documentir.TypeCode:
				typstCode.WriteString(fmt.Sprintf("```\n%s\n```\n\n", text))
			case documentir.TypeList:
				typstCode.WriteString(fmt.Sprintf("- %s\n\n", escapedText))
			default:
				typstCode.WriteString(fmt.Sprintf("%s\n\n", escapedText))
			}
		}
		typstCode.WriteString("#pagebreak()\n")
	}

	tempDir, err := os.MkdirTemp("", "typst-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	typPath := filepath.Join(tempDir, "main.typ")
	pdfPath := filepath.Join(tempDir, "main.pdf")

	if err := os.WriteFile(typPath, typstCode.Bytes(), 0644); err != nil {
		return nil, fmt.Errorf("failed to write typst file: %w", err)
	}

	typstCmd := "typst"
	if _, err := os.Stat("typst.exe"); err == nil {
		typstCmd = "./typst.exe"
	} else if _, err := os.Stat("typst"); err == nil {
		typstCmd = "./typst"
	}

	cmd := exec.CommandContext(ctx, typstCmd, "compile", typPath, pdfPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		// FALLBACK: If Typst compilation failed due to LaTeX math syntax or unescaped characters,
		// compile a safe fallback where equations and complex blocks are rendered as raw code/text blocks.
		var fallbackCode bytes.Buffer
		fallbackCode.WriteString("#set page(margin: 1in)\n\n")
		for _, page := range doc.Pages {
			for _, block := range page.Blocks {
				text := ""
				if t, ok := block.Content["text"].(string); ok {
					text = t
				}
				safe := strings.ReplaceAll(text, "@", "\\@")
				safe = strings.ReplaceAll(safe, "#", "\\#")

				switch block.Type {
				case documentir.TypeHeading:
					fallbackCode.WriteString(fmt.Sprintf("= %s\n\n", safe))
				case documentir.TypeEquation:
					// Wrap in raw block so complex LaTeX matrices/fractions render without syntax errors
					fallbackCode.WriteString(fmt.Sprintf("```\n%s\n```\n\n", text))
				default:
					fallbackCode.WriteString(fmt.Sprintf("%s\n\n", safe))
				}
			}
			fallbackCode.WriteString("#pagebreak()\n")
		}

		if writeErr := os.WriteFile(typPath, fallbackCode.Bytes(), 0644); writeErr == nil {
			cmdFallback := exec.CommandContext(ctx, typstCmd, "compile", typPath, pdfPath)
			if _, fbErr := cmdFallback.CombinedOutput(); fbErr != nil {
				return nil, fmt.Errorf("typst compile failed: %w, output: %s", err, string(output))
			}
		} else {
			return nil, fmt.Errorf("typst compile failed: %w, output: %s", err, string(output))
		}
	}

	pdfBytes, err := os.ReadFile(pdfPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read pdf output: %w", err)
	}

	return pdfBytes, nil
}
