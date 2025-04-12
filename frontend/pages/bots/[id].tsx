// pages/bots/[id].tsx

import { useRouter } from "next/router";
import { useUser } from "@/context/UserContext";
import Spinner from "@/components/Spinner";
import TradeChart from "@/components/TradeChart";
import Layout from "@/components/Layout";
import { JSX, useEffect, useState } from "react";
import { Bot } from "@/types/bot";
import { BotService } from "@/services/botService";

const BotChartPage = () => {
  const router = useRouter();
  const { id } = router.query;
  const { user, loading } = useUser();
  const [bot, setBot] = useState<Bot>();

  useEffect(() => {
    if (!user || !id || typeof id !== "string") return;

    const loadBot = async () => {
      const bot = await BotService.getById(id);
      if (bot) {
        setBot(bot);
      } else {
        console.error("Erro ao carregar bot");
      }
    };

    loadBot();
  }, [user, id]);

  if (loading) return <Spinner />;
  if (!user || !id || typeof id !== "string") return null;
  if (!bot) return null;

  const token = localStorage.getItem("token") || "";

  return (
    <div className="p-4 max-w-7xl mx-auto min-h-[500px]">
      <h1 className="text-2xl font-bold mb-4">{bot.symbol}</h1>
      <div className="w-full h-[500px]">
        <TradeChart botID={id} token={token} />
      </div>
    </div>
  );
};

BotChartPage.getLayout = (page: JSX.Element) => <Layout>{page}</Layout>;

export default BotChartPage;
