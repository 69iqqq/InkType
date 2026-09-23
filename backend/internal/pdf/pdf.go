package pdf

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"io"
	"os"

	"inktype-backend/internal/storage"

	"github.com/gen2brain/go-fitz"
)

type Renderer interface {
	GetPageCount(ctx context.Context, documentPath string) (int, error)
	RenderPage(ctx context.Context, documentPath string, pageNumber int) ([]byte, error)
	ValidatePDF(ctx context.Context, documentPath string) error
	ExtractText(ctx context.Context, documentPath string, pageNumber int) (string, error)
}

type SystemRenderer struct {
	storageClient *storage.Client
}

func NewSystemRenderer(storageClient *storage.Client) *SystemRenderer {
	return &SystemRenderer{
		storageClient: storageClient,
	}
}

func (r *SystemRenderer) downloadToTemp(ctx context.Context, documentPath string) (string, error) {
	reader, err := r.storageClient.GetFile(ctx, documentPath)
	if err != nil {
		return "", fmt.Errorf("failed to get file from storage: %w", err)
	}
	defer reader.Close()

	tempFile, err := os.CreateTemp("", "inktype-pdf-*.pdf")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tempFile.Close()

	if _, err := io.Copy(tempFile, reader); err != nil {
		os.Remove(tempFile.Name())
		return "", fmt.Errorf("failed to download file: %w", err)
	}

	return tempFile.Name(), nil
}

func (r *SystemRenderer) ValidatePDF(ctx context.Context, documentPath string) error {
	tempPath, err := r.downloadToTemp(ctx, documentPath)
	if err != nil {
		return err
	}
	defer os.Remove(tempPath)

	doc, err := fitz.New(tempPath)
	if err != nil {
		return fmt.Errorf("invalid PDF: %w", err)
	}
	doc.Close()
	return nil
}

func (r *SystemRenderer) GetPageCount(ctx context.Context, documentPath string) (int, error) {
	tempPath, err := r.downloadToTemp(ctx, documentPath)
	if err != nil {
		return 0, err
	}
	defer os.Remove(tempPath)

	doc, err := fitz.New(tempPath)
	if err != nil {
		return 0, fmt.Errorf("failed to open PDF: %w", err)
	}
	defer doc.Close()

	return doc.NumPage(), nil
}

func (r *SystemRenderer) RenderPage(ctx context.Context, documentPath string, pageNumber int) ([]byte, error) {
	tempPath, err := r.downloadToTemp(ctx, documentPath)
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempPath)

	doc, err := fitz.New(tempPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open PDF: %w", err)
	}
	defer doc.Close()

	img, err := doc.Image(pageNumber - 1)
	if err != nil {
		return nil, fmt.Errorf("failed to render page: %w", err)
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("failed to encode PNG: %w", err)
	}

	return buf.Bytes(), nil
}

func (r *SystemRenderer) ExtractText(ctx context.Context, documentPath string, pageNumber int) (string, error) {
	tempPath, err := r.downloadToTemp(ctx, documentPath)
	if err != nil {
		return "", err
	}
	defer os.Remove(tempPath)

	doc, err := fitz.New(tempPath)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF: %w", err)
	}
	defer doc.Close()

	text, err := doc.Text(pageNumber - 1)
	if err != nil {
		return "", fmt.Errorf("failed to extract text: %w", err)
	}
	return text, nil
}
