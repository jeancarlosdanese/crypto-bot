// File: pages/api/auth/verify-otp.ts

import type { NextApiRequest, NextApiResponse } from "next";
import axios from "axios";

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  if (req.method !== "POST") {
    return res.status(405).json({ message: "Método não permitido" });
  }

  try {
    const { identifier, otp } = req.body;

    if (!identifier || !otp) {
      return res.status(400).json({ message: "E-mail/WhatsApp e OTP são obrigatórios" });
    }

    // Chama a API do backend `go-marketing` para verificar OTP
    const response = await axios.post(`${process.env.NEXT_PUBLIC_API_URL}/auth/verify-otp`, {
      identifier,
      otp,
    });

    // Salva o token JWT na resposta
    return res.status(200).json(response.data);
  } catch (error: unknown) {
    if (axios.isAxiosError(error)) {
      console.error("Erro ao verificar OTP:", error.response?.data || error.message);
      return res.status(401).json({ message: "OTP inválido ou expirado" });
    } else {
      console.error("Erro ao verificar OTP:", error);
      return res.status(500).json({ message: "Erro interno do servidor" });
    }
  }
}
