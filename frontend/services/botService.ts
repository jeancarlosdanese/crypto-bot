// services/botService.ts

import axios from "axios"
import { Bot } from "@/types/bot"

const API_URL = process.env.NEXT_PUBLIC_API_URL

const getAuthHeaders = () => ({
  headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
})

export const BotService = {
  async getAll(): Promise<Bot[] | null> {
    try {
      const response = await axios.get(`${API_URL}/bots`, getAuthHeaders())
      return response.data
    } catch (error: any) {
      console.error("Erro ao buscar bots:", error)
      return null
    }
  },
  async getHistotical(botID: string): Promise<any> {
    try {
      const response = await axios.get(`${API_URL}/bots/${botID}/candles`, getAuthHeaders())
      return response.data
    } catch (error: any) {
      console.error("Erro ao buscar candles:", error)
      return null
    }
  },
}

