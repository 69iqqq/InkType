"use client";

import { Sidebar } from "@/components/Sidebar";
import { Workspace } from "@/components/Workspace";

export default function Home() {
  return (
    <main className="flex h-full w-full overflow-hidden">
      <Sidebar />
      <Workspace />
    </main>
  );
}
