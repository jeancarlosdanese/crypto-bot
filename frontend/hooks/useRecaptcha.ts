// File: hooks/useRecaptcha.ts

import { useEffect, useState } from "react";

const RECAPTCHA_SITE_KEY = process.env.NEXT_PUBLIC_RECAPTCHA_SITE_KEY || "";

export function useRecaptcha(action: string) {
  const [token, setToken] = useState<string | null>(null);

  useEffect(() => {
    if (!RECAPTCHA_SITE_KEY) {
      console.error("❌ reCAPTCHA site key não foi definida.");
      return;
    }

    const executeRecaptcha = () => {
      if (window.grecaptcha) {
        window.grecaptcha.ready(() => {
          window.grecaptcha
            .execute(RECAPTCHA_SITE_KEY, { action })
            .then((token: string) => {
              setToken(token);
            })
            .catch((err: any) => console.error("❌ Erro ao executar reCAPTCHA", err));
        });
      } else {
        console.error("❌ reCAPTCHA não está carregado ainda.");
      }
    };

    // Aguarda até o script ser carregado
    const interval = setInterval(() => {
      if (window.grecaptcha) {
        executeRecaptcha();
        clearInterval(interval);
      }
    }, 500);

    return () => clearInterval(interval);
  }, [action]);

  return token;
}
