"use client";

import React, { useState, useEffect } from "react";
import Image from "next/image";
import { useAuth, UserButton } from "@clerk/nextjs";
import { deleteDocument } from "../lib/api";

export function NewDocumentButton({ onClick }: { onClick?: () => void }) {
  return (
    <button
      onClick={onClick}
      className="w-full text-left py-2 px-3 text-sm font-medium text-foreground hover:bg-hover rounded transition-colors border border-border"
    >
      + New document
    </button>
  );
}

export function HistoryItem({ 
  title, 
  active, 
  onClick,
  onDelete,
}: { 
  title: string; 
  active?: boolean; 
  onClick?: () => void;
  onDelete?: () => void;
}) {
  return (
    <div
      onClick={onClick}
      className={`group flex items-center justify-between py-1.5 px-3 cursor-pointer rounded text-sm transition-colors ${
        active ? "bg-hover font-medium text-foreground" : "text-secondary hover:bg-hover hover:text-foreground"
      }`}
    >
      <span className="truncate pr-2">{title}</span>
      {onDelete && (
        <button 
          onClick={(e) => {
            e.stopPropagation();
            onDelete();
          }}
          title="Delete document"
          className="opacity-0 group-hover:opacity-100 hover:text-red-500 hover:bg-red-500/10 rounded px-1.5 py-0.5 text-xs transition-all text-secondary"
        >
          ✕
        </button>
      )}
    </div>
  );
}

export function HistoryList({ 
  selectedDocId, 
  onSelectDocument,
  onDeleteDocument,
  refreshTrigger,
}: { 
  selectedDocId?: string | null; 
  onSelectDocument?: (id: string) => void;
  onDeleteDocument?: (id: string) => void;
  refreshTrigger?: number;
}) {
  const { isLoaded, isSignedIn, getToken } = useAuth();
  const [docs, setDocs] = useState<{ id: string; original_filename: string; created_at: string; status: string }[]>([]);

  const fetchDocs = async () => {
    if (!isLoaded || !isSignedIn) return;
    try {
      const token = await getToken();
      if (!token) return;
      const apiBase = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";
      const res = await fetch(`${apiBase}/documents`, {
        headers: { "Authorization": `Bearer ${token}` }
      });
      if (res.ok) {
        const data = await res.json();
        setDocs(Array.isArray(data) ? data : []);
      }
    } catch (e) {
      console.error('Failed to fetch documents:', e);
    }
  };

  useEffect(() => {
    if (isLoaded && isSignedIn) {
      fetchDocs();
    }
  }, [isLoaded, isSignedIn, refreshTrigger]);

  const handleDelete = async (id: string, name: string) => {
    if (!confirm(`Delete "${name || 'Untitled'}" from database?`)) return;
    try {
      const success = await deleteDocument(id);
      if (success) {
        setDocs(prev => prev.filter(d => d.id !== id));
        onDeleteDocument?.(id);
      } else {
        alert("Failed to delete document from database.");
      }
    } catch (err) {
      console.error("Error deleting document:", err);
    }
  };

  if (docs.length === 0) {
    return (
      <div className="px-3 py-4 text-xs text-secondary">
        No documents yet
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-0.5">
      {docs.map((doc) => (
        <HistoryItem 
          key={doc.id} 
          title={doc.original_filename || 'Untitled'} 
          active={doc.id === selectedDocId}
          onClick={() => onSelectDocument?.(doc.id)}
          onDelete={() => handleDelete(doc.id, doc.original_filename)}
        />
      ))}
    </div>
  );
}

export function SettingsButton() {
  return (
    <button className="w-full text-left py-1.5 px-3 text-sm text-secondary hover:text-foreground hover:bg-hover rounded transition-colors mb-1">
      Settings
    </button>
  );
}

export function UserMenu() {
  return (
    <div className="flex items-center gap-2 py-2 px-3 text-sm font-medium text-foreground hover:bg-hover rounded cursor-pointer transition-colors">
      <UserButton />
    </div>
  );
}

export function Sidebar({ 
  selectedDocId,
  onSelectDocument,
  onDeleteDocument,
  onNewDocument,
  refreshTrigger,
}: { 
  selectedDocId?: string | null;
  onSelectDocument?: (id: string) => void;
  onDeleteDocument?: (id: string) => void;
  onNewDocument?: () => void;
  refreshTrigger?: number;
}) {
  return (
    <div className="hidden md:flex flex-col w-[260px] h-full bg-sidebar border-r border-border p-3 flex-shrink-0">
      <div className="px-3 py-3 mb-4 flex items-center">
        <Image src="/InkTypeLogo-v2.png" alt="InkType" width={120} height={40} className="h-8 w-auto object-contain" priority />
      </div>
      <NewDocumentButton onClick={onNewDocument} />
      
      <div className="flex-1 overflow-y-auto mt-6 scrollbar-hide">
        <div className="text-xs font-semibold text-secondary mb-2 px-3 uppercase tracking-wider">
          History
        </div>
        <HistoryList 
          selectedDocId={selectedDocId} 
          onSelectDocument={onSelectDocument}
          onDeleteDocument={onDeleteDocument}
          refreshTrigger={refreshTrigger}
        />
      </div>

      <div className="pt-4 border-t border-border mt-auto">
        <SettingsButton />
        <UserMenu />
      </div>
    </div>
  );
}
