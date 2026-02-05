"use client";
import React from "react";
import { User } from "@/types";
import "./ClientCard.css";
import ButtonPrimary from "@/components/Buttons/ButtonPrimary/ButtonPrimary";
import { useRouter } from "next/navigation";

interface ClientCardProps {
  client: User;
}

const ClientCard: React.FC<ClientCardProps> = ({ client }) => {
  const router = useRouter();

  const handleAddRequest = () => {
    router.push(`/dashboard/clients/${client.id}/add-request`);
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString("en-US", {
      year: "numeric",
      month: "long",
      day: "numeric",
    });
  };

  return (
    <div className="client-card">
      <h3 className="client-name">
        {client.first_name} {client.last_name}
      </h3>

      <div className="client-body">
        <div className="client-info">
          <p className="client-info-item">
            <span className="client-label">Email:</span>
            <span className="client-value">{client.email}</span>
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
      </div>
    </div>
  );
};

export default ClientCard;
