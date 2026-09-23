-- name: CreateDocument :one
INSERT INTO documents (
    clerk_user_id, original_filename, input_object_key, status
) VALUES (
    $1, $2, $3, 'uploaded'
)
RETURNING id, clerk_user_id, original_filename, input_object_key, output_object_key, page_count, status, error, created_at, updated_at;

-- name: GetDocument :one
SELECT id, clerk_user_id, original_filename, input_object_key, output_object_key, page_count, status, error, created_at, updated_at 
FROM documents 
WHERE id = $1 AND clerk_user_id = $2;

-- name: ListDocuments :many
SELECT id, clerk_user_id, original_filename, input_object_key, output_object_key, page_count, status, error, created_at, updated_at 
FROM documents 
WHERE clerk_user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: DeleteDocument :exec
DELETE FROM documents
WHERE id = $1 AND clerk_user_id = $2;

-- name: GetDocumentByID :one
SELECT id, clerk_user_id, original_filename, input_object_key, output_object_key, page_count, status, error, created_at, updated_at 
FROM documents 
WHERE id = $1;

-- name: UpdateDocumentStatus :exec
UPDATE documents
SET status = $2,
    page_count = COALESCE(NULLIF($3::integer, 0), page_count),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: CompleteDocument :exec
UPDATE documents
SET status = 'completed',
    output_object_key = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: FailDocument :exec
UPDATE documents
SET status = 'failed',
    error = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: EnqueueJob :one
INSERT INTO jobs (
    document_id, document_page_id, type, status
) VALUES (
    $1, $2, $3, 'pending'
)
RETURNING id, document_id, document_page_id, type, status, created_at;

-- name: ClaimJob :one
UPDATE jobs
SET status = 'processing',
    attempts = attempts + 1,
    started_at = CURRENT_TIMESTAMP,
    available_at = CURRENT_TIMESTAMP + ($1::integer * interval '1 second')
WHERE id = (
    SELECT id
    FROM jobs
    WHERE status = 'pending' OR (status = 'processing' AND available_at < CURRENT_TIMESTAMP)
    ORDER BY created_at ASC
    FOR UPDATE SKIP LOCKED
    LIMIT 1
)
RETURNING id, document_id, document_page_id, type, status, attempts, available_at, started_at, completed_at, error, created_at, updated_at;

-- name: CompleteJob :exec
UPDATE jobs
SET status = 'completed',
    completed_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: FailJob :exec
UPDATE jobs
SET status = 'failed',
    error = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: CreateDocumentPage :one
INSERT INTO document_pages (
    document_id, page_number, image_object_key, status
) VALUES (
    $1, $2, $3, 'pending'
) RETURNING id, document_id, page_number, image_object_key, extracted_json, ocr_result, embedded_text, status, error, created_at, updated_at;

-- name: GetDocumentPage :one
SELECT id, document_id, page_number, image_object_key, extracted_json, ocr_result, embedded_text, status, error, created_at, updated_at
FROM document_pages
WHERE id = $1;

-- name: GetPendingDocumentPages :many
SELECT id, document_id, page_number, image_object_key, extracted_json, ocr_result, embedded_text, status, error, created_at, updated_at
FROM document_pages
WHERE document_id = $1 AND status = 'pending'
ORDER BY page_number ASC;

-- name: GetDocumentPages :many
SELECT id, document_id, page_number, image_object_key, extracted_json, ocr_result, embedded_text, status, error, created_at, updated_at
FROM document_pages
WHERE document_id = $1
ORDER BY page_number ASC;

-- name: CountUnfinishedPages :one
SELECT COUNT(*)::integer AS count
FROM document_pages
WHERE document_id = $1 AND status NOT IN ('completed', 'failed');

-- name: UpdateDocumentPageStatus :exec
UPDATE document_pages
SET status = $2,
    extracted_json = COALESCE($3, extracted_json),
    ocr_result = COALESCE($4, ocr_result),
    embedded_text = COALESCE($5, embedded_text),
    error = COALESCE($6, error),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;
