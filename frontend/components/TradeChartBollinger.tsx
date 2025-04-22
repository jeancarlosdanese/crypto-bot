// components/TradeChartBollinger.tsx

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

export default function TradeChartBollinger({ botID, token }: Props) {
  const chartRef = useRef<HTMLDivElement | null>(null);
  const candleSeriesRef = useRef<ISeriesApi<"Candlestick"> | null>(null);
  const bbUpperRef = useRef<ISeriesApi<"Line"> | null>(null);
  const bbLowerRef = useRef<ISeriesApi<"Line"> | null>(null);
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

    const candleSeries = chart.addCandlestickSeries();
    const bbUpperSeries = chart.addLineSeries({
      color: "#0ea5e9", // azul
      lineWidth: 2,
      priceScaleId: "right",
    });
    const bbLowerSeries = chart.addLineSeries({
      color: "#f43f5e", // vermelho
      lineWidth: 2,
      priceScaleId: "right",
    });

    candleSeriesRef.current = candleSeries;
    bbUpperRef.current = bbUpperSeries;
    bbLowerRef.current = bbLowerSeries;

    const loadHistorical = async () => {
      try {
        type BBCandle = CandlestickData & {
          bb_upper?: number;
          bb_lower?: number;
        };

        const historical: BBCandle[] = await BotService.getHistotical(botID);
        if (!historical) return;

        candleSeries.setData(historical);

        bbUpperSeries.setData(
          historical
            .filter((c) => c.bb_upper !== undefined)
            .map((c) => ({ time: c.time, value: c.bb_upper! }))
        );

        bbLowerSeries.setData(
          historical
            .filter((c) => c.bb_lower !== undefined)
            .map((c) => ({ time: c.time, value: c.bb_lower! }))
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
        if (candle.bb_upper) bbUpperSeries.update({ time: candle.time, value: candle.bb_upper });
        if (candle.bb_lower) bbLowerSeries.update({ time: candle.time, value: candle.bb_lower });
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
