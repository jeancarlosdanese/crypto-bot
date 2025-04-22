// components/TradeChartVolumeSpike.tsx

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

export default function TradeChartVolumeSpike({ botID, token }: Props) {
  const chartRef = useRef<HTMLDivElement | null>(null);
  const candleSeriesRef = useRef<ISeriesApi<"Candlestick"> | null>(null);
  const volumeSeriesRef = useRef<ISeriesApi<"Histogram"> | null>(null);
  const volumeAvgSeriesRef = useRef<ISeriesApi<"Line"> | null>(null);
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

    // Price chart (top)
    chart.priceScale("right").applyOptions({
      scaleMargins: { top: 0.2, bottom: 0.35 },
    });

    const candleSeries = chart.addCandlestickSeries({ priceScaleId: "right" });

    // Volume histogram (bottom overlay)
    const volumeSeries = chart.addHistogramSeries({
      priceFormat: { type: "volume" },
      color: "#c084fc",
    });

    // Average line
    const avgLineSeries = chart.addLineSeries({
      color: "#6366f1",
      lineWidth: 1,
      priceScaleId: "",
    });

    candleSeriesRef.current = candleSeries;
    volumeSeriesRef.current = volumeSeries;
    volumeAvgSeriesRef.current = avgLineSeries;

    const loadHistorical = async () => {
      try {
        type VolumeCandle = CandlestickData & {
          volume: number;
          avg_volume?: number;
        };

        const historical: VolumeCandle[] = await BotService.getHistotical(botID);
        if (!historical) return;

        candleSeries.setData(historical);

        volumeSeries.setData(
          historical.map((c) => ({
            time: c.time,
            value: c.volume,
            color: "#c084fc",
          }))
        );

        avgLineSeries.setData(
          historical
            .filter((c) => c.avg_volume !== undefined)
            .map((c) => ({
              time: c.time,
              value: c.avg_volume!,
            }))
        );
      } catch (error) {
        console.error("Erro ao carregar histórico", error);
      }
    };

    loadHistorical();

    const socket = new WebSocket(`ws://localhost:8080/ws/${botID}?token=${token}`);

    socket.onmessage = (event) => {
      const msg = JSON.parse(event.data);
      if (msg.type === "candle") {
        const candle = msg.data;
        candleSeries.update(candle);
        if (candle.volume) {
          volumeSeries.update({
            time: candle.time,
            value: candle.volume,
            color: "#c084fc",
          });
        }
        if (candle.avg_volume) {
          avgLineSeries.update({
            time: candle.time,
            value: candle.avg_volume,
          });
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
