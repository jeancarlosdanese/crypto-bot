// components/TradeChartEMAFan.tsx

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

export default function TradeChartEMAFan({ botID, token }: Props) {
  const chartRef = useRef<HTMLDivElement | null>(null);
  const chartApiRef = useRef<IChartApi | null>(null);
  const candleSeriesRef = useRef<ISeriesApi<"Candlestick"> | null>(null);
  const markersRef = useRef<Decision[]>([]);
  const emaRefs = useRef<Record<number, ISeriesApi<"Line">>>({});

  useEffect(() => {
    if (!chartRef.current) return;

    const isDark = document.documentElement.classList.contains("dark");

    const chart = createChart(chartRef.current, {
      width: chartRef.current.clientWidth,
      layout: {
        background: { type: ColorType.Solid, color: isDark ? "#121212" : "#ffffff" },
        textColor: isDark ? "#f5f5f5" : "#333",
      },
      grid: {
        vertLines: { color: isDark ? "#444" : "#eee" },
        horzLines: { color: isDark ? "#444" : "#eee" },
      },
      crosshair: {
        mode: CrosshairMode.Normal,
      },
      timeScale: {
        borderColor: isDark ? "#666" : "#ccc",
        timeVisible: true,
        secondsVisible: true,
        rightOffset: 20,
      },
    });

    const resizeObserver = new ResizeObserver(() => {
      if (chartRef.current) {
        chart.resize(chartRef.current.clientWidth, chartRef.current.clientHeight);
      }
    });

    const candleSeries = chart.addCandlestickSeries();
    chartApiRef.current = chart;
    candleSeriesRef.current = candleSeries;

    const emaPeriods = [10, 15, 20, 25, 30, 35, 40];
    const colors = ["#fbbf24", "#22c55e", "#3b82f6", "#8b5cf6", "#ec4899", "#ef4444", "#14b8a6"];

    emaPeriods.forEach((period, i) => {
      const series = chart.addLineSeries({
        color: colors[i],
        lineWidth: 2,
      });
      emaRefs.current[period] = series;
    });

    const loadHistorical = async () => {
      try {
        type EMACandle = CandlestickData & {
          [key: string]: number | Time;
        };

        const historical: EMACandle[] = await BotService.getHistotical(botID);
        if (!historical) {
          console.error("Erro ao carregar dados históricos");
          return;
        }

        candleSeries.setData(historical);

        emaPeriods.forEach((period) => {
          const key = `ema${period}`;
          const data = historical
            .filter((c) => c[key] !== undefined)
            .map((c) => ({
              time: c.time,
              value: c[key] as number,
            }));

          emaRefs.current[period]?.setData(data);
        });
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
        chart.timeScale().scrollToRealTime();

        emaPeriods.forEach((period) => {
          const key = `ema${period}`;
          if (candle[key]) {
            emaRefs.current[period]?.update({ time: candle.time, value: candle[key] });
          }
        });
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

    resizeObserver.observe(chartRef.current);

    return () => {
      socket.close();
      chart.remove();
      resizeObserver.disconnect();
    };
  }, [botID, token]);

  return <div ref={chartRef} className="relative w-full min-h-[500px] h-[500px]" />;
}
