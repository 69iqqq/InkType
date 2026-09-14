"use client";

import React, { useState } from "react";
import Link from "next/link";
import Image from "next/image";
import { useRouter } from "next/navigation";

export default function SignInPage() {
  const [email, setEmail] = useState("");
  const [submitted, setSubmitted] = useState(false);
  const router = useRouter();

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (email) {
      setSubmitted(true);
      // Mock redirect after 1.5s
      setTimeout(() => {
        router.push("/welcome");
      }, 1500);
    }
  };

  return (
    <main className="min-h-full w-full flex flex-col bg-background text-foreground">
      <nav className="p-6 md:p-12">
        <Link href="/">
          <Image src="/InkTypeLogo-v2.png" alt="InkType" width={120} height={40} className="h-8 w-auto object-contain" priority />
        </Link>
      </nav>

      <div className="flex-1 flex items-center justify-center p-6">
        <div className="w-full max-w-md">
          {submitted ? (
            <div className="animate-in fade-in duration-300">
              <h1 className="text-3xl font-bold mb-4">Check your email</h1>
              <p className="text-secondary mb-8">
                We sent a magic link to <span className="text-foreground font-bold">{email}</span>.
              </p>
              <div className="w-full h-1 bg-foreground animate-pulse"></div>
            </div>
          ) : (
            <form onSubmit={handleSubmit} className="animate-in fade-in duration-300">
              <h1 className="text-3xl font-bold mb-12">Sign In</h1>
              
              <div className="mb-12">
                <label htmlFor="email" className="block text-xs font-bold text-secondary uppercase tracking-widest mb-4">
                  Email Address
                </label>
                <input
                  type="email"
                  id="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="name@example.com"
                  className="w-full bg-transparent border-b-2 border-border focus:border-foreground py-2 outline-none transition-colors text-xl"
                  required
                />
              </div>

              <div className="flex items-center justify-between">
                <Link href="/" className="text-sm font-bold text-secondary hover:text-foreground transition-colors">
                  &lt;- Back
                </Link>
                <button 
                  type="submit"
                  className="bg-foreground text-background px-8 py-3 font-bold hover:bg-transparent hover:text-foreground border-2 border-foreground transition-colors"
                >
                  Send Link
                </button>
              </div>
            </form>
          )}
        </div>
      </div>
    </main>
  );
}
