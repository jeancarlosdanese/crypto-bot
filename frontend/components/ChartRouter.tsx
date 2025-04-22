// components/ChartRouter.tsx

"use client";

import { Bot } from "@/types/bot";
import TradeChartCrossOver from "@/components/TradeChartCrossOver";
import TradeChartEMAFan from "./TradeChartEMAFan";
import TradeChartRSI from "./TradeChartRSI";
import TradeChartMACD from "./TradeChartMACD";
import TradeChartVolumeSpike from "./TradeChartVolumeSpike";
import TradeChartBollinger from "./TradeChartBollinger";

type Props = {
  bot: Bot;
  token: string;
};

export default function ChartRouter({ bot, token }: Props) {
  switch (bot.strategy_name) {
    case "CROSSOVER":
    case "CROSSOVER_ADVANCED":
      return <TradeChartCrossOver botID={bot.id} token={token} />;

    case "EMA_FAN":
      return <TradeChartEMAFan botID={bot.id} token={token} />;

    case "RSI2":
      return <TradeChartRSI botID={bot.id} token={token} />;

    case "MACD_CROSS":
      return <TradeChartMACD botID={bot.id} token={token} />;

    case "VOLUME_SPIKE":
      return <TradeChartVolumeSpike botID={bot.id} token={token} />;

    case "BB_REBOUND":
      return <TradeChartBollinger botID={bot.id} token={token} />;

    default:
      return (
        <div className="text-sm text-red-600 dark:text-red-400">
          Estratégia não suportada no gráfico ainda: <strong>{bot.strategy_name}</strong>
        </div>
      );
  }
}
