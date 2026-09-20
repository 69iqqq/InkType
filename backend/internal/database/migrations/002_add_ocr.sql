ALTER TABLE document_pages
ADD COLUMN ocr_result JSONB,
ADD COLUMN embedded_text TEXT;
