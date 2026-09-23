package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"inktype-backend/internal/ai"
	"inktype-backend/internal/config"
	"inktype-backend/internal/documentir"
	"inktype-backend/internal/ocr"
	"inktype-backend/internal/pdf"
	"inktype-backend/internal/render"
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
	case "DOCUMENT_RECONCILE":
		return w.reconcileDocument(ctx, job, log)
	default:
		return fmt.Errorf("unknown job type: %s", job.Type)
	}
}

func (w *Worker) processPage(ctx context.Context, job repository.Job, log zerolog.Logger) error {
	p, err := w.repo.GetDocumentPage(ctx, job.DocumentPageID)
	if err != nil {
		return fmt.Errorf("failed to get page: %w", err)
	}

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
	imgReader, err := w.storage.GetFile(ctx, p.ImageObjectKey)
	if err != nil {
		return fmt.Errorf("failed to get page image from storage: %w", err)
	}
	defer imgReader.Close()

	imageBytes, err = io.ReadAll(imgReader)
	if err != nil {
		return fmt.Errorf("failed to read page image: %w", err)
	}

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
		PageNumber:     int(p.PageNumber),
	})

	if err != nil {
		errStr := err.Error()
		_, _ = w.db.Exec(ctx, `
			UPDATE document_pages
			SET status = $2, error = $3, extracted_json = '{}'::jsonb, ocr_result = '{}'::jsonb, updated_at = CURRENT_TIMESTAMP
			WHERE id = $1
		`, p.ID, "failed", errStr)
		return fmt.Errorf("page analysis failed for page %d: %w", p.PageNumber, err)
	}

	// Save extracted JSON
	pageJson, _ := json.Marshal(res.Page)

	ocrJsonStr := "{}"
	if ocrResult != nil {
		ocrJsonBytes, _ := json.Marshal(ocrResult)
		ocrJsonStr = string(ocrJsonBytes)
	}

	// IMPORTANT: We bypass the sqlc-generated function here because pgx Simple Protocol
	// (required for Neon PgBouncer) sends []byte as hex-encoded bytea, which Postgres
	// rejects when parsing as JSONB. By passing string instead, pgx sends it as text.
	_, err = w.db.Exec(ctx, `
		UPDATE document_pages
		SET status = $2,
		    extracted_json = $3::jsonb,
		    ocr_result = $4::jsonb,
		    embedded_text = COALESCE($5, embedded_text),
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`, p.ID, "completed", string(pageJson), ocrJsonStr,
		pgtype.Text{String: embeddedText, Valid: embeddedText != ""})
	if err != nil {
		log.Error().Err(err).Int32("page", p.PageNumber).Msg("failed to save page results")
		return fmt.Errorf("failed to save page results: %w", err)
	}

	count, err := w.repo.CountUnfinishedPages(ctx, job.DocumentID)
	if err != nil {
		return fmt.Errorf("failed to count unfinished pages: %w", err)
	}
	if count == 0 {
		_, err = w.repo.EnqueueJob(ctx, repository.EnqueueJobParams{
			DocumentID:     job.DocumentID,
			DocumentPageID: pgtype.UUID{}, // empty
			Type:           "DOCUMENT_RECONCILE",
		})
		if err != nil {
			return fmt.Errorf("failed to enqueue reconcile job: %w", err)
		}
	}

	return nil
}

func (w *Worker) reconcileDocument(ctx context.Context, job repository.Job, log zerolog.Logger) error {
	pages, err := w.repo.GetDocumentPages(ctx, job.DocumentID)
	if err != nil {
		return fmt.Errorf("failed to get document pages: %w", err)
	}

	docIR := &documentir.Document{}
	for _, p := range pages {
		var pageIR documentir.Page
		if len(p.ExtractedJson) > 0 {
			if err := json.Unmarshal(p.ExtractedJson, &pageIR); err != nil {
				log.Warn().Err(err).Int32("page", p.PageNumber).Msg("failed to parse page json")
				continue
			}
			docIR.Pages = append(docIR.Pages, pageIR)
		}
	}

	renderer := render.NewTypstRenderer()
	pdfBytes, err := renderer.CompileDocument(ctx, docIR, "minimal")
	if err != nil {
		return fmt.Errorf("failed to compile document: %w", err)
	}

	outputKey := "output/" + job.DocumentID.String() + "/result.pdf"
	if err := w.storage.Upload(ctx, outputKey, pdfBytes); err != nil {
		return fmt.Errorf("failed to upload compiled document: %w", err)
	}

	mdRenderer := render.NewMarkdownRenderer()
	mdBytes, err := mdRenderer.CompileDocument(ctx, docIR)
	if err != nil {
		return fmt.Errorf("failed to compile markdown: %w", err)
	}

	mdOutputKey := "output/" + job.DocumentID.String() + "/result.md"
	if err := w.storage.Upload(ctx, mdOutputKey, mdBytes); err != nil {
		return fmt.Errorf("failed to upload markdown document: %w", err)
	}

	log.Info().Int("page_count", len(pages)).Msg("reconciliation complete")

	return w.repo.CompleteDocument(ctx, repository.CompleteDocumentParams{
		ID:              job.DocumentID,
		OutputObjectKey: pgtype.Text{String: outputKey, Valid: true},
	})
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
		_ = w.repo.FailDocument(ctx, repository.FailDocumentParams{ID: doc.ID, Error: pgtype.Text{String: err.Error(), Valid: true}})
		return fmt.Errorf("failed to update document status: %w", err)
	}

	// 2. Validate PDF (stub)
	err = w.pdfRender.ValidatePDF(ctx, doc.InputObjectKey)
	if err != nil {
		_ = w.repo.FailDocument(ctx, repository.FailDocumentParams{ID: doc.ID, Error: pgtype.Text{String: err.Error(), Valid: true}})
		return fmt.Errorf("pdf validation failed: %w", err)
	}

	// 3. Render pages & 4. Upload to storage (stub)
	pageCount, err := w.pdfRender.GetPageCount(ctx, doc.InputObjectKey)
	if err != nil {
		_ = w.repo.FailDocument(ctx, repository.FailDocumentParams{ID: doc.ID, Error: pgtype.Text{String: err.Error(), Valid: true}})
		return fmt.Errorf("failed to get page count: %w", err)
	}

	// Update document with page count
	err = w.repo.UpdateDocumentStatus(ctx, repository.UpdateDocumentStatusParams{
		ID:      doc.ID,
		Status:  "processing",
		Column3: int32(pageCount),
	})
	if err != nil {
		_ = w.repo.FailDocument(ctx, repository.FailDocumentParams{ID: doc.ID, Error: pgtype.Text{String: err.Error(), Valid: true}})
		return fmt.Errorf("failed to update document page count: %w", err)
	}

	// Get existing pages to skip duplicates on retry
	existingPages, err := w.repo.GetDocumentPages(ctx, doc.ID)
	if err != nil {
		_ = w.repo.FailDocument(ctx, repository.FailDocumentParams{ID: doc.ID, Error: pgtype.Text{String: err.Error(), Valid: true}})
		return fmt.Errorf("failed to check existing pages: %w", err)
	}
	existingPageMap := make(map[int32]bool)
	for _, p := range existingPages {
		existingPageMap[p.PageNumber] = true
	}

	// 5. Create document_page records and Enqueue PAGE_ANALYZE jobs
	for i := 1; i <= pageCount; i++ {
		if existingPageMap[int32(i)] {
			log.Info().Int("page", i).Msg("page already exists, skipping creation")
			continue
		}

		imgData, err := w.pdfRender.RenderPage(ctx, doc.InputObjectKey, i)
		if err != nil {
			_ = w.repo.FailDocument(ctx, repository.FailDocumentParams{ID: doc.ID, Error: pgtype.Text{String: err.Error(), Valid: true}})
			return fmt.Errorf("failed to render page %d: %w", i, err)
		}

		imageObjectKey := fmt.Sprintf("users/%s/documents/%s/pages/%d.png", doc.ClerkUserID, doc.ID.String(), i)
		if err := w.storage.Upload(ctx, imageObjectKey, imgData); err != nil {
			_ = w.repo.FailDocument(ctx, repository.FailDocumentParams{ID: doc.ID, Error: pgtype.Text{String: err.Error(), Valid: true}})
			return fmt.Errorf("failed to upload page %d: %w", i, err)
		}

		page, err := w.repo.CreateDocumentPage(ctx, repository.CreateDocumentPageParams{
			DocumentID:     doc.ID,
			PageNumber:     int32(i),
			ImageObjectKey: imageObjectKey,
		})
		if err != nil {
			_ = w.repo.FailDocument(ctx, repository.FailDocumentParams{ID: doc.ID, Error: pgtype.Text{String: err.Error(), Valid: true}})
			return fmt.Errorf("failed to create document_page record: %w", err)
		}

		_, err = w.repo.EnqueueJob(ctx, repository.EnqueueJobParams{
			DocumentID:     doc.ID,
			DocumentPageID: page.ID,
			Type:           "PAGE_ANALYZE",
		})
		if err != nil {
			_ = w.repo.FailDocument(ctx, repository.FailDocumentParams{ID: doc.ID, Error: pgtype.Text{String: err.Error(), Valid: true}})
			return fmt.Errorf("failed to enqueue PAGE_ANALYZE job: %w", err)
		}
	}

	return nil
}
