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

  return (
    <div className="client-card">
      <h3 className="client-name">{client.email.split("@")[0]}</h3>

      <div className="client-body">
        <p className="client-email">
          <strong>Email:</strong> {client.email}
        </p>
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
