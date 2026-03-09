"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { DocumentRequest, RequestStatus, User } from "@/types";
import StatusBadge from "../../Pages/RequestsComponents/StatusBadge";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import HighlightText from "../../OtherComponents/HighlightText/HighlightText";
import { formatDate } from "@/lib/client/formatDate";
import BaseDashboardCard from "@/components/CardComponents/BaseDashboardCard/BaseDashboardCard";
import { useRequestActions } from "@/hooks/useRequestActions";
import RequestBodyProfessional from "@/components/Pages/RequestsComponents/RequestBodyProfessional";
import RequestBodyClient from "@/components/Pages/RequestsComponents/RequestBodyClient";
import Modal from "@/components/Modals/Modal";

interface RequestProps {
        request: DocumentRequest;
        user: User;
        searchTerm?: string;
        archived?: boolean;
}

export default function RequestCard({ request, searchTerm, user, archived }: RequestProps) {
        const router = useRouter();
        const { closeReq, reopenReq } = useRequestActions(request.id);
        const [isCloseModalOpen, setIsCloseModalOpen] = useState(false);

        const isOverdue = request.status === "overdue";
        const isScheduledFuture =
                request.is_scheduled &&
                request.scheduled_for &&
                new Date(request.scheduled_for) > new Date();

        return (
                <>
                        <BaseDashboardCard
                                header={
                                        <>
                                                <StatusBadge
                                                        status={
                                                                request.status as RequestStatus
                                                        }
                                                />
                                                {isScheduledFuture && (
                                                        <span
                                                                className="scheduled-badge"
                                                                title={`Programată pentru ${formatDate(request.scheduled_for!)}`}
                                                        >
                                                                SCHEDULED
                                                        </span>
                                                )}
                                        </>
                                }
                                title={
                                        <HighlightText
                                                text={request.title}
                                                search={searchTerm}
                                        />
                                }
                                footer={
                                        <RequestFooter
                                                archived={archived}
                                                onView={() =>
                                                        router.push(
                                                                `/dashboard/requests/${request.id}`,
                                                        )
                                                }
                                                onClose={() => setIsCloseModalOpen(true)}
                                                onReopen={reopenReq}
                                        />
                                }
                                isHighlighted={isOverdue}
                        >
                                {user.role === "PROFESSIONAL" ? (
                                        <RequestBodyProfessional
                                                request={request}
                                                searchTerm={searchTerm}
                                        />
                                ) : (
                                        <RequestBodyClient
                                                request={request}
                                                searchTerm={searchTerm}
                                        />
                                )}
                        </BaseDashboardCard>

                        <Modal
                                isOpen={isCloseModalOpen}
                                onClose={() => setIsCloseModalOpen(false)}
                                onConfirm={closeReq}
                                title={"Arhivează dosarul"}
                        >
                                <ArchiveRequestContent title={request.title} />
                        </Modal>
                </>
        );
}

function ArchiveRequestContent({ title }: { title: string }) {
        return (
                <>
                        <p className="modal-text">
                                Eşti sigur că vrei să arhivezi dosarul <strong>{title}</strong>
                                ?
                        </p>

                        <p className="modal-subtext">
                                Această acțiune va marca dosarul ca arhivat. Solicitantul nu va
                                mai putea adăuga documente. Acțiunea este reversibilă.
                        </p>
                </>
        );
}

interface RequestFooterProps {
        archived?: boolean;
        onView: () => void;
        onClose: () => void;
        onReopen: () => void;
}

const RequestFooter: React.FC<RequestFooterProps> = ({
        archived,
        onView,
        onClose,
        onReopen,
}) => {
        if (archived) {
                return (
                        <ButtonPrimary
                                text="Redeschide dosar"
                                variant="ghost"
                                fullWidth
                                onClick={onReopen}
                        />
                );
        }

        return (
                <>
                        <ButtonPrimary
                                text="Vezi detalii"
                                variant="ghost"
                                fullWidth
                                onClick={onView}
                        />
                        <ButtonPrimary
                                text="Închide dosar"
                                variant="ghost"
                                fullWidth
                                onClick={onClose}
                        />
                </>
        );
};
