// File: components/Campaign/CampaignActions.tsx

import { Button } from "@/components/ui/button";
import { useRouter } from "next/router";
import { ConfirmDelete } from "@/components/ui/ConfirmDelete";
import { CampaignService } from "@/services/campaign";

const CampaignActions = ({ campaign }: { campaign: any }) => {
  const router = useRouter();

  const deleteCampaign = async (id: string) => {
    try {
      await CampaignService.delete(id);
      router.push("/campaigns");
    } catch (error) {
      console.error("Erro ao excluir campanha", error);
    }
  };

  return (
    <div className="grid grid-cols-2 gap-2">
      <Button variant="outline" size="sm" onClick={() => router.push(`/campaigns/${campaign.id}`)}>
        Ver detalhes
      </Button>
      <Button
        variant="outline"
        size="sm"
        onClick={() => router.push(`/campaigns/manage/${campaign.id}`)}
      >
        Editar
      </Button>
      {campaign.status === "pendente" && (
        <ConfirmDelete onConfirm={() => deleteCampaign(campaign.id)} entityName="Campanha" />
      )}
    </div>
  );
};

export default CampaignActions;
