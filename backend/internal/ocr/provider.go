package ocr

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
)

type StubOCRProvider struct{}

func NewStubOCRProvider() *StubOCRProvider {
	return &StubOCRProvider{}
}

func (p *StubOCRProvider) ExtractText(ctx context.Context, image []byte) (OCRResult, error) {
	tempFile, err := os.CreateTemp("", "inktype-ocr-*.png")
	if err != nil {
		return OCRResult{}, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.Write(image); err != nil {
		tempFile.Close()
		return OCRResult{}, fmt.Errorf("failed to write image to temp file: %w", err)
	}
	tempFile.Close()

	cmd := exec.CommandContext(ctx, "tesseract", tempFile.Name(), "stdout", "-l", "eng", "--psm", "3")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return OCRResult{}, fmt.Errorf("tesseract failed: %w", err)
	}

	return OCRResult{
		ExtractedText: stdout.String(),
	}, nil
}
