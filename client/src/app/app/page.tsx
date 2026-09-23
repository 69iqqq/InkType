"use client";

import { useState } from "react";
import { Sidebar } from "@/components/Sidebar";
import { Workspace } from "@/components/Workspace";

export default function Home() {
  const [selectedDocId, setSelectedDocId] = useState<string | null>(null);
  const [refreshTrigger, setRefreshTrigger] = useState<number>(0);

  const handleDelete = (id: string) => {
    if (selectedDocId === id) {
      setSelectedDocId(null);
    }
    setRefreshTrigger(prev => prev + 1);
  };

  return (
    <main className="flex h-full w-full overflow-hidden">
      <Sidebar 
        selectedDocId={selectedDocId}
        onSelectDocument={(id) => setSelectedDocId(id)}
        onDeleteDocument={handleDelete}
        onNewDocument={() => setSelectedDocId(null)}
        refreshTrigger={refreshTrigger}
      />
      <Workspace 
        selectedDocId={selectedDocId}
        onResetSelected={() => setSelectedDocId(null)}
        onDeleteDocument={handleDelete}
      />
    </main>
  );
}
