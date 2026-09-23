import type { Metadata } from "next";
import { IBM_Plex_Mono } from "next/font/google";
import "./globals.css";

const ibmPlexMono = IBM_Plex_Mono({
  variable: "--font-ibm-plex-mono",
  subsets: ["latin"],
  weight: ["400", "500", "600", "700"],
});

export const metadata: Metadata = {
  title: "InkType",
  description: "Turn handwriting into clean, beautifully typeset PDFs.",
};

import { ClerkProvider } from "@clerk/nextjs";

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <ClerkProvider>
      <html lang="en" className="h-full">
        <body
          className={`${ibmPlexMono.variable} ${ibmPlexMono.className} antialiased h-full text-foreground bg-background`}
        >
          {children}
        </body>
      </html>
    </ClerkProvider>
  );
}
