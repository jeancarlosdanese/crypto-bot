// components/BotCard.tsx

import { Bot } from "@/types/bot";
import { Badge } from "@/components/ui/badge";
import { Card } from "@/components/ui/card";
import Link from "next/link";

type Props = {
  bot: Bot;
};

export function BotCard({ bot }: Props) {
  return (
    <Card className="p-4 shadow-md bg-white dark:bg-zinc-900">
      <div className="flex justify-between items-center mb-2">
        <h2 className="text-xl font-bold">{bot.symbol}</h2>
        <span className="flex items-center gap-2">
          <span
            className={`w-2 h-2 rounded-full ${bot.active ? "bg-green-500" : "bg-red-500"}`}
          ></span>
          {bot.active ? "ativo" : "inativo"}
        </span>
      </div>
      <p className="text-sm text-muted-foreground mb-1">
        Estratégia: <strong>{bot.strategy_name}</strong>
      </p>
      <p className="text-sm text-muted-foreground mb-1">
        Intervalo: <strong>{bot.interval}</strong>
      </p>
      <p className="text-sm text-muted-foreground">
        Autonomia: {bot.autonomous ? "Automático" : "Manual"}
      </p>
      <Link href={`/bots/${bot.id}`}>
        <button className="mt-4 w-full bg-blue-600 text-white py-2 px-4 rounded hover:bg-blue-700 transition">
          📈 Ver Gráfico
        </button>
      </Link>
    </Card>
  );
}
