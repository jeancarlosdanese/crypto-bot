// components/TradeChartRSI.tsx

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
import { BotService } from "@/services/botService";

type Decision = {
  time: Time;
  price: number;
  decision: "BUY" | "SELL";
};

type Props = {
  botID: string;
  token: string;
};

export default function TradeChartRSI({ botID, token }: Props) {
  const chartRef = useRef<HTMLDivElement | null>(null);
  const candleSeriesRef = useRef<ISeriesApi<"Candlestick"> | null>(null);
  const rsiSeriesRef = useRef<ISeriesApi<"Line"> | null>(null);
  const markersRef = useRef<Decision[]>([]);

  useEffect(() => {
    if (!chartRef.current) return;

    const chart = createChart(chartRef.current, {
      width: chartRef.current.clientWidth,
      height: 500,
      layout: {
        background: { type: ColorType.Solid, color: "#fff" },
        textColor: "#000",
      },
      crosshair: {
        mode: CrosshairMode.Normal,
      },
      timeScale: {
        timeVisible: true,
        secondsVisible: true,
        borderColor: "#ccc",
        rightOffset: 20,
      },
    });

    // Aplica margens visuais para RSI
    chart.priceScale("right").applyOptions({
      scaleMargins: { top: 0.2, bottom: 0.2 },
    });

    const candleSeries = chart.addCandlestickSeries();
    const rsiSeries = chart.addLineSeries({
      color: "#3b82f6",
      lineWidth: 2,
      priceScaleId: "right", // mesmo painel do candle
    });

    candleSeriesRef.current = candleSeries;
    rsiSeriesRef.current = rsiSeries;

    const loadHistorical = async () => {
      try {
        type RSICandle = CandlestickData & { rsi?: number };

        const historical: RSICandle[] = await BotService.getHistotical(botID);
        if (!historical) return;

        candleSeries.setData(historical);
        rsiSeries.setData(
          historical
            .filter((c) => c.rsi !== undefined)
            .map((c) => ({ time: c.time, value: c.rsi! }))
        );
      } catch (error) {
        console.error("Erro ao carregar dados históricos", error);
      }
    };

    loadHistorical();

    const socket = new WebSocket(`ws://localhost:8080/ws/${botID}?token=${token}`);

    socket.onmessage = (event) => {
      const msg = JSON.parse(event.data);
      if (msg.type === "candle") {
        const candle = msg.data;
        candleSeries.update(candle);
        if (candle.rsi) {
          rsiSeries.update({ time: candle.time, value: candle.rsi });
        }
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

    return () => {
      socket.close();
      chart.remove();
    };
  }, [botID, token]);

  return <div ref={chartRef} className="relative w-full min-h-[500px] h-[500px]" />;
}
