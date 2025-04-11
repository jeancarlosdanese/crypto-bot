// File: pages/dashboard.tsx

import { useUser } from "@/context/UserContext";
import { useRouter } from "next/router";
import { Card } from "@/components/ui/card";
import axios from "axios";
import { useEffect, useState } from "react";
import Spinner from "@/components/Spinner";
import Layout from "@/components/Layout";

const Dashboard = () => {
  const { user, loading } = useUser();
  const [campaigns, setCampaigns] = useState([]);
  const router = useRouter();

  // Redireciona para login caso não esteja autenticado
  useEffect(() => {
    if (!loading && !user) {
      router.push("/auth/login");
    }
  }, [loading, user, router]);

  // Carrega as campanhas ativas
  useEffect(() => {
    if (loading || !user) return;

    axios
      .get(`${process.env.NEXT_PUBLIC_API_URL}/health`, {
        headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
      })
      .then((response) => setCampaigns(response.data))
      .catch((error) => console.error("Erro ao carregar campanhas", error));
  }, [loading, user]);

  // Mostra um spinner enquanto carrega os dados
  if (loading) return <Spinner />;
  if (!user) return null;

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold">Dashboard</h1>
      <div className="grid gap-4 grid-cols-1 md:grid-cols-2 lg:grid-cols-3 mt-4">
        <Card className="w-full p-4 bg-card shadow-lg">
          <h3 className="text-lg font-semibold">Campanhas Ativas</h3>
          <p className="text-2xl font-bold">{campaigns?.length}</p>
        </Card>
      </div>
    </div>
  );
};

// Define o Layout global para a Dashboard
Dashboard.getLayout = (page: JSX.Element) => <Layout>{page}</Layout>;

export default Dashboard;
