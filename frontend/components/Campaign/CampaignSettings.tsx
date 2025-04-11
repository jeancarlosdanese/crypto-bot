// File: components/Campaign/CampaignSettings.tsx

import { useState, useEffect, useMemo } from "react";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import { CampaignService } from "@/services/campaign";
import { se } from "date-fns/locale";

const CampaignSettings = ({ campaignId, status }: { campaignId: string; status: string }) => {
  const [settings, setSettings] = useState<any | null>(null);
  const [isSaving, setIsSaving] = useState(false);
  const [isCloning, setIsCloning] = useState(false);
  // Se estiver pendente ou cancelada, permitir edição
  const isEditing = useMemo(() => status === "pendente" || status === "cancelada", [status]);

  useEffect(() => {
    async function fetchSettings() {
      try {
        const data = await CampaignService.getSettings(campaignId);
        setSettings(data || {});
      } catch (error) {
        console.warn("Nenhuma configuração encontrada.");
        setSettings({});
      }
    }

    fetchSettings();
  }, [campaignId]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    setSettings({ ...settings, [e.target.name]: e.target.value });
  };

  const handleSave = async () => {
    if (!settings.brand || !settings.subject || !settings.email_from || !settings.whatsapp_from) {
      toast.error("Todos os campos obrigatórios devem ser preenchidos.");
      return;
    }

    setIsSaving(true);

    try {
      if (settings.id) {
        await CampaignService.updateSettings(campaignId, settings);
        toast.success("Configuração atualizada com sucesso!");
      } else {
        await CampaignService.createSettings(campaignId, { ...settings, campaign_id: campaignId });
        toast.success("Configuração criada com sucesso!");
      }
    } catch (error) {
      console.error("Erro ao salvar configurações", error);
      toast.error("Erro ao salvar configurações. Tente novamente.");
    }

    setIsSaving(false);
  };

  const handleCloneSettings = async () => {
    setIsCloning(true);
    try {
      const clonedSettings = await CampaignService.cloneLastSettings(campaignId);
      setSettings(clonedSettings);
      toast.success("Configuração clonada com sucesso!");
    } catch (error) {
      console.error("Erro ao clonar configurações", error);
      toast.error("Erro ao clonar configurações.");
    }
    setIsCloning(false);
  };

  if (!settings) return <p className="text-gray-500 text-center">Carregando configurações...</p>;

  return (
    <div className="p-4">
      <h2 className="text-xl font-bold mb-4">⚙️ Configuração de Envio</h2>

      <div className="grid gap-4">
        <label>Marca:</label>
        <Input
          disabled={!isEditing}
          name="brand"
          value={settings.brand || ""}
          onChange={handleChange}
          placeholder="Nome da marca"
        />

        <label>Assunto do E-mail:</label>
        <Input
          disabled={!isEditing}
          name="subject"
          value={settings.subject || ""}
          onChange={handleChange}
          placeholder="Assunto do e-mail"
        />

        <label>Tom de Voz:</label>
        <Input
          disabled={!isEditing}
          name="tone"
          value={settings.tone || ""}
          onChange={handleChange}
          placeholder="Tom de voz"
        />

        <label>Remetente do E-mail:</label>
        <Input
          disabled={!isEditing}
          name="email_from"
          value={settings.email_from || ""}
          onChange={handleChange}
          placeholder="E-mail do remetente"
        />

        <label>Resposta do E-mail:</label>
        <Input
          disabled={!isEditing}
          name="email_reply"
          value={settings.email_reply || ""}
          onChange={handleChange}
          placeholder="E-mail para respostas"
        />

        <label>Rodapé do E-mail:</label>
        <Textarea
          disabled={!isEditing}
          name="email_footer"
          value={settings.email_footer || ""}
          onChange={handleChange}
          placeholder="Rodapé do e-mail"
        />

        <label>Instruções do E-mail:</label>
        <Textarea
          disabled={!isEditing}
          name="email_instructions"
          value={settings.email_instructions || ""}
          onChange={handleChange}
          placeholder="Instruções do e-mail"
        />

        <label>Remetente do WhatsApp:</label>
        <Input
          disabled={!isEditing}
          name="whatsapp_from"
          value={settings.whatsapp_from || ""}
          onChange={handleChange}
          placeholder="Número do WhatsApp"
        />

        <label>Resposta do WhatsApp:</label>
        <Input
          disabled={!isEditing}
          name="whatsapp_reply"
          value={settings.whatsapp_reply || ""}
          onChange={handleChange}
          placeholder="Número para resposta"
        />

        <label>Rodapé do WhatsApp:</label>
        <Textarea
          disabled={!isEditing}
          name="whatsapp_footer"
          value={settings.whatsapp_footer || ""}
          onChange={handleChange}
          placeholder="Rodapé do WhatsApp"
        />

        <label>Instruções do WhatsApp:</label>
        <Textarea
          disabled={!isEditing}
          name="whatsapp_instructions"
          value={settings.whatsapp_instructions || ""}
          onChange={handleChange}
          placeholder="Instruções do WhatsApp"
        />
      </div>

      <div className="mt-4 flex gap-3">
        <Button onClick={handleSave} disabled={!isEditing || isSaving}>
          {isSaving ? "Salvando..." : "💾 Salvar Configuração"}
        </Button>

        {!settings.id && !isCloning && (
          <Button
            variant="default"
            onClick={handleCloneSettings}
            disabled={!isEditing || settings.id || isCloning}
          >
            {isCloning ? "Clonando..." : "🔄 Clonar Última Configuração"}
          </Button>
        )}
      </div>
    </div>
  );
};

export default CampaignSettings;
