"use client";
import React from "react";
import { User } from "@/types";
import "./ClientCard.css";
import ButtonPrimary from "@/components/Buttons/ButtonPrimary/ButtonPrimary";

interface ClientCardProps {
  client: User;
}

const ClientCard: React.FC<ClientCardProps> = ({ client }) => {
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
          onClick={() => console.log("New request for:", client.id)}
        />
      </div>
    </div>
  );
};

export default ClientCard;
