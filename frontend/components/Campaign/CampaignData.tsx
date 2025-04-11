// File: components/Campaign/CampaignInfo.tsx

import { useState, useEffect, useMemo } from "react";
import { Campaign } from "@/types/campaign";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import { CampaignService } from "@/services/campaign";
import { TemplateService } from "@/services/template";

const CampaignSettings = ({ campaign }: { campaign: Campaign }) => {
  const [updatedCampaign, setUpdatedCampaign] = useState(campaign);
  const [emailTemplates, setEmailTemplates] = useState<any[]>([]);
  const [whatsappTemplates, setWhatsappTemplates] = useState<any[]>([]);
  const [isSaving, setIsSaving] = useState(false);
  // Se estiver pendente ou cancelada, permitir edição
  const isEditing = useMemo(
    () => campaign.status === "pendente" || campaign.status === "cancelada",
    [campaign.status]
  );

  useEffect(() => {
    async function fetchTemplates() {
      try {
        setEmailTemplates(await TemplateService.getByChannel("email"));
        setWhatsappTemplates(await TemplateService.getByChannel("whatsapp"));
      } catch (error) {
        console.error("Erro ao carregar templates", error);
        toast.error("Erro ao carregar templates.");
      }
    }

    fetchTemplates();
  }, []);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    setUpdatedCampaign({ ...updatedCampaign, [e.target.name]: e.target.value });
  };

  const handleSave = async () => {
    if (!updatedCampaign.name.trim()) {
      toast.error("O nome da campanha é obrigatório.");
      return;
    }

    if (Object.keys(updatedCampaign.channels).length === 0) {
      toast.error("Pelo menos um canal de comunicação é obrigatório.");
      return;
    }

    setIsSaving(true);

    try {
      await CampaignService.update(updatedCampaign.id, updatedCampaign);
      toast.success("Configurações salvas com sucesso!");
    } catch (error) {
      console.error("Erro ao salvar campanha", error);
      toast.error("Erro ao salvar campanha. Tente novamente.");
    }

    setIsSaving(false);
  };

  return (
    <div className="p-4">
      <h2 className="text-xl font-bold mb-4">📝 Dados da Campanha</h2>

      <div className="grid gap-4">
        <label>Nome:</label>
        <Input
          disabled={!isEditing}
          name="name"
          value={updatedCampaign.name}
          onChange={handleChange}
          placeholder="Nome da Campanha *"
        />

        <label>Descrição:</label>
        <Textarea
          disabled={!isEditing}
          name="description"
          value={updatedCampaign.description || ""}
          onChange={handleChange}
          placeholder="Descrição"
        />

        <label>Template Email:</label>
        <Select
          disabled={!isEditing}
          value={updatedCampaign.channels.email.template}
          onValueChange={(value) =>
            setUpdatedCampaign((prev) => ({
              ...prev!,
              channels: { ...prev!.channels, email: { ...prev!.channels.email, template: value } },
            }))
          }
        >
          <SelectTrigger className="w-full">
            <SelectValue placeholder="Selecionar Template de Email" />
          </SelectTrigger>
          <SelectContent>
            {emailTemplates.length > 0 ? (
              emailTemplates.map((t) => (
                <SelectItem key={t.id} value={t.id}>
                  {t.name}
                </SelectItem>
              ))
            ) : (
              <SelectItem value="none" disabled>
                Nenhum template disponível
              </SelectItem>
            )}
          </SelectContent>
        </Select>

        <label>Template WhatsApp:</label>
        <Select
          disabled={!isEditing}
          value={updatedCampaign.channels.whatsapp?.template}
          onValueChange={(value) =>
            setUpdatedCampaign((prev) => ({
              ...prev!,
              channels: {
                ...prev!.channels,
                whatsapp: { ...prev!.channels.whatsapp, template: value },
              },
            }))
          }
        >
          <SelectTrigger className="w-full">
            <SelectValue placeholder="Selecionar Template de WhatsApp" />
          </SelectTrigger>
          <SelectContent>
            {whatsappTemplates?.length > 0 ? (
              whatsappTemplates?.map((t) => (
                <SelectItem key={t.id} value={t.id}>
                  {t.name}
                </SelectItem>
              ))
            ) : (
              <SelectItem value="none" disabled>
                Nenhum template disponível
              </SelectItem>
            )}
          </SelectContent>
        </Select>
      </div>

      <div className="mt-4">
        <Button onClick={handleSave} disabled={!isEditing || isSaving}>
          {isSaving ? "Salvando..." : "💾 Salvar Dados"}
        </Button>
      </div>
    </div>
  );
};

export default CampaignSettings;
