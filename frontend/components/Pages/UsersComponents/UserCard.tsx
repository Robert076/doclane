"use client";
import { useState } from "react";
import { User } from "@/types";
import BaseDashboardCard from "@/components/CardComponents/BaseDashboardCard/BaseDashboardCard";
import InfoList from "@/components/CardComponents/InfoList/InfoList";
import InfoItem from "@/components/CardComponents/InfoItem/InfoItem";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import Modal from "@/components/Modals/Modal";
import { formatDate } from "@/lib/client/formatDate";

interface Props {
        user: User;
        footer?: React.ReactNode;
}

export default function UserCard({ user, footer }: Props) {
        const [isModalOpen, setIsModalOpen] = useState(false);

        return (
                <>
                        <BaseDashboardCard
                                title={`${user.first_name} ${user.last_name}`}
                                isHighlighted={!user.is_active}
                                footer={
                                        <>
                                                <ButtonPrimary
                                                        text="Vezi detalii"
                                                        variant="ghost"
                                                        fullWidth
                                                        onClick={() => setIsModalOpen(true)}
                                                />
                                                {footer}
                                        </>
                                }
                        >
                                <InfoList>
                                        <InfoItem label="Email" value={user.email} />
                                </InfoList>
                        </BaseDashboardCard>

                        <Modal
                                isOpen={isModalOpen}
                                onClose={() => setIsModalOpen(false)}
                                title={`${user.first_name} ${user.last_name}`}
                                hideFooter
                        >
                                <InfoList>
                                        <InfoItem label="Email" value={user.email} />
                                        <InfoItem label="Rol" value={user.role} />
                                        <InfoItem
                                                label="Status"
                                                value={user.is_active ? "Activ" : "Dezactivat"}
                                        />
                                        <InfoItem
                                                label="Înregistrat la"
                                                value={formatDate(user.created_at)}
                                        />
                                        {user.last_notified && (
                                                <InfoItem
                                                        label="Ultima notificare"
                                                        value={formatDate(user.last_notified)}
                                                />
                                        )}
                                </InfoList>
                        </Modal>
                </>
        );
}
