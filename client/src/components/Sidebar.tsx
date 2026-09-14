"use client";

import React from "react";
import Image from "next/image";

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

export function HistoryItem({ title, active }: { title: string; active?: boolean }) {
  return (
    <div
      className={`group flex items-center justify-between py-1.5 px-3 cursor-pointer rounded text-sm transition-colors ${
        active ? "bg-hover font-medium text-foreground" : "text-secondary hover:bg-hover hover:text-foreground"
      }`}
    >
      <span className="truncate pr-2">{title}</span>
      <button className="opacity-0 group-hover:opacity-100 transition-opacity text-xs tracking-widest text-secondary hover:text-foreground">
        ...
      </button>
    </div>
  );
}

export function HistoryList() {
  return (
    <div className="flex flex-col gap-6 mt-6">
      <div>
        <div className="text-xs font-semibold text-secondary mb-2 px-3 uppercase tracking-wider">
          Today
        </div>
        <div className="flex flex-col gap-0.5">
          <HistoryItem title="foo1" active />
        </div>
      </div>
      <div>
        <div className="text-xs font-semibold text-secondary mb-2 px-3 uppercase tracking-wider">
          Past History
        </div>
        <div className="flex flex-col gap-0.5">
          <HistoryItem title="foo2" />
          <HistoryItem title="foo3" />
        </div>
      </div>
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
      <div className="w-5 h-5 rounded-full bg-border flex items-center justify-center text-xs">
        F
      </div>
      foo1
    </div>
  );
}

export function Sidebar({ onNewDocument }: { onNewDocument?: () => void }) {
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
        <HistoryList />
      </div>

      <div className="pt-4 border-t border-border mt-auto">
        <SettingsButton />
        <UserMenu />
      </div>
    </div>
  );
}
