"use client";

import React, { useState } from "react";
import Link from "next/link";
import Image from "next/image";
import { useRouter } from "next/navigation";

export default function WelcomePage() {
  const [step, setStep] = useState(1);
  const router = useRouter();

  const handleNext = () => {
    if (step === 1) {
      setStep(2);
    } else {
      router.push("/app");
    }
  };

  return (
    <main className="min-h-full w-full flex flex-col bg-background text-foreground">
      <nav className="p-6 md:p-12">
        <Image src="/InkTypeLogo-v2.png" alt="InkType" width={120} height={40} className="h-8 w-auto object-contain" priority />
      </nav>

      <div className="flex-1 flex items-center justify-center p-6">
        <div className="w-full max-w-2xl">
          <div className="text-xs font-bold text-secondary tracking-widest mb-12">
            STEP 0{step} / 02
          </div>

          {step === 1 && (
            <div className="animate-in fade-in duration-300">
              <h1 className="text-4xl md:text-5xl font-bold mb-12 tracking-tighter">
                How do you plan to use InkType?
              </h1>
              
              <div className="grid md:grid-cols-3 gap-6 mb-12">
                {["Academic", "Professional", "Personal"].map((intent) => (
                  <button 
                    key={intent}
                    onClick={handleNext}
                    className="border-2 border-border hover:border-foreground p-8 text-left transition-colors flex flex-col justify-between h-40"
                  >
                    <span className="font-bold text-xl">{intent}</span>
                    <span className="text-secondary">-&gt;</span>
                  </button>
                ))}
              </div>
            </div>
          )}

          {step === 2 && (
            <div className="animate-in fade-in duration-300">
              <h1 className="text-4xl md:text-5xl font-bold mb-12 tracking-tighter">
                Select your default export font.
              </h1>
              
              <div className="grid md:grid-cols-2 gap-6 mb-12">
                <button 
                  onClick={handleNext}
                  className="border-2 border-border hover:border-foreground p-8 text-left transition-colors group"
                >
                  <div className="font-serif text-3xl mb-4 group-hover:underline">Serif</div>
                  <div className="text-sm text-secondary font-sans">Best for academic papers and essays.</div>
                </button>
                <button 
                  onClick={handleNext}
                  className="border-2 border-border hover:border-foreground p-8 text-left transition-colors group"
                >
                  <div className="font-sans text-3xl mb-4 group-hover:underline">Sans-Serif</div>
                  <div className="text-sm text-secondary font-sans">Best for technical notes and modern reports.</div>
                </button>
              </div>
            </div>
          )}

        </div>
      </div>
    </main>
  );
}
