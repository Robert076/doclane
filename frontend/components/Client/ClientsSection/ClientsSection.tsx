import React from "react";
import { User } from "@/types";
import ClientCard from "../ClientCard/ClientCard";
import "./ClientsSection.css";

interface ClientsSectionProps {
  clients: User[];
}

const ClientsSection: React.FC<ClientsSectionProps> = ({ clients }) => {
  if (clients.length === 0) {
    return (
      <div className="clients-empty">
        <p>No clients found. Start by adding your first client.</p>
      </div>
    );
  }

  return (
    <div className="clients-grid">
      {clients.map((client) => (
        <ClientCard key={client.id} client={client} />
      ))}
    </div>
  );
};

export default ClientsSection;
