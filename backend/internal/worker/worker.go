package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	
	"inktype-backend/internal/ai"
	"inktype-backend/internal/config"
	"inktype-backend/internal/ocr"
	"inktype-backend/internal/pdf"
	"inktype-backend/internal/repository"
	"inktype-backend/internal/storage"
)

type Worker struct {
	cfg        *config.Config
	db         *pgxpool.Pool
	repo       *repository.Queries
	storage    *storage.Client
	pdfRender  pdf.Renderer
	vision     ai.VisionProvider
	ocr        ocr.OCRProvider
	logger     zerolog.Logger
	cancelFunc context.CancelFunc
	wg         sync.WaitGroup
}

func NewWorker(cfg *config.Config, db *pgxpool.Pool, repo *repository.Queries, storageClient *storage.Client, pdfRender pdf.Renderer, vision ai.VisionProvider, ocrProvider ocr.OCRProvider, logger zerolog.Logger) *Worker {
	return &Worker{
		cfg:       cfg,
		db:        db,
		repo:      repo,
		storage:   storageClient,
		pdfRender: pdfRender,
		vision:    vision,
		ocr:       ocrProvider,
		logger:    logger,
	}
}

func (w *Worker) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	w.cancelFunc = cancel

	concurrency := 5 // Default concurrency, should come from config
	w.logger.Info().Int("concurrency", concurrency).Msg("starting worker pool")

	for i := 0; i < concurrency; i++ {
		w.wg.Add(1)
		go w.workerLoop(ctx, i)
	}

	return nil
}

func (w *Worker) Stop() {
	if w.cancelFunc != nil {
		w.logger.Info().Msg("stopping workers...")
		w.cancelFunc()
		w.wg.Wait()
		w.logger.Info().Msg("all workers stopped")
	}
}

func (w *Worker) workerLoop(ctx context.Context, id int) {
	defer w.wg.Done()
	log := w.logger.With().Int("worker_id", id).Logger()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.processNextJob(ctx, log)
		}
	}
}

func (w *Worker) processNextJob(ctx context.Context, log zerolog.Logger) {
	// Try to claim a job with 30s lease
	job, err := w.repo.ClaimJob(ctx, 30)
	if err != nil {
		// pgx returns err no rows if no job is available
		if err.Error() == "no rows in result set" {
			return
		}
		log.Error().Err(err).Msg("failed to claim job")
		return
	}

	log.Info().Str("job_id", job.ID.String()).Str("type", job.Type).Msg("claimed job")
	start := time.Now()

	err = w.handleJob(ctx, job, log)
	duration := time.Since(start)

	if err != nil {
		log.Error().Err(err).Str("job_id", job.ID.String()).Dur("duration", duration).Msg("job failed")
		
		// If attempts < max_attempts, it will naturally be retried after lease expires
		// If it's a permanent error, we should fail it immediately.
		// For now, let's mark it as failed if attempts >= 3
		if job.Attempts >= 3 {
			failErr := w.repo.FailJob(ctx, repository.FailJobParams{
				ID:    job.ID,
				Error: pgtype.Text{String: err.Error(), Valid: true},
			})
			if failErr != nil {
				log.Error().Err(failErr).Msg("failed to mark job as failed")
			}
		}
	} else {
		log.Info().Str("job_id", job.ID.String()).Dur("duration", duration).Msg("job completed successfully")
		completeErr := w.repo.CompleteJob(ctx, job.ID)
		if completeErr != nil {
			log.Error().Err(completeErr).Msg("failed to mark job as completed")
		}
	}
}

func (w *Worker) handleJob(ctx context.Context, job repository.Job, log zerolog.Logger) error {
	switch job.Type {
	case "DOCUMENT_PROCESS":
		return w.processDocument(ctx, job, log)
	case "PAGE_ANALYZE":
		return w.processPage(ctx, job, log)
	default:
		return fmt.Errorf("unknown job type: %s", job.Type)
	}
}

func (w *Worker) processPage(ctx context.Context, job repository.Job, log zerolog.Logger) error {
	pages, err := w.repo.GetPendingDocumentPages(ctx, job.DocumentID)
	if err != nil {
		return fmt.Errorf("failed to get pending pages: %w", err)
	}

	for _, p := range pages {
		log.Info().Int32("page_number", p.PageNumber).Msg("analyzing page")

		doc, err := w.repo.GetDocumentByID(ctx, p.DocumentID)
		if err != nil {
			return fmt.Errorf("failed to get document %w", err)
		}

		embeddedText, err := w.pdfRender.ExtractText(ctx, doc.InputObjectKey, int(p.PageNumber))
		if err != nil {
			log.Warn().Err(err).Int32("page", p.PageNumber).Msg("failed to extract embedded text, continuing")
		}

		imageBytes := []byte("dummy image")
		var ocrResult *ocr.OCRResult

		if embeddedText == "" {
			log.Info().Int32("page", p.PageNumber).Msg("no embedded text found, running OCR")
			ocrRes, err := w.ocr.ExtractText(ctx, imageBytes)
			if err != nil {
				log.Error().Err(err).Int32("page", p.PageNumber).Msg("OCR failed, continuing without OCR")
			} else {
				ocrResult = &ocrRes
			}
		}

		res, err := w.vision.AnalyzePage(ctx, ai.PageAnalysisRequest{
			ImageObjectKey: p.ImageObjectKey,
			ImageBytes:     imageBytes,
			OCRResult:      ocrResult,
			EmbeddedText:   embeddedText,
		})

		if err != nil {
			errStr := err.Error()
			_ = w.repo.UpdateDocumentPageStatus(ctx, repository.UpdateDocumentPageStatusParams{
				ID:     p.ID,
				Status: "failed",
				Error:  pgtype.Text{String: errStr, Valid: true},
			})
			return fmt.Errorf("page analysis failed for page %d: %w", p.PageNumber, err)
		}

		// Save extracted JSON
		pageJson, _ := json.Marshal(res.Page)
		
		var ocrJson []byte
		if ocrResult != nil {
			ocrJson, _ = json.Marshal(ocrResult)
		}

		_ = w.repo.UpdateDocumentPageStatus(ctx, repository.UpdateDocumentPageStatusParams{
			ID:     p.ID,
			Status: "completed",
			ExtractedJson: pageJson,
			OcrResult: ocrJson,
			EmbeddedText: pgtype.Text{String: embeddedText, Valid: embeddedText != ""},
		})
	}

	// Enqueue DOCUMENT_RECONCILE when all pages are done
	// We could check if there are any pending pages left
	_, err = w.repo.EnqueueJob(ctx, repository.EnqueueJobParams{
		DocumentID: job.DocumentID,
		Type:       "DOCUMENT_RECONCILE",
	})
	if err != nil {
		return fmt.Errorf("failed to enqueue reconcile job: %w", err)
	}

	return nil
}

func (w *Worker) processDocument(ctx context.Context, job repository.Job, log zerolog.Logger) error {
	doc, err := w.repo.GetDocumentByID(ctx, job.DocumentID)
	if err != nil {
		return fmt.Errorf("failed to get document: %w", err)
	}
	
	// Update status to processing
	err = w.repo.UpdateDocumentStatus(ctx, repository.UpdateDocumentStatusParams{
		ID:     doc.ID,
		Status: "processing",
	})
	if err != nil {
		return fmt.Errorf("failed to update document status: %w", err)
	}

	// 2. Validate PDF (stub)
	err = w.pdfRender.ValidatePDF(ctx, doc.InputObjectKey)
	if err != nil {
		return fmt.Errorf("pdf validation failed: %w", err)
	}

	// 3. Render pages & 4. Upload to storage (stub)
	pageCount, err := w.pdfRender.GetPageCount(ctx, doc.InputObjectKey)
	if err != nil {
		return fmt.Errorf("failed to get page count: %w", err)
	}

	// Update document with page count
	err = w.repo.UpdateDocumentStatus(ctx, repository.UpdateDocumentStatusParams{
		ID:       doc.ID,
		Status:   "processing",
		Column3:  int32(pageCount),
	})
	if err != nil {
		return fmt.Errorf("failed to update document page count: %w", err)
	}

	// 5. Create document_page records and Enqueue PAGE_ANALYZE jobs
	for i := 1; i <= pageCount; i++ {
		imgData, err := w.pdfRender.RenderPage(ctx, doc.InputObjectKey, i)
		if err != nil {
			return fmt.Errorf("failed to render page %d: %w", i, err)
		}
		
		imageObjectKey := fmt.Sprintf("users/%s/documents/%s/pages/%d.png", doc.ClerkUserID, doc.ID.String(), i)
		if err := w.storage.Upload(ctx, imageObjectKey, imgData); err != nil {
			return fmt.Errorf("failed to upload page %d: %w", i, err)
		}

		_, err = w.repo.CreateDocumentPage(ctx, repository.CreateDocumentPageParams{
			DocumentID:     doc.ID,
			PageNumber:     int32(i),
			ImageObjectKey: imageObjectKey,
		})
		if err != nil {
			return fmt.Errorf("failed to create document_page record: %w", err)
		}

		_, err = w.repo.EnqueueJob(ctx, repository.EnqueueJobParams{
			DocumentID: doc.ID,
			Type:       "PAGE_ANALYZE",
		})
		if err != nil {
			return fmt.Errorf("failed to enqueue PAGE_ANALYZE job: %w", err)
		}
	}

	return nil
}
