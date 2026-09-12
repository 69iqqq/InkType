"use client";

import React, { useState, useRef } from "react";

export function DownloadButton() {
  return (
    <button className="px-4 py-1.5 text-sm font-medium bg-foreground text-background hover:opacity-90 rounded transition-opacity">
      Download PDF
    </button>
  );
}

export function DocumentHeader({ title, onReset }: { title: string; onReset: () => void }) {
  return (
    <div className="flex items-center justify-between px-8 py-4 border-b border-border flex-shrink-0">
      <div className="flex items-center gap-4">
        <button className="md:hidden px-2 py-1 -ml-2 text-secondary hover:text-foreground hover:bg-hover rounded transition-colors text-sm">
          Menu
        </button>
        <h1 className="text-sm font-bold">{title}</h1>
      </div>
      <div className="flex items-center gap-4">
        <button onClick={onReset} className="text-sm text-secondary hover:text-foreground transition-colors">
          Reset
        </button>
        <button className="text-sm text-secondary hover:text-foreground transition-colors">
          Options
        </button>
        <DownloadButton />
      </div>
    </div>
  );
}

export function ViewerToolbar({ pages = 12, currentPage = 1, zoom = 100 }) {
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
  content, 
  isGenerated 
}: { 
  title: string; 
  content: React.ReactNode; 
  isGenerated?: boolean;
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

export function SplitDocumentViewer() {
  const originalContent = (
    <div className="opacity-50 blur-[0.5px] select-none pointer-events-none grayscale flex flex-col gap-4">
      <div className="w-3/4 h-8 bg-border rounded-sm mb-4"></div>
      <div className="w-full h-4 bg-border rounded-sm"></div>
      <div className="w-full h-4 bg-border rounded-sm"></div>
      <div className="w-5/6 h-4 bg-border rounded-sm"></div>
      <div className="w-1/2 h-16 bg-border rounded-sm my-4"></div>
      <div className="w-full h-4 bg-border rounded-sm"></div>
      <div className="w-4/5 h-4 bg-border rounded-sm"></div>
    </div>
  );

  const generatedContent = (
    <div className="text-foreground">
      <h1 className="text-2xl font-bold mb-8">foo1</h1>
      <p className="mb-6 leading-relaxed text-sm">
        foo2
      </p>
      <p className="mb-6 leading-relaxed text-sm">
        foo3
      </p>
    </div>
  );

  return (
    <div className="flex flex-1 overflow-hidden">
      <DocumentViewer title="Original PDF" content={originalContent} />
      <DocumentViewer title="Typeset PDF" content={generatedContent} isGenerated />
    </div>
  );
}

export function ProcessingView({ onComplete }: { onComplete: () => void }) {
  const [step, setStep] = React.useState(0);
  
  React.useEffect(() => {
    const timer1 = setTimeout(() => setStep(1), 1500);
    const timer2 = setTimeout(() => setStep(2), 3000);
    const timer3 = setTimeout(() => setStep(3), 4500);
    const timer4 = setTimeout(() => onComplete(), 6000);
    
    return () => {
      clearTimeout(timer1);
      clearTimeout(timer2);
      clearTimeout(timer3);
      clearTimeout(timer4);
    };
  }, [onComplete]);

  const steps = [
    "Reading handwriting",
    "Understanding equations",
    "Reconstructing layout",
    "Typesetting document"
  ];

  return (
    <div className="flex-1 flex flex-col items-center justify-center p-8 bg-background">
      <div className="w-full max-w-sm">
        <h2 className="text-xl font-bold mb-12 text-center">Preparing document...</h2>
        
        <div className="flex flex-col gap-6">
          {steps.map((label, index) => {
            const isCompleted = step > index;
            const isCurrent = step === index;
            const isPending = step < index;
            
            return (
              <div key={label} className={`flex items-center justify-between transition-opacity duration-500 ${isPending ? 'opacity-30' : 'opacity-100'}`}>
                <span className={`text-sm ${isCurrent ? 'text-foreground font-bold' : 'text-secondary'}`}>
                  {label}
                </span>
                <div className="text-sm font-bold w-12 text-right">
                  {isCompleted ? (
                    <span className="text-foreground">Done</span>
                  ) : isCurrent ? (
                    <span className="text-foreground animate-pulse">Wait</span>
                  ) : (
                    <span className="text-secondary">-</span>
                  )}
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}

export function UploadArea({ onUpload }: { onUpload: () => void }) {
  const [isDragging, setIsDragging] = useState(false);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
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
    if (file.type === "application/pdf" || file.name.endsWith('.pdf')) {
      setSelectedFile(file);
    } else {
      setSelectedFile(new File([""], "foo1.pdf", { type: "application/pdf" }));
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
          <p className="text-secondary text-sm mt-2">24 pages · 8.2 MB</p>
        </div>
        <button 
          onClick={onUpload}
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
        className={`w-72 h-48 border flex flex-col items-center justify-center gap-4 cursor-pointer transition-colors ${
          isDragging 
            ? 'border-foreground bg-sidebar border-solid' 
            : 'border-border border-dashed hover:border-foreground hover:bg-sidebar'
        }`}
      >
        <div className="text-center">
          <p className="font-bold mb-2">Drop PDF here</p>
          <p className="text-sm text-secondary underline">or browse files</p>
        </div>
        <input 
          type="file" 
          ref={fileInputRef}
          onChange={handleChange}
          accept="application/pdf"
          className="hidden" 
        />
      </div>
      <p className="text-xs font-bold text-secondary mt-6 uppercase tracking-widest">Max 50 MB</p>
    </div>
  );
}

export function EmptyState({ onUpload }: { onUpload: () => void }) {
  return (
    <div className="flex-1 flex flex-col items-center justify-center p-8 bg-background">
      <div className="max-w-xl w-full flex flex-col items-center text-center">
        <h2 className="text-3xl md:text-4xl font-bold tracking-tight mb-6">Turn handwriting<br />into something clean.</h2>
        <p className="text-secondary text-lg mb-10 max-w-sm">
          Upload a handwritten PDF and we'll turn it into beautiful, typeset notes.
        </p>
        
        <UploadArea onUpload={onUpload} />
      </div>
    </div>
  );
}

export type AppState = 'empty' | 'processing' | 'result';

export function Workspace() {
  const [appState, setAppState] = useState<AppState>('empty');

  const handleUpload = () => {
    setAppState('processing');
  };

  const handleProcessComplete = () => {
    setAppState('result');
  };

  const handleReset = () => {
    setAppState('empty');
  };

  if (appState === 'empty') {
    return <EmptyState onUpload={handleUpload} />;
  }

  if (appState === 'processing') {
    return <ProcessingView onComplete={handleProcessComplete} />;
  }

  return (
    <div className="flex-1 flex flex-col h-full overflow-hidden bg-background animate-in fade-in duration-500">
      <DocumentHeader title="foo1" onReset={handleReset} />
      <SplitDocumentViewer />
    </div>
  );
}
