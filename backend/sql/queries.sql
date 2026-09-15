-- name: CreateDocument :one
INSERT INTO documents (
    clerk_user_id, original_filename, input_object_key, status
) VALUES (
    $1, $2, $3, 'uploaded'
)
RETURNING id, clerk_user_id, original_filename, input_object_key, status, created_at, updated_at;

-- name: GetDocument :one
SELECT id, clerk_user_id, original_filename, input_object_key, output_object_key, page_count, status, created_at, updated_at 
FROM documents 
WHERE id = $1 AND clerk_user_id = $2;

-- name: ListDocuments :many
SELECT id, clerk_user_id, original_filename, input_object_key, output_object_key, page_count, status, created_at, updated_at 
FROM documents 
WHERE clerk_user_id = $1
ORDER BY created_at DESC;

-- name: DeleteDocument :exec
DELETE FROM documents
WHERE id = $1 AND clerk_user_id = $2;

-- name: EnqueueJob :one
INSERT INTO jobs (
    document_id, type, status
) VALUES (
    $1, $2, 'pending'
)
RETURNING id, document_id, type, status, created_at;
