// File: components/Campaign/CampaignDetails.tsx

import { Campaign } from "@/types/campaign";

const CampaignDetails = ({ campaign }: { campaign: Campaign }) => {
  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold">Detalhes da Campanha</h1>
      <p>Status: {campaign.status}</p>
    </div>
  );
};

export default CampaignDetails;
