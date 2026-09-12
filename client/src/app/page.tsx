"use client";

import Link from "next/link";

function Navbar() {
  return (
    <nav className="flex items-center justify-between px-6 py-4 md:px-12 border-b border-border bg-background/80 backdrop-blur-md sticky top-0 z-50">
      <Link href="/" className="font-bold tracking-tight text-xl">
        INKTYPE
      </Link>
      <div className="hidden md:flex items-center gap-8 text-sm font-medium">
        <a href="#about" className="text-secondary hover:text-foreground transition-colors">About</a>
        <Link href="/pricing" className="text-secondary hover:text-foreground transition-colors">Pricing</Link>
        <Link href="/signin" className="text-secondary hover:text-foreground transition-colors">Sign In</Link>
        <Link 
          href="/signin" 
          className="bg-foreground text-background px-4 py-2 hover:opacity-90 transition-opacity"
        >
          Get Started -&gt;
        </Link>
      </div>
    </nav>
  );
}

function Hero() {
  return (
    <section className="px-6 py-24 md:py-32 md:px-12 flex flex-col items-center text-center">
      <h1 className="text-6xl md:text-8xl lg:text-[10rem] font-bold tracking-tighter leading-none mb-8">
        INKTYPE
      </h1>
      <p className="text-xl md:text-2xl text-secondary max-w-2xl mb-12 leading-relaxed">
        Turn your messy handwritten notes into beautiful, clean, professionally typeset PDFs instantly.
      </p>
      <Link 
        href="/signin" 
        className="bg-foreground text-background px-8 py-4 text-lg font-bold hover:opacity-90 transition-opacity flex items-center gap-4"
      >
        Try InkType for free
        <span>-&gt;</span>
      </Link>
    </section>
  );
}

function About() {
  const features = [
    {
      title: "Reads any handwriting",
      description: "Our AI understands cursive, block letters, math equations, and messy margins."
    },
    {
      title: "Smart typesetting",
      description: "We don't just transcribe. We format headings, lists, and paragraphs with editorial precision."
    },
    {
      title: "Export to PDF",
      description: "Download a beautiful, searchable PDF ready for publishing, sharing, or studying."
    }
  ];

  return (
    <section id="about" className="px-6 py-24 md:px-12 bg-sidebar border-y border-border">
      <div className="max-w-6xl mx-auto">
        <div className="mb-20">
          <h2 className="text-3xl md:text-5xl font-bold tracking-tight mb-6">How it works</h2>
          <p className="text-xl text-secondary max-w-2xl">
            A minimal tool designed for maximum clarity. We do the heavy lifting so you can focus on writing.
          </p>
        </div>
        
        <div className="grid md:grid-cols-3 gap-16">
          {features.map((feature, i) => (
            <div key={i} className="flex flex-col border-l-2 border-foreground pl-6">
              <h3 className="text-xl font-bold mb-4">{(i+1).toString().padStart(2, '0')}. {feature.title}</h3>
              <p className="text-secondary leading-relaxed">{feature.description}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

function Footer() {
  return (
    <footer className="px-6 py-12 md:px-12 border-t border-border flex flex-col md:flex-row items-center justify-between gap-6 text-sm text-secondary font-medium">
      <div>
        © {new Date().getFullYear()} InkType. All rights reserved.
      </div>
      <div className="flex gap-8">
        <a href="#" className="hover:text-foreground transition-colors">Twitter</a>
        <a href="#" className="hover:text-foreground transition-colors">GitHub</a>
        <a href="#" className="hover:text-foreground transition-colors">Terms</a>
        <a href="#" className="hover:text-foreground transition-colors">Privacy</a>
      </div>
    </footer>
  );
}

export default function LandingPage() {
  return (
    <main className="min-h-full w-full flex flex-col bg-background text-foreground overflow-y-auto overflow-x-hidden">
      <Navbar />
      <Hero />
      <About />
      
      <div className="flex-1 min-h-[15vh]"></div>
      
      <Footer />
    </main>
  );
}
