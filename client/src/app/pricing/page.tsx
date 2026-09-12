"use client";

import Link from "next/link";

export default function PricingPage() {
  return (
    <main className="min-h-full w-full flex flex-col bg-background text-foreground overflow-y-auto overflow-x-hidden">
      <nav className="flex items-center justify-between px-6 py-4 md:px-12 border-b border-border bg-background/80 backdrop-blur-md sticky top-0 z-50">
        <Link href="/" className="font-bold tracking-tight text-xl">
          INKTYPE
        </Link>
        <Link href="/" className="text-sm font-medium text-secondary hover:text-foreground transition-colors">
          back
        </Link>
      </nav>
      
      <section className="px-6 py-24 md:py-32 md:px-12 max-w-4xl mx-auto w-full">
        <h1 className="text-4xl md:text-6xl font-bold tracking-tighter mb-12 text-center">
          Pricing
        </h1>
        <div className="grid md:grid-cols-2 gap-8">
          <div className="border border-border p-8 rounded flex flex-col">
            <h2 className="text-2xl font-bold mb-4">Free</h2>
            <p className="text-secondary mb-8">Perfect for occasional use.</p>
            <div className="text-4xl font-bold mb-8">$0<span className="text-lg text-secondary font-normal">/mo</span></div>
            <ul className="flex flex-col gap-4 mb-8 flex-1">
              <li className="flex items-center gap-3"><span className="text-secondary">+</span> 5 documents per month</li>
              <li className="flex items-center gap-3"><span className="text-secondary">+</span> Standard processing speed</li>
              <li className="flex items-center gap-3"><span className="text-secondary">+</span> Basic exports</li>
            </ul>
            <button className="w-full py-3 bg-sidebar border border-border font-medium hover:bg-hover transition-colors rounded">
              Current Plan
            </button>
          </div>
          <div className="border-2 border-foreground p-8 rounded flex flex-col bg-sidebar">
            <h2 className="text-2xl font-bold mb-4">Pro</h2>
            <p className="text-secondary mb-8">For students and professionals.</p>
            <div className="text-4xl font-bold mb-8">$10<span className="text-lg text-secondary font-normal">/mo</span></div>
            <ul className="flex flex-col gap-4 mb-8 flex-1">
              <li className="flex items-center gap-3"><span className="text-foreground">+</span> Unlimited documents</li>
              <li className="flex items-center gap-3"><span className="text-foreground">+</span> Priority processing</li>
              <li className="flex items-center gap-3"><span className="text-foreground">+</span> Advanced math recognition</li>
              <li className="flex items-center gap-3"><span className="text-foreground">+</span> Custom fonts</li>
            </ul>
            <button className="w-full py-3 bg-foreground text-background font-medium hover:opacity-90 transition-opacity rounded">
              Upgrade to Pro
            </button>
          </div>
        </div>
      </section>
    </main>
  );
}
