"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";

export default function CheckoutPage() {
  const [loading, setLoading] = useState(true);
  const router = useRouter();

  useEffect(() => {
    const timer = setTimeout(() => {
      setLoading(false);
      setTimeout(() => {
        router.push("/app");
      }, 2000);
    }, 1500);
    
    return () => clearTimeout(timer);
  }, [router]);

  return (
    <main className="min-h-full w-full flex flex-col items-center justify-center bg-background text-foreground p-6">
      <div className="w-full max-w-md text-center animate-in fade-in duration-300">
        {loading ? (
          <div>
            <h1 className="text-2xl font-bold mb-6 uppercase tracking-widest">Processing Payment</h1>
            <div className="w-full h-1 bg-border relative overflow-hidden">
              <div className="absolute top-0 left-0 h-full w-1/2 bg-foreground animate-[pulse_1s_ease-in-out_infinite]"></div>
            </div>
          </div>
        ) : (
          <div>
            <div className="text-6xl font-bold mb-6">+</div>
            <h1 className="text-3xl font-bold mb-4 uppercase tracking-tighter">Welcome to Pro</h1>
            <p className="text-secondary font-bold tracking-widest uppercase text-sm mb-12">
              Redirecting to workspace...
            </p>
          </div>
        )}
      </div>
    </main>
  );
}
