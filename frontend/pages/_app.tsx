// File: pages/_app.tsx

import "@/styles/globals.css";
import type { AppProps } from "next/app";
import { useEffect, useState } from "react";
import { UserProvider } from "@/context/UserContext";
import { Toaster } from "@/components/ui/sonner";
import Footer from "@/components/Footer";

export default function App({ Component, pageProps }: AppProps) {
  const getLayout = (Component as any).getLayout || ((page: JSX.Element) => page);

  // 🌙 Gerencia o Tema
  const [theme, setTheme] = useState<"light" | "dark">("light");

  useEffect(() => {
    if (typeof window !== "undefined") {
      const storedTheme = localStorage.getItem("theme") as "light" | "dark" | null;
      const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
      const initialTheme = storedTheme || (prefersDark ? "dark" : "light");

      setTheme(initialTheme);
      document.documentElement.classList.toggle("dark", initialTheme === "dark");

      // 🔹 Carregar script do reCAPTCHA v3
      const script = document.createElement("script");
      script.src = `https://www.google.com/recaptcha/api.js?render=${process.env.NEXT_PUBLIC_RECAPTCHA_SITE_KEY}`;
      script.async = true;
      script.defer = true;
      document.body.appendChild(script);

      script.onload = () => {
        console.log("✅ reCAPTCHA carregado com sucesso!");
      };

      return () => {
        document.body.removeChild(script);
      };
    }
  }, []);

  return (
    <UserProvider>
      {getLayout(<Component {...pageProps} />)}
      <Toaster />
      <Footer />
    </UserProvider>
  );
}
