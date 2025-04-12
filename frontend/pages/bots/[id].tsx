// pages/bots/[id].tsx

import { useRouter } from "next/router";
import { useUser } from "@/context/UserContext";
import Spinner from "@/components/Spinner";
import TradeChart from "@/components/TradeChart";
import Layout from "@/components/Layout";
import { JSX } from "react";

const BotChartPage = () => {
  const router = useRouter();
  const { id } = router.query;
  const { user, loading } = useUser();

  if (loading) return <Spinner />;
  if (!user || !id || typeof id !== "string") return null;

  const token = localStorage.getItem("token") || "";

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-4">Gráfico do Bot</h1>
      <TradeChart botID={id} token={token} />
    </div>
  );
};

BotChartPage.getLayout = (page: JSX.Element) => <Layout>{page}</Layout>;

export default BotChartPage;
