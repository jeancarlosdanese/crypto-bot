// File: components/Campaign/CampaignStatusBadge.tsx

import { Clock, RefreshCw, Send, CheckCircle, XCircle } from "lucide-react";

const statusIcons = {
  pendente: { icon: <Clock className="w-4 h-4 text-yellow-600" />, label: "Pendente" },
  processando: {
    icon: <RefreshCw className="w-4 h-4 text-orange-500 animate-spin" />,
    label: "Processando",
  },
  enviando: { icon: <Send className="w-4 h-4 text-blue-500" />, label: "Enviando" },
  concluida: { icon: <CheckCircle className="w-4 h-4 text-green-600" />, label: "Concluída" },
  cancelada: { icon: <XCircle className="w-4 h-4 text-red-600" />, label: "Cancelada" },
};

const CampaignStatusBadge = ({ status }: { status: string }) => {
  const statusData = statusIcons[status as keyof typeof statusIcons] || {
    icon: <Clock className="w-4 h-4 text-gray-600" />,
    label: "Desconhecido",
  };

  return (
    <div className="flex items-center gap-1 text-gray-700">
      {statusData.icon}
      <span className="text-xs">{statusData.label}</span>
    </div>
  );
};

export default CampaignStatusBadge;
