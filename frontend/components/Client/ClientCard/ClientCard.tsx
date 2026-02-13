"use client";
import React, { useState } from "react";
import { User } from "@/types";
import "./ClientCard.css";
import ButtonPrimary from "@/components/Buttons/ButtonPrimary/ButtonPrimary";
import HighlightText from "@/components/HighlightText/HighlightText";
import { useRouter } from "next/navigation";
import toast from "react-hot-toast";
import DeactivateClientModal from "./_components/DeactivateClientModal";

interface ClientCardProps {
  client: User;
  searchTerm?: string;
}

const ClientCard: React.FC<ClientCardProps> = ({ client, searchTerm }) => {
  const router = useRouter();
  const [isDeactivateModalOpen, setIsDeactivateModalOpen] = useState(false);

  const handleAddRequest = () => {
    router.push(`/dashboard/clients/${client.id}/add-request`);
  };

  const handleDeactivateClick = () => {
    setIsDeactivateModalOpen(true);
  };

  const handleDeactivateConfirm = async () => {
    const deactivatePromise = fetch(`/api/backend/users/deactivate/${client.id}`, {
      method: "POST",
      credentials: "include",
    }).then(async (res) => {
      if (!res.ok) {
        const error = await res.json();
        console.log(error);
        throw new Error(error.error || "Failed to deactivate user");
      }
    });

    toast.promise(deactivatePromise, {
      loading: "Deactivating user...",
      success: "Deactivated user.",
      error: "Failed to deactivate user",
    });
  };

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
    <>
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
            text="Deactivate their account"
            variant="ghost"
            fullWidth={true}
            onClick={handleDeactivateClick}
          />
        </div>
      </div>

      <DeactivateClientModal
        isOpen={isDeactivateModalOpen}
        onClose={() => setIsDeactivateModalOpen(false)}
        onConfirm={handleDeactivateConfirm}
        clientName={fullName}
      />
    </>
  );
};

export default ClientCard;
