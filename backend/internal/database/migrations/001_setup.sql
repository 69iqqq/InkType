-- Initial schema setup for InkType

CREATE TABLE documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clerk_user_id VARCHAR(255) NOT NULL,
    original_filename TEXT NOT NULL,
    input_object_key TEXT NOT NULL,
    output_object_key TEXT,
    page_count INTEGER,
    status VARCHAR(50) NOT NULL DEFAULT 'uploaded',
    error TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE document_pages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    page_number INTEGER NOT NULL,
    image_object_key TEXT NOT NULL,
    extracted_json JSONB,
    ocr_result JSONB,
    embedded_text TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    error TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(document_id, page_number)
);

CREATE TABLE jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    document_page_id UUID REFERENCES document_pages(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    attempts INTEGER NOT NULL DEFAULT 0,
    available_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    error TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_documents_user_created ON documents (clerk_user_id, created_at DESC);
CREATE INDEX idx_document_pages_doc_id ON document_pages (document_id);
CREATE INDEX idx_document_pages_pending ON document_pages (document_id, status) WHERE status = 'pending';
CREATE INDEX idx_jobs_document_id ON jobs (document_id);
CREATE INDEX idx_jobs_claimable ON jobs (status, available_at, created_at) WHERE status IN ('pending', 'processing');

-- Auto-update updated_at timestamps
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_documents_updated_at BEFORE UPDATE ON documents FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_document_pages_updated_at BEFORE UPDATE ON document_pages FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_jobs_updated_at BEFORE UPDATE ON jobs FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

---- create above / drop below ----

DROP TRIGGER IF EXISTS update_jobs_updated_at ON jobs;
DROP TRIGGER IF EXISTS update_document_pages_updated_at ON document_pages;
DROP TRIGGER IF EXISTS update_documents_updated_at ON documents;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP TABLE IF EXISTS jobs;
DROP TABLE IF EXISTS document_pages;
DROP TABLE IF EXISTS documents;
