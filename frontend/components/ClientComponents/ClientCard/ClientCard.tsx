"use client";
import React, { useState } from "react";
import { User } from "@/types";
import "./ClientCard.css";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import HighlightText from "@/components/OtherComponents/HighlightText/HighlightText";
import { useRouter } from "next/navigation";
import toast from "react-hot-toast";
import DeactivateClientModal from "./_components/DeactivateClientModal";
import ClientInfoItem from "./ClientInfoItem";
import { formatDate } from "@/lib/client/formatDate";
import { deactivateUser } from "@/lib/api/api";
import { UI_TEXT } from "@/locales/ro";

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
                const deactivatePromise = await deactivateUser(+client.id);
                const loadingToast = toast.loading("Deactivating user");

                if (deactivatePromise.success === false) {
                        toast.dismiss(loadingToast);
                        toast.error(deactivatePromise.message);
                } else {
                        toast.dismiss(loadingToast);
                        toast.success(deactivatePromise.message);
                }
        };

        return (
                <>
                        <div className="client-card">
                                <h3 className="client-name">
                                        <HighlightText
                                                text={`${client.first_name} ${client.last_name}`}
                                                search={searchTerm}
                                        />
                                </h3>
                                <div className="client-body">
                                        <div className="client-info">
                                                {ClientInfoItem(
                                                        UI_TEXT.client.card.clientEmail,
                                                        client.email,
                                                        searchTerm,
                                                )}
                                                {ClientInfoItem(
                                                        UI_TEXT.client.card.joinedAt,
                                                        formatDate(client.created_at),
                                                        searchTerm,
                                                )}
                                        </div>
                                </div>
                                <div className="client-footer">
                                        <ButtonPrimary
                                                text={UI_TEXT.client.actions.newRequest}
                                                variant="ghost"
                                                fullWidth={true}
                                                onClick={handleAddRequest}
                                        />
                                        <ButtonPrimary
                                                text={UI_TEXT.client.actions.deactivateAccount}
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
                                clientName={`${client.first_name} ${client.last_name}`}
                        />
                </>
        );
};

export default ClientCard;
