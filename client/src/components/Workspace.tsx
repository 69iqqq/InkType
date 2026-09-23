"use client";

import React, { useState, useRef, useEffect } from "react";
import { uploadDocument, getAuthToken, deleteDocument } from "../lib/api";

export function DownloadDropdown({ documentInfo }: { documentInfo?: any }) {
  const [open, setOpen] = useState(false);

  const handleDownload = async (format: 'pdf' | 'md') => {
    const apiBase = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";
    const docId = documentInfo?.id;
    if (!docId) {
        console.log('No document ID for download');
        setOpen(false);
        return;
    }
    const url = `${apiBase}/documents/${docId}/download?format=${format}`;
    try {
        const { getAuthToken } = await import('../lib/api');
        const token = await getAuthToken();
        const res = await fetch(url, { headers: { 'Authorization': `Bearer ${token}` } });
        if (!res.ok) throw new Error('Download failed');
        const blob = await res.blob();
        const a = window.document.createElement('a');
        a.href = URL.createObjectURL(blob);
        a.download = `document.${format === 'md' ? 'md' : 'pdf'}`;
        a.click();
        URL.revokeObjectURL(a.href);
    } catch (e) {
        console.error('Download error:', e);
    }
    setOpen(false);
  };

  return (
    <div className="relative">
      <button 
        onClick={() => setOpen(!open)}
        className="flex items-center px-4 py-2 text-sm font-bold bg-foreground text-background hover:bg-transparent hover:text-foreground border-2 border-foreground transition-colors"
      >
        <span>Export</span>
        <span className="ml-4 text-xs font-normal opacity-50">▼</span>
      </button>
      
      {open && (
        <div className="absolute right-0 mt-2 w-48 bg-background border border-border shadow-lg z-50">
          <button 
            onClick={() => handleDownload('pdf')}
            className="w-full text-left px-4 py-2 text-sm hover:bg-hover transition-colors"
          >
            Download PDF
          </button>
          <button 
            onClick={() => handleDownload('md')}
            className="w-full text-left px-4 py-2 text-sm hover:bg-hover transition-colors border-t border-border"
          >
            Download Markdown
          </button>
        </div>
      )}
    </div>
  );
}

export function DocumentHeader({ 
  title, 
  onReset, 
  onDelete, 
  documentInfo 
}: { 
  title: string; 
  onReset: () => void; 
  onDelete?: () => void; 
  documentInfo?: any; 
}) {
  return (
    <div className="flex items-center justify-between px-8 py-4 border-b border-border flex-shrink-0">
      <div className="flex items-center gap-4">
        <button className="md:hidden px-2 py-1 -ml-2 text-secondary hover:text-foreground hover:bg-hover rounded transition-colors text-sm">
          Menu
        </button>
        <h1 className="text-sm font-bold">{title}</h1>
      </div>
      <div className="flex items-center gap-4">
        {onDelete && (
          <button 
            onClick={onDelete} 
            className="text-sm text-red-500 hover:text-red-700 hover:bg-red-500/10 px-3 py-1.5 rounded transition-colors font-medium"
          >
            Delete
          </button>
        )}
        <button onClick={onReset} className="text-sm text-secondary hover:text-foreground transition-colors">
          Reset
        </button>
        <DownloadDropdown documentInfo={documentInfo} />
      </div>
    </div>
  );
}

export function ViewerToolbar({ pages = 1, currentPage = 1, zoom = 100 }) {
  return (
    <div className="flex items-center justify-between px-6 py-3 border-t border-border bg-background text-sm text-secondary flex-shrink-0">
      <div className="flex items-center gap-4">
        <button className="px-2 py-1 hover:text-foreground hover:bg-hover rounded transition-colors">&lt; Prev</button>
        <span>{currentPage} / {pages}</span>
        <button className="px-2 py-1 hover:text-foreground hover:bg-hover rounded transition-colors">Next &gt;</button>
      </div>
      <div className="flex items-center gap-4">
        <button className="px-2 py-1 hover:text-foreground hover:bg-hover rounded transition-colors">-</button>
        <span>{zoom}%</span>
        <button className="px-2 py-1 hover:text-foreground hover:bg-hover rounded transition-colors">+</button>
      </div>
    </div>
  );
}

export function DocumentViewer({ 
  title, 
  content
}: { 
  title: string; 
  content: React.ReactNode;
}) {
  return (
    <div className="flex flex-col h-full bg-sidebar border-r border-border last:border-r-0 w-full relative">
      <div className="absolute top-0 left-0 right-0 p-4 flex justify-center z-10 pointer-events-none">
        <span className="text-xs font-bold text-secondary uppercase tracking-widest bg-sidebar/90 px-3 py-1 rounded backdrop-blur-sm border border-border">
          {title}
        </span>
      </div>
      <div className="flex-1 overflow-auto p-12 pt-20 flex justify-center">
        <div className="w-full max-w-[500px] h-fit bg-background border border-border shadow-sm rounded-sm p-10 min-h-[700px]">
          {content}
        </div>
      </div>
      <ViewerToolbar />
    </div>
  );
}

export function SplitDocumentViewer({ document }: { document?: any }) {
  const originalContent = (
    <div className="flex flex-col items-center justify-center gap-4 h-full">
        <p className="text-secondary text-sm">Original document</p>
        <p className="text-foreground font-bold">{document?.original_filename || 'Untitled'}</p>
        <p className="text-xs text-secondary">Status: {document?.status || 'unknown'}</p>
        {document?.page_count && (
            <p className="text-xs text-secondary">{document.page_count} pages</p>
        )}
    </div>
  );

  const blocks = (document?.pages || []).flatMap((page: { blocks?: { type: string; content?: { text?: string }; confidence?: number }[] }) => 
    (page.blocks || []).map(b => ({ ...b, text: b.content?.text || '' }))
  );

  const generatedContent = (
    <div className="text-foreground">
      {blocks.length > 0 ? (
        blocks.map((block: any, idx: number) => {
          let content = block.text || "";
          
          if (block.confidence !== undefined && block.confidence < 0.8) {
            content = <span className="bg-yellow-500/20">{content}</span>;
          }

          if (block.type === 'heading' || block.type === 'h1') {
            return <h1 key={idx} className="text-2xl font-bold mb-8">{content}</h1>;
          } else if (block.type === 'h2') {
            return <h2 key={idx} className="text-xl font-bold mb-6">{content}</h2>;
          } else if (block.type === 'h3') {
            return <h3 key={idx} className="text-lg font-bold mb-4">{content}</h3>;
          }
          return <p key={idx} className="mb-6 leading-relaxed text-sm">{content}</p>;
        })
      ) : (
        <>
          <h1 className="text-2xl font-bold mb-8">Waiting for content...</h1>
          <p className="mb-6 leading-relaxed text-sm text-secondary">
            No structured text blocks found.
          </p>
        </>
      )}
    </div>
  );

  return (
    <div className="flex flex-1 overflow-hidden">
      <DocumentViewer title="Original PDF" content={originalContent} />
      <DocumentViewer title="Typeset PDF" content={generatedContent} />
    </div>
  );
}

export function ProcessingView({ message }: { message?: string }) {
  return (
    <div className="flex-1 flex flex-col items-center justify-center p-8 bg-background animate-in fade-in duration-300">
      <div className="w-full max-w-sm flex flex-col items-center">
        <div className="w-12 h-12 rounded-full border-4 border-border border-t-foreground animate-spin mb-8"></div>
        <h2 className="text-xl font-bold mb-4 text-center">Processing document...</h2>
        <p className="text-secondary text-sm text-center mb-12">
          {message || "This may take a few moments depending on the length of your document."}
        </p>
      </div>
    </div>
  );
}

export function UploadArea({ onUpload }: { onUpload: (file: File) => void }) {
  const [isDragging, setIsDragging] = useState(false);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [uploadError, setUploadError] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = () => {
    setIsDragging(false);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
    if (e.dataTransfer.files && e.dataTransfer.files[0]) {
      handleFileSelection(e.dataTransfer.files[0]);
    }
  };

  const handleFileSelection = (file: File) => {
    setUploadError(null);
    if (file.size > 50 * 1024 * 1024) {
      setUploadError('File exceeds 50 MB limit.');
      return;
    }
    if (file.type === "application/pdf" || file.name.endsWith('.pdf')) {
      setSelectedFile(file);
    } else {
      setUploadError('Only PDF files are supported.');
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      handleFileSelection(e.target.files[0]);
    }
  };

  if (selectedFile) {
    return (
      <div className="flex flex-col items-center gap-6 mt-12 animate-in fade-in duration-300">
        <div className="text-center p-6 border border-border bg-sidebar rounded">
          <p className="text-foreground font-bold">{selectedFile.name}</p>
          <p className="text-secondary text-sm mt-2">Ready to process</p>
        </div>
        <button 
          onClick={() => onUpload(selectedFile)}
          className="px-6 py-3 bg-foreground text-background font-medium rounded hover:opacity-90 transition-opacity"
        >
          Process Document
        </button>
      </div>
    );
  }

  return (
    <div className="mt-12 flex flex-col items-center">
      <div 
        onClick={() => fileInputRef.current?.click()}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
        className={`w-full max-w-lg p-16 border-2 flex flex-col items-center justify-center gap-4 cursor-pointer transition-colors ${
          isDragging 
            ? 'border-foreground bg-sidebar border-solid' 
            : 'border-border border-dashed hover:border-foreground hover:bg-sidebar'
        }`}
      >
        <div className="text-center">
          <p className="font-bold text-xl mb-4">Drop PDF here</p>
          <p className="text-sm text-secondary uppercase tracking-widest"><span className="underline">or browse files</span> &nbsp; [ ⌘ U ]</p>
        </div>
        <input 
          type="file" 
          ref={fileInputRef}
          onChange={handleChange}
          accept="application/pdf"
          className="hidden" 
        />
      </div>
      {uploadError && (
        <p className="text-xs font-bold text-red-500 mt-4">{uploadError}</p>
      )}
      <p className="text-xs font-bold text-secondary mt-6 uppercase tracking-widest">Max 50 MB</p>
    </div>
  );
}

export function EmptyState({ onUpload }: { onUpload: (file: File) => void }) {
  return (
    <div className="flex-1 flex flex-col items-center justify-center p-8 bg-background">
      <div className="max-w-xl w-full flex flex-col items-center text-center">
        <h2 className="text-3xl md:text-4xl font-bold tracking-tight mb-6">Turn handwriting<br />into something clean.</h2>
        <p className="text-secondary text-lg mb-10 max-w-sm">
          Upload a handwritten PDF and we&apos;ll turn it into beautiful, typeset notes.
        </p>
        
        <UploadArea onUpload={onUpload} />
      </div>
    </div>
  );
}

export type AppState = 'empty' | 'processing' | 'result';

export function Workspace({ 
  selectedDocId, 
  onResetSelected,
  onDeleteDocument,
}: { 
  selectedDocId?: string | null; 
  onResetSelected?: () => void; 
  onDeleteDocument?: (id: string) => void;
} = {}) {
  const [appState, setAppState] = useState<AppState>('empty');
  const [documentInfo, setDocumentInfo] = useState<any>(null);
  const [progressMsg, setProgressMsg] = useState<string>("");

  useEffect(() => {
    if (!selectedDocId) {
      return;
    }
    const fetchSelectedDoc = async () => {
      try {
        const token = await getAuthToken();
        const apiBase = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";
        const res = await fetch(`${apiBase}/documents/${selectedDocId}`, {
          headers: { "Authorization": `Bearer ${token}` }
        });
        if (res.ok) {
          const doc = await res.json();
          setDocumentInfo(doc);
          if (doc.status === 'completed') {
            setAppState('result');
          } else if (doc.status === 'processing') {
            setAppState('processing');
          } else {
            setAppState('result');
          }
        }
      } catch (err) {
        console.error("Failed to fetch selected document:", err);
      }
    };
    fetchSelectedDoc();
  }, [selectedDocId]);

  const handleUpload = async (file: File) => {
    setAppState('processing');
    setProgressMsg("");
    try {
      const doc = await uploadDocument(file);
      console.log("Uploaded file to R2, DB record created:", doc);
      setDocumentInfo(doc);
    } catch (e) {
      console.error(e);
      setAppState('empty');
    }
  };

  useEffect(() => {
    if (appState !== 'processing' || !documentInfo || !documentInfo.id) return;

    let ws: WebSocket;
    const connectWs = async () => {
      try {
        const token = await getAuthToken();
        const apiBase = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";
        // Convert http:// to ws://
        const wsBase = apiBase.replace(/^http/, 'ws');
        
        ws = new WebSocket(`${wsBase}/documents/${documentInfo.id}/ws?token=${token}`);

        ws.onmessage = async (event) => {
          try {
            const data = JSON.parse(event.data);
            
            if (data.status === 'completed') {
              // Fetch the fully complete document which should contain the structured content
              try {
                const res = await fetch(`${apiBase}/documents/${documentInfo.id}`, {
                  headers: { "Authorization": `Bearer ${token}` }
                });
                if (res.ok) {
                  const finalDoc = await res.json();
                  setDocumentInfo(finalDoc);
                } else {
                  setDocumentInfo((prev: any) => ({ ...prev, status: 'completed' }));
                }
              } catch (e) {
                setDocumentInfo((prev: any) => ({ ...prev, status: 'completed' }));
              }
              setAppState('result');
            } else if (data.status === 'failed' || data.error) {
              console.error("Document processing failed:", data.error || data);
              setAppState('empty');
            } else {
              // Active processing
              const total = data.pages_total || '?';
              const completed = data.pages_completed || 0;
              setProgressMsg(`Processing page ${completed} of ${total}`);
            }
          } catch (err) {
            console.error("Failed to parse WS message", err);
          }
        };

        ws.onerror = (err) => {
          console.error("WebSocket error:", err);
        };
      } catch (err) {
        console.error("Error setting up WS:", err);
      }
    };

    connectWs();
    return () => {
      if (ws) {
        ws.close();
      }
    };
  }, [appState, documentInfo?.id]);

  const handleReset = () => {
    setAppState('empty');
    setDocumentInfo(null);
    setProgressMsg("");
    onResetSelected?.();
  };

  const handleDeleteCurrent = async () => {
    if (!documentInfo?.id) return;
    if (!confirm(`Delete "${documentInfo.original_filename || 'this document'}" from database?`)) return;
    try {
      const success = await deleteDocument(documentInfo.id);
      if (success) {
        const deletedId = documentInfo.id;
        handleReset();
        onDeleteDocument?.(deletedId);
      } else {
        alert("Failed to delete document from database.");
      }
    } catch (err) {
      console.error("Failed to delete document:", err);
    }
  };

  if (appState === 'empty') {
    return <EmptyState onUpload={handleUpload} />;
  }

  if (appState === 'processing') {
    return <ProcessingView message={progressMsg} />;
  }

  return (
    <div className="flex-1 flex flex-col h-full overflow-hidden bg-background animate-in fade-in duration-500">
      <DocumentHeader 
        title={documentInfo?.original_filename || "Untitled Document"} 
        onReset={handleReset} 
        onDelete={handleDeleteCurrent}
        documentInfo={documentInfo}
      />
      <SplitDocumentViewer document={documentInfo} />
    </div>
  );
}
