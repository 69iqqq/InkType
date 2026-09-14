"use client";

import { useEffect } from "react";
import Link from "next/link";
import Image from "next/image";
import { motion, useMotionValue, useSpring } from "framer-motion";

function Navbar() {
  return (
    <nav className="flex items-center justify-between px-6 py-4 md:px-12 border-b border-border bg-background/80 backdrop-blur-md fixed w-full top-0 z-50">
      <Link href="/">
        <Image src="/InkTypeLogo-v2.png" alt="InkType" width={120} height={40} className="h-8 w-auto object-contain" priority />
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
    <section className="relative px-6 py-24 md:py-40 md:px-12 flex flex-col items-center justify-center text-center min-h-screen overflow-hidden bg-[#F3F3F2]">
      
      {/* --- SEPARATED ARTWORK LAYERS --- */}
      {/* These pieces assemble from the edges and float independently to create a dynamic, living composition on all screens. */}
      
      {/* Left Artwork (Statue) */}
      <motion.div 
        className="absolute left-0 bottom-0 h-[40%] md:h-[60%] lg:h-[90%] w-[45vw] md:w-[35vw] lg:w-[30vw] z-0"
        initial={{ x: "-10vw", opacity: 0 }}
        animate={{ x: 0, opacity: 1 }}
        transition={{ type: "spring", bounce: 0.2, duration: 1.5, delay: 0.1 }}
      >
        <motion.div 
          className="relative w-full h-full"
          animate={{ y: [0, -15, 0] }}
          transition={{ duration: 6, repeat: Infinity, ease: "easeInOut" }}
        >
          <Image
            src="/hero-left.webp"
            alt="Decorative statue"
            fill
            priority
            className="object-contain object-left-bottom lg:object-left"
            sizes="(max-width: 1024px) 50vw, 30vw"
          />
        </motion.div>
      </motion.div>

      {/* Right Top Artwork (Eye/Doodles) */}
      <motion.div 
        className="absolute right-0 top-0 h-[25%] sm:h-[30%] md:h-[40%] lg:h-[55%] w-[45vw] sm:w-[50vw] md:w-[30vw] lg:w-[35vw] z-0"
        initial={{ x: "10vw", opacity: 0 }}
        animate={{ x: 0, opacity: 1 }}
        transition={{ type: "spring", bounce: 0.2, duration: 1.5, delay: 0.2 }}
      >
        <motion.div 
          className="relative w-full h-full"
          animate={{ y: [0, 10, 0] }}
          transition={{ duration: 7, repeat: Infinity, ease: "easeInOut", delay: 1 }}
        >
          <Image
            src="/hero-right-top.webp"
            alt="Decorative eye"
            fill
            priority
            className="object-contain object-right-top"
            sizes="(max-width: 1024px) 50vw, 35vw"
          />
        </motion.div>
      </motion.div>

      {/* Right Bottom Artwork (Grid/Mountain) */}
      <motion.div 
        className="absolute right-0 bottom-0 h-[35%] md:h-[50%] lg:h-[60%] w-[55vw] md:w-[40vw] lg:w-[40vw] z-0"
        initial={{ x: "10vw", opacity: 0 }}
        animate={{ x: 0, opacity: 1 }}
        transition={{ type: "spring", bounce: 0.2, duration: 1.5, delay: 0.3 }}
      >
        <motion.div 
          className="relative w-full h-full"
          animate={{ y: [0, -12, 0] }}
          transition={{ duration: 8, repeat: Infinity, ease: "easeInOut", delay: 0.5 }}
        >
          <Image
            src="/hero-right-bottom.webp"
            alt="Decorative sketches"
            fill
            priority
            className="object-contain object-right-bottom"
            sizes="(max-width: 1024px) 50vw, 40vw"
          />
        </motion.div>
      </motion.div>

      {/* --- Foreground Content --- */}
      <motion.div 
        className="relative z-10 flex flex-col items-center mt-12 w-full max-w-4xl mx-auto"
        initial={{ opacity: 0, scale: 0.98 }}
        animate={{ opacity: 1, scale: 1 }}
        transition={{ duration: 1, ease: "easeOut", delay: 0.4 }}
      >
        <motion.h1 
          className="text-6xl md:text-8xl lg:text-[10rem] font-bold tracking-tighter leading-none mb-8 text-foreground"
          initial={{ opacity: 0, scale: 0.95 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ duration: 0.8, delay: 0.6, ease: "easeOut" }}
        >
          INKTYPE
        </motion.h1>
        <motion.p 
          className="text-xl md:text-2xl text-foreground max-w-2xl mb-12 leading-relaxed font-medium"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.8, delay: 0.8, ease: "easeOut" }}
        >
          Turn your messy handwritten notes into beautiful, clean, professionally typeset PDFs instantly.
        </motion.p>
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.8, delay: 1.0, ease: "easeOut" }}
        >
          <Link 
            href="/signin" 
            className="bg-foreground text-background px-8 py-4 text-lg font-bold hover:scale-105 hover:shadow-xl transition-all duration-300 flex items-center gap-4 group"
          >
            Try InkType for free
            <span className="group-hover:translate-x-1 transition-transform">-&gt;</span>
          </Link>
        </motion.div>
      </motion.div>
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

function CustomCursor() {
  const cursorX = useMotionValue(-100);
  const cursorY = useMotionValue(-100);
  
  const springConfig = { damping: 25, stiffness: 400 };
  const cursorXSpring = useSpring(cursorX, springConfig);
  const cursorYSpring = useSpring(cursorY, springConfig);

  useEffect(() => {
    const moveCursor = (e: MouseEvent) => {
      cursorX.set(e.clientX - 12);
      cursorY.set(e.clientY - 12);
    };
    window.addEventListener("mousemove", moveCursor);
    return () => {
      window.removeEventListener("mousemove", moveCursor);
    };
  }, [cursorX, cursorY]);

  return (
    <motion.div
      className="hidden md:block fixed top-0 left-0 w-6 h-6 border-2 border-foreground bg-foreground/10 backdrop-invert rounded-full pointer-events-none z-[9999]"
      style={{
        translateX: cursorXSpring,
        translateY: cursorYSpring,
      }}
    />
  );
}

export default function LandingPage() {
  return (
    <main className="min-h-full w-full flex flex-col bg-background text-foreground overflow-x-hidden md:cursor-none md:[&_a]:cursor-none md:[&_button]:cursor-none">
      <CustomCursor />
      <Navbar />
      <Hero />
      <About />
      
      <div className="flex-1 min-h-[15vh]"></div>
      
      <Footer />
    </main>
  );
}
