// File: services/user.ts

import axios from "axios";

const API_URL = process.env.NEXT_PUBLIC_API_URL;

const getAuthHeaders = () => ({
  headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
});

export const UserService = {
  async getAuthenticatedUser() {
    try {
      const token = localStorage.getItem("token");
      if (!token) {
        console.warn("Nenhum token encontrado, redirecionando para login.");
        return null;
      }

      const response = await axios.get(`${API_URL}/auth/me`, getAuthHeaders());
      return response.data;
    } catch (error: any) {
      if (error.response?.status === 401) {
        console.warn("Token inválido ou expirado. Redirecionando para login...");
        localStorage.removeItem("token");
      }
      console.error("Erro ao buscar usuário autenticado", error);
      return null;
    }
  },
};
