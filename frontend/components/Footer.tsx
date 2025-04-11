// File: components/Footer.tsx

import Link from "next/link";

export default function Footer() {
  return (
    <footer className="w-full bg-background border-t border-border text-muted-foreground text-sm text-center py-3">
      <p>
        Next-Marketing v1.0.0 - Desenvolvido por{" "}
        <Link
          href="https://github.com/jeancarlosdanese/"
          target="_blank"
          className="underline hover:text-primary transition"
        >
          Jean Carlos Danese
        </Link>
      </p>
      <p className="opacity-75">© {new Date().getFullYear()} Todos os direitos reservados.</p>

      {/* 🔹 Mensagem do reCAPTCHA */}
      <p className="opacity-50 text-xs mt-2">
        Este site é protegido pelo reCAPTCHA e está sujeito à{" "}
        <a
          href="https://policies.google.com/privacy"
          className="underline"
          target="_blank"
          rel="noopener noreferrer"
        >
          Política de Privacidade
        </a>{" "}
        e{" "}
        <a
          href="https://policies.google.com/terms"
          className="underline"
          target="_blank"
          rel="noopener noreferrer"
        >
          Termos de Serviço
        </a>{" "}
        do Google.
      </p>
    </footer>
  );
}
