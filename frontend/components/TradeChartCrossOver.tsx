// components/TradeChartCrossOver.tsx

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

export default function TradeChartCrossOver({ botID, token }: Props) {
  const chartRef = useRef<HTMLDivElement | null>(null);
  const chartApiRef = useRef<IChartApi | null>(null);
  const candleSeriesRef = useRef<ISeriesApi<"Candlestick"> | null>(null);
  const markersRef = useRef<Decision[]>([]);
  const ma9SeriesRef = useRef<ISeriesApi<"Line"> | null>(null);
  const ma26SeriesRef = useRef<ISeriesApi<"Line"> | null>(null);

  useEffect(() => {
    if (!chartRef.current) return;

    const isDark = document.documentElement.classList.contains("dark");

    const chart = createChart(chartRef.current, {
      width: chartRef.current.clientWidth,
      // height: 400,
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
        rightOffset: 20, // 🔥 mostra 5 candles de margem à direita
      },
    });

    // 🌓 Detectar mudanças no tema e aplicar ao gráfico
    const modeThemeObserver = new MutationObserver(() => {
      const isDark = document.documentElement.classList.contains("dark");

      chart.applyOptions({
        layout: {
          background: { type: ColorType.Solid, color: isDark ? "#121212" : "#ffffff" },
          textColor: isDark ? "#f5f5f5" : "#333",
        },
        grid: {
          vertLines: { color: isDark ? "#444" : "#eee" },
          horzLines: { color: isDark ? "#444" : "#eee" },
        },
        timeScale: {
          borderColor: isDark ? "#666" : "#ccc",
          timeVisible: true,
          secondsVisible: true,
        },
      });

      // 🔥 Atualizar cor das médias
      ma9SeriesRef.current?.applyOptions({
        color: isDark ? "#fbbf24" : "orange", // Amarelo no dark, laranja no light
      });

      ma26SeriesRef.current?.applyOptions({
        color: isDark ? "#9370DB" : "#4682B4", // Azul claro no dark, steelblue no light
      });
    });

    // Observar mudanças na classe `dark` no elemento HTML
    modeThemeObserver.observe(document.documentElement, {
      attributes: true,
      attributeFilter: ["class"],
    });

    // 🖥️ Redimensionar automaticamente com a janela
    const resizeObserver = new ResizeObserver(() => {
      if (chartRef.current) {
        chart.resize(chartRef.current.clientWidth, chartRef.current.clientHeight);
      }
    });

    const candleSeries = chart.addCandlestickSeries();
    const ma9Series = chart.addLineSeries({ color: "orange", lineWidth: 2 });
    const ma26Series = chart.addLineSeries({ color: "#9370DB", lineWidth: 2 }); // SteelBlue

    chartApiRef.current = chart;
    candleSeriesRef.current = candleSeries;
    ma9SeriesRef.current = ma9Series;
    ma26SeriesRef.current = ma26Series;

    const tooltip = document.getElementById("tooltip");

    chart.subscribeCrosshairMove((param) => {
      if (!param?.time || !tooltip) {
        tooltip!.style.display = "none";
        return;
      }

      const date = new Date((param.time as number) * 1000);
      const formattedTime = date.toLocaleString("pt-BR", {
        timeZone: "America/Sao_Paulo",
        hour: "2-digit",
        minute: "2-digit",
        second: "2-digit",
      });

      const prices = (param as any).seriesData as Map<any, any>;
      if (!prices || prices.size === 0) {
        tooltip.style.display = "none";
        return;
      }

      const candle = prices.get(candleSeries);
      const ma9 = prices.get(ma9Series);
      const ma26 = prices.get(ma26Series);

      if (!candle) return;

      tooltip.style.display = "block";
      tooltip.innerHTML = `
        <div><strong>🕒 Horário:</strong> ${formattedTime}</div>
        <div><strong>Open:</strong> ${candle.open}</div>
        <div><strong>High:</strong> ${candle.high}</div>
        <div><strong>Low:</strong> ${candle.low}</div>
        <div><strong>Close:</strong> ${candle.close}</div>
        ${
          ma9
            ? `<div style="color:${ma9SeriesRef.current?.options().color};"><strong>MA9:</strong> ${ma9.value.toFixed(2)}</div>`
            : ""
        }
        ${
          ma26
            ? `<div style="color:${ma26SeriesRef.current?.options().color};"><strong>MA26:</strong> ${ma26.value.toFixed(2)}</div>`
            : ""
        }
      `;

      const x = param.point?.x ?? 0;
      const y = param.point?.y ?? 0;

      tooltip.style.left = `${x + 20}px`;
      tooltip.style.top = `${y + 20}px`;
    });

    const loadHistorical = async () => {
      try {
        type HistoricalCandle = CandlestickData & { ma9?: number; ma26?: number };

        const historical: HistoricalCandle[] = await BotService.getHistotical(botID);
        if (!historical) {
          console.error("Erro ao carregar dados históricos");
          return;
        }

        candleSeries.setData(historical);
        ma9Series.setData(
          historical
            .filter((c) => c.ma9 !== undefined)
            .map((c) => ({ time: c.time, value: c.ma9! }))
        );
        ma26Series.setData(
          historical
            .filter((c) => c.ma26 !== undefined)
            .map((c) => ({ time: c.time, value: c.ma26! }))
        );
      } catch (error) {
        console.error("Erro ao carregar dados históricos", error);
      }
    };

    loadHistorical();

    const socket = new WebSocket(`ws://localhost:8080/ws/${botID}?token=${token}`);

    socket.onopen = () => {
      console.log("✅ WebSocket aberto", { botID });
    };

    socket.onclose = () => {
      console.log("⚠️ WebSocket connection closed");
    };

    socket.onerror = (error) => {
      console.error("❌ WebSocket erro", error);
    };

    socket.onmessage = (event) => {
      const msg = JSON.parse(event.data);
      if (msg.type === "candle") {
        const candle = msg.data;
        candleSeries.update(candle);
        chart.timeScale().scrollToRealTime(); // 🔄 rola para o final
        if (candle.ma9 && candle.ma26) {
          ma9SeriesRef.current?.update({ time: candle.time, value: candle.ma9 });
          ma26SeriesRef.current?.update({ time: candle.time, value: candle.ma26 });
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
      modeThemeObserver.disconnect();
      resizeObserver.disconnect();
    };
  }, [botID, token]);

  return (
    <div ref={chartRef} className="relative w-full min-h-[500px] h-[500px]">
      <div
        id="tooltip"
        className="absolute z-10 bg-white text-black border border-gray-300 shadow-md text-sm p-2 rounded hidden dark:bg-gray-800 dark:text-white dark:border-gray-600"
        style={{ pointerEvents: "none" }}
      ></div>
    </div>
  );
}
