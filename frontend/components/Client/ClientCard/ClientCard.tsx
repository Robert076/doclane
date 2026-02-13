"use client";
import React from "react";
import { User } from "@/types";
import "./ClientCard.css";
import ButtonPrimary from "@/components/Buttons/ButtonPrimary/ButtonPrimary";
import HighlightText from "@/components/HighlightText/HighlightText";
import { useRouter } from "next/navigation";

interface ClientCardProps {
  client: User;
  searchTerm?: string;
}

const ClientCard: React.FC<ClientCardProps> = ({ client, searchTerm }) => {
  const router = useRouter();

  const handleAddRequest = () => {
    router.push(`/dashboard/clients/${client.id}/add-request`);
  };

  const handleDeactivateClient = () => {};

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString("en-US", {
      year: "numeric",
      month: "long",
      day: "numeric",
    });
  };

  const fullName = `${client.first_name} ${client.last_name}`;

  return (
    <div className="client-card">
      <h3 className="client-name">
        <HighlightText text={fullName} search={searchTerm} />
      </h3>
      <div className="client-body">
        <div className="client-info">
          <p className="client-info-item">
            <span className="client-label">Email:</span>
            <span className="client-value">
              <HighlightText text={client.email} search={searchTerm} />
            </span>
          </p>
          <p className="client-info-item">
            <span className="client-label">Joined at:</span>
            <span className="client-value">{formatDate(client.created_at)}</span>
          </p>
        </div>
      </div>
      <div className="client-footer">
        <ButtonPrimary
          text="New Request"
          variant="ghost"
          fullWidth={true}
          onClick={handleAddRequest}
        />
        <ButtonPrimary
          text="Deactivate thier account"
          variant="ghost"
          fullWidth={true}
          onClick={handleDeactivateClient}
        />
      </div>
    </div>
  );
};

export default ClientCard;
