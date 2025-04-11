// File: pages/api/auth/request-otp.ts

import type { NextApiRequest, NextApiResponse } from "next";
import axios from "axios";

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  if (req.method !== "POST") {
    return res.status(405).json({ message: "Método não permitido" });
  }

  try {
    const { identifier } = req.body;

    if (!identifier) {
      return res.status(400).json({ message: "E-mail ou WhatsApp é obrigatório" });
    }

    // Chama a API do backend `go-marketing` para gerar e enviar OTP
    await axios.post(`${process.env.NEXT_PUBLIC_API_URL}/auth/request-otp`, { identifier });

    return res.status(200).json({ message: "OTP enviado com sucesso" });
  } catch (error) {
    console.error("Erro ao solicitar OTP:", error);
    return res.status(500).json({ message: "Erro interno ao solicitar OTP" });
  }
}
