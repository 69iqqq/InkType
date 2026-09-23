import { SignIn } from "@clerk/nextjs";
import Image from "next/image";
import Link from "next/link";

export default function SignInPage() {
  return (
    <div className="min-h-screen flex flex-col justify-between bg-background text-foreground font-mono">
      {/* Top Navigation */}
      <header className="px-6 py-6 md:px-12 flex items-center justify-between border-b border-border">
        <Link href="/" className="flex items-center gap-3 hover:opacity-80 transition-opacity">
          <Image 
            src="/InkTypeLogo-v2.png" 
            alt="InkType" 
            width={120} 
            height={40} 
            className="h-7 w-auto object-contain" 
            priority 
          />
        </Link>
        <Link 
          href="/" 
          className="text-xs text-secondary hover:text-foreground transition-colors uppercase tracking-widest font-medium"
        >
          [ ← Home ]
        </Link>
      </header>

      {/* Main Content */}
      <main className="flex-1 flex flex-col items-center justify-center px-4 py-12">
        <div className="w-full max-w-md flex flex-col items-center">
          {/* Header Info */}
          <div className="text-center mb-8">
            <span className="text-[11px] font-bold text-secondary uppercase tracking-widest bg-sidebar px-3 py-1 rounded-sm border border-border inline-block mb-3">
              Authentication
            </span>
            <h1 className="text-2xl font-bold tracking-tight text-foreground">
              Sign in to InkType
            </h1>
            <p className="text-xs text-secondary mt-1">
              Turn handwriting into something clean.
            </p>
          </div>

          {/* Styled Clerk SignIn Card */}
          <div className="w-full">
            <SignIn 
              routing="hash" 
              fallbackRedirectUrl="/app" 
              signUpFallbackRedirectUrl="/welcome"
              appearance={{
                variables: {
                  colorPrimary: "#111111",
                  colorBackground: "#FAFAFA",
                  borderRadius: "0px",
                },
                elements: {
                  rootBox: "w-full",
                  card: "shadow-none border-2 border-border bg-background p-6 md:p-8 rounded-none w-full",
                  header: "hidden",
                  socialButtonsBlockButton: "border-2 border-border bg-background hover:bg-hover hover:border-foreground text-foreground text-xs font-mono font-medium py-3 rounded-none transition-colors shadow-none",
                  socialButtonsBlockButtonText: "font-mono text-xs font-semibold",
                  dividerRow: "my-6",
                  dividerLine: "bg-border h-[1px]",
                  dividerText: "text-[10px] font-mono text-secondary uppercase tracking-widest px-3 bg-background",
                  formFieldLabel: "text-xs font-mono text-secondary uppercase tracking-wider mb-1.5",
                  formFieldInput: "bg-background border-2 border-border rounded-none text-foreground font-mono text-sm py-2.5 px-3 focus:border-foreground focus:outline-none transition-colors shadow-none",
                  formButtonPrimary: "bg-foreground text-background hover:bg-neutral-800 font-mono text-xs uppercase tracking-widest font-bold py-3.5 rounded-none shadow-none transition-colors border-2 border-foreground mt-3",
                  footer: "pt-6 border-t border-border mt-6 bg-transparent",
                  footerAction: "font-mono text-xs text-secondary",
                  footerActionLink: "font-mono text-xs font-bold text-foreground underline hover:opacity-70 transition-opacity ml-1",
                  identityPreviewText: "font-mono text-xs text-foreground",
                  identityPreviewEditButton: "font-mono text-xs text-secondary hover:text-foreground",
                  formFieldAction: "font-mono text-xs text-secondary hover:text-foreground underline",
                }
              }}
            />
          </div>
        </div>
      </main>

      {/* Footer */}
      <footer className="px-6 py-6 text-center border-t border-border text-xs text-secondary uppercase tracking-widest">
        InkType &copy; {new Date().getFullYear()} &bull; Minimal Handwritten Notes
      </footer>
    </div>
  );
}
