// components/TradeChartMACD.tsx

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

export default function TradeChartMACD({ botID, token }: Props) {
  const chartRef = useRef<HTMLDivElement | null>(null);
  const candleSeriesRef = useRef<ISeriesApi<"Candlestick"> | null>(null);
  const macdLineRef = useRef<ISeriesApi<"Line"> | null>(null);
  const signalLineRef = useRef<ISeriesApi<"Line"> | null>(null);
  const histogramRef = useRef<ISeriesApi<"Histogram"> | null>(null);
  const markersRef = useRef<Decision[]>([]);

  useEffect(() => {
    if (!chartRef.current) return;

    const chart = createChart(chartRef.current, {
      width: chartRef.current.clientWidth,
      height: 500,
      layout: {
        background: { type: ColorType.Solid, color: "#ffffff" },
        textColor: "#000000",
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

    chart.priceScale("right").applyOptions({
      scaleMargins: { top: 0.2, bottom: 0.4 },
    });

    const candleSeries = chart.addCandlestickSeries();
    const macdLine = chart.addLineSeries({
      color: "#1e40af",
      lineWidth: 2,
      priceScaleId: "right",
    });
    const signalLine = chart.addLineSeries({
      color: "#f59e0b",
      lineWidth: 2,
      priceScaleId: "right",
    });

    const histogram = chart.addHistogramSeries({
      priceScaleId: "right",
      color: "#8884d8",
      base: 0,
    });

    candleSeriesRef.current = candleSeries;
    macdLineRef.current = macdLine;
    signalLineRef.current = signalLine;
    histogramRef.current = histogram;

    const loadHistorical = async () => {
      try {
        type MACDCandle = CandlestickData & {
          macd?: number;
          macd_signal?: number;
          macd_histogram?: number;
        };

        const historical: MACDCandle[] = await BotService.getHistotical(botID);
        if (!historical) return;

        candleSeries.setData(historical);

        macdLine.setData(
          historical
            .filter((c) => c.macd !== undefined)
            .map((c) => ({ time: c.time, value: c.macd! }))
        );

        signalLine.setData(
          historical
            .filter((c) => c.macd_signal !== undefined)
            .map((c) => ({ time: c.time, value: c.macd_signal! }))
        );

        histogram.setData(
          historical
            .filter((c) => c.macd_histogram !== undefined)
            .map((c) => ({
              time: c.time,
              value: c.macd_histogram!,
              color: c.macd_histogram! >= 0 ? "#22c55e" : "#ef4444",
            }))
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

        if (candle.macd !== undefined) macdLine.update({ time: candle.time, value: candle.macd });
        if (candle.macd_signal !== undefined)
          signalLine.update({ time: candle.time, value: candle.macd_signal });
        if (candle.macd_histogram !== undefined)
          histogram.update({
            time: candle.time,
            value: candle.macd_histogram,
            color: candle.macd_histogram >= 0 ? "#22c55e" : "#ef4444",
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

    return () => {
      socket.close();
      chart.remove();
    };
  }, [botID, token]);

  return <div ref={chartRef} className="relative w-full min-h-[500px] h-[500px]" />;
}
