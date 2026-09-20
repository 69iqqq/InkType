package pdf

import (
	"context"
	"fmt"
	"os"
)

type Renderer interface {
	GetPageCount(ctx context.Context, documentPath string) (int, error)
	RenderPage(ctx context.Context, documentPath string, pageNumber int) ([]byte, error)
	ValidatePDF(ctx context.Context, documentPath string) error
	ExtractText(ctx context.Context, documentPath string, pageNumber int) (string, error)
}

type SystemRenderer struct {
	// e.g. path to poppler pdftoppm or similar
}

func NewSystemRenderer() *SystemRenderer {
	return &SystemRenderer{}
}

func (r *SystemRenderer) ValidatePDF(ctx context.Context, documentPath string) error {
	// TODO: read magic bytes, check mime type, etc.
	info, err := os.Stat(documentPath)
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}
	if info.Size() == 0 {
		return fmt.Errorf("file is empty")
	}
	return nil
}

func (r *SystemRenderer) GetPageCount(ctx context.Context, documentPath string) (int, error) {
	// TODO: use real poppler pdfinfo
	return 1, nil // stub
}

func (r *SystemRenderer) RenderPage(ctx context.Context, documentPath string, pageNumber int) ([]byte, error) {
	// TODO: use real poppler pdftoppm
	return []byte("dummy image data"), nil // stub
}

func (r *SystemRenderer) ExtractText(ctx context.Context, documentPath string, pageNumber int) (string, error) {
	// TODO: use real poppler pdftotext
	return "", nil // returning empty string means no embedded text found
}
