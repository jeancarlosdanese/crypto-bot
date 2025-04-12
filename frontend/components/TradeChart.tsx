// components/TradeChart.tsx

"use client";

import { useEffect, useRef } from "react";
import {
  createChart,
  ColorType,
  CrosshairMode,
  CandlestickData,
  Time,
  IChartApi,
  ISeriesApi,
} from "lightweight-charts";

type Decision = {
  time: Time;
  price: number;
  decision: "BUY" | "SELL";
};

type Props = {
  botID: string;
  token: string;
};

export default function TradeChart({ botID, token }: Props) {
  const chartRef = useRef<HTMLDivElement | null>(null);
  const chartApiRef = useRef<IChartApi | null>(null);
  const candleSeriesRef = useRef<ISeriesApi<"Candlestick"> | null>(null);
  const markersRef = useRef<Decision[]>([]);

  useEffect(() => {
    if (!chartRef.current) return;

    // Criar o gráfico
    const chart = createChart(chartRef.current, {
      width: chartRef.current.clientWidth,
      height: 400,
      layout: {
        background: { type: ColorType.Solid, color: "#ffffff" },
        textColor: "#333",
      },
      grid: {
        vertLines: { color: "#eee" },
        horzLines: { color: "#eee" },
      },
      crosshair: {
        mode: CrosshairMode.Normal,
      },
      timeScale: {
        borderColor: "#ccc",
      },
    });

    const candleSeries = chart.addCandlestickSeries();
    chartApiRef.current = chart;
    candleSeriesRef.current = candleSeries;

    // Conectar ao WebSocket
    const socket = new WebSocket(`ws://localhost:8080/ws/${botID}?token=${token}`);

    socket.onopen = () => {
      console.log("✅ WebSocket connection established");
    };

    socket.onclose = () => {
      console.log("❌ WebSocket connection closed");
    };

    socket.onerror = (error) => {
      console.error("⚠️ WebSocket error:", error);
    };

    socket.onmessage = (event) => {
      const msg = JSON.parse(event.data);
      console.log("📨 WebSocket message received:", msg);

      if (msg.type === "candle") {
        const candle = msg.data as CandlestickData;
        candleSeries.update(candle);
      }

      if (msg.type === "decision") {
        const d = msg.data as Decision;
        markersRef.current.push(d);

        candleSeries.setMarkers(
          markersRef.current.map((m) => ({
            time: m.time,
            position: m.decision === "BUY" ? "belowBar" : "aboveBar",
            color: m.decision === "BUY" ? "green" : "red",
            shape: "arrowUp",
            text: m.decision,
          }))
        );
      }
    };

    // Cleanup
    return () => {
      socket.close();
      chart.remove();
    };
  }, [botID, token]);

  return <div ref={chartRef} className="w-full h-[400px]" />;
}
