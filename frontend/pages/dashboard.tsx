// pages/dashboard.tsx

import { useUser } from "@/context/UserContext";
import { useRouter } from "next/router";
import axios from "axios";
import { JSX, useEffect, useState } from "react";
import Spinner from "@/components/Spinner";
import Layout from "@/components/Layout";
import { Bot } from "@/types/bot";
import { BotCard } from "@/components/BotCard";

const Dashboard = () => {
  const { user, loading } = useUser();
  const [bots, setBots] = useState<Bot[]>([]);
  const router = useRouter();

  useEffect(() => {
    if (!loading && !user) {
      router.push("/auth/login");
    }
  }, [loading, user, router]);

  useEffect(() => {
    if (loading || !user) return;

    axios
      .get(`${process.env.NEXT_PUBLIC_API_URL}/bots`, {
        headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
      })
      .then((response) => setBots(response.data))
      .catch((error) => console.error("Erro ao carregar bots", error));
  }, [loading, user]);

  if (loading) return <Spinner />;
  if (!user) return null;

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold">Bots Ativos</h1>
      {bots.length > 0 && (
        <div className="grid gap-4 grid-cols-1 md:grid-cols-2 lg:grid-cols-3 mt-4">
          {bots.map((bot) => (
            <BotCard key={bot.id} bot={bot} />
          ))}
        </div>
      )}
    </div>
  );
};

Dashboard.getLayout = (page: JSX.Element) => <Layout>{page}</Layout>;

export default Dashboard;
