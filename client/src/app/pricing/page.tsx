"use client";

import Link from "next/link";
import Image from "next/image";
import { useState } from "react";

export default function PricingPage() {
  const [openFaq, setOpenFaq] = useState<number | null>(null);

  const faqs = [
    { q: "Can I cancel anytime?", a: "Yes. No contracts, no friction." },
    { q: "What counts as a document?", a: "Any PDF up to 50 pages. Larger files require chunking." },
    { q: "Do you offer educational discounts?", a: "Students get 50% off the Pro plan with a valid .edu email." }
  ];

  return (
    <main className="min-h-full w-full flex flex-col bg-background text-foreground overflow-y-auto overflow-x-hidden">
      <nav className="flex items-center justify-between px-6 py-6 md:px-12 border-b-2 border-border sticky top-0 z-50 bg-background/90 backdrop-blur">
        <Link href="/">
          <Image src="/InkTypeLogo-v2.png" alt="InkType" width={120} height={40} className="h-8 w-auto object-contain" priority />
        </Link>
        <Link href="/" className="text-sm font-bold uppercase tracking-widest text-secondary hover:text-foreground transition-colors">
          &lt;- Back
        </Link>
      </nav>
      
      <section className="px-6 py-24 md:py-32 md:px-12 max-w-4xl mx-auto w-full">
        <div className="flex flex-col items-center mb-24 text-center">
          <h1 className="text-5xl md:text-7xl font-bold tracking-tighter mb-6 uppercase">
            Service Menu
          </h1>
          <p className="text-xl text-secondary max-w-lg leading-relaxed">
            Simple, predictable pricing. Built for heavy usage and deep work.
          </p>
        </div>

        <div className="grid md:grid-cols-2 gap-0 border-2 border-foreground mb-32">
          {/* Free Plan */}
          <div className="p-12 md:border-r-2 border-border flex flex-col bg-background">
            <h2 className="text-3xl font-bold mb-2 uppercase tracking-tight">Free</h2>
            <p className="text-secondary mb-12 text-sm font-bold tracking-widest uppercase">Casual</p>
            
            <div className="text-6xl font-bold mb-12 tracking-tighter">$0</div>
            
            <ul className="flex flex-col gap-6 mb-12 flex-1 text-sm font-bold">
              <li className="flex justify-between border-b border-border pb-4">
                <span>Monthly Limit</span>
                <span className="text-secondary">5 Docs</span>
              </li>
              <li className="flex justify-between border-b border-border pb-4">
                <span>Speed</span>
                <span className="text-secondary">Standard</span>
              </li>
              <li className="flex justify-between border-b border-border pb-4">
                <span>Export Formats</span>
                <span className="text-secondary">PDF</span>
              </li>
            </ul>
            
            <button className="w-full py-4 bg-sidebar border-2 border-border font-bold hover:bg-hover transition-colors">
              Current Plan
            </button>
          </div>

          {/* Pro Plan */}
          <div className="p-12 flex flex-col bg-sidebar">
            <h2 className="text-3xl font-bold mb-2 uppercase tracking-tight">Pro</h2>
            <p className="text-secondary mb-12 text-sm font-bold tracking-widest uppercase">Professional</p>
            
            <div className="text-6xl font-bold mb-12 tracking-tighter">$10<span className="text-xl text-secondary font-bold tracking-normal">/mo</span></div>
            
            <ul className="flex flex-col gap-6 mb-12 flex-1 text-sm font-bold">
              <li className="flex justify-between border-b border-foreground pb-4">
                <span>Monthly Limit</span>
                <span className="text-foreground">Unlimited</span>
              </li>
              <li className="flex justify-between border-b border-foreground pb-4">
                <span>Speed</span>
                <span className="text-foreground">Priority</span>
              </li>
              <li className="flex justify-between border-b border-foreground pb-4">
                <span>Math Recognition</span>
                <span className="text-foreground">Advanced</span>
              </li>
              <li className="flex justify-between border-b border-foreground pb-4">
                <span>Custom Fonts</span>
                <span className="text-foreground">Included</span>
              </li>
            </ul>
            
            <Link href="/checkout" className="text-center w-full py-4 bg-foreground text-background font-bold hover:bg-transparent hover:text-foreground border-2 border-foreground transition-colors">
              Upgrade to Pro
            </Link>
          </div>
        </div>

        {/* FAQ Section */}
        <div className="max-w-2xl mx-auto">
          <h3 className="text-2xl font-bold mb-8 uppercase tracking-tighter">FAQ</h3>
          <div className="border-t-2 border-foreground">
            {faqs.map((faq, i) => (
              <div key={i} className="border-b border-border">
                <button 
                  onClick={() => setOpenFaq(openFaq === i ? null : i)}
                  className="w-full py-6 flex justify-between items-center text-left hover:text-secondary transition-colors"
                >
                  <span className="font-bold text-lg">{faq.q}</span>
                  <span className="font-bold text-xl">{openFaq === i ? "-" : "+"}</span>
                </button>
                {openFaq === i && (
                  <div className="pb-6 text-secondary leading-relaxed">
                    {faq.a}
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>
      </section>
    </main>
  );
}
