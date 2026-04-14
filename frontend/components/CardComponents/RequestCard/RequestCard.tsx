"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { Request, RequestStatus, User } from "@/types";
import StatusBadge from "../../Pages/RequestsComponents/StatusBadge";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import { formatDate } from "@/lib/client/formatDate";
import BaseDashboardCard from "@/components/CardComponents/BaseDashboardCard/BaseDashboardCard";
import InfoList from "@/components/CardComponents/InfoList/InfoList";
import InfoItem from "@/components/CardComponents/InfoItem/InfoItem";
import { useRequestActions } from "@/hooks/useRequestActions";
import Modal from "@/components/Modals/Modal";
import HighlightText from "@/components/OtherComponents/HighlightText/HighlightText";

interface RequestProps {
        request: Request;
        user: User;
        searchTerm?: string;
        archived?: boolean;
        cancelled?: boolean;
}

export default function RequestCard({ request, searchTerm, user, archived }: RequestProps) {
        const router = useRouter();
        const { closeReq, reopenReq } = useRequestActions(request.id);
        const [isCloseModalOpen, setIsCloseModalOpen] = useState(false);

        console.log(request);
        const canManage = user.role === "admin" || user.department_id !== null;
        const isOverdue = request.status === "overdue";
        const isScheduledFuture =
                request.is_scheduled &&
                request.scheduled_for &&
                new Date(request.scheduled_for) > new Date();

        return (
                <>
                        <BaseDashboardCard
                                header={
                                        !request.is_cancelled && (
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
                                                                        Programată
                                                                </span>
                                                        )}
                                                </>
                                        )
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
                                                canManage={canManage}
                                                cancelled={request.is_cancelled}
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
                                <InfoList>
                                        <InfoItem
                                                label="Departament"
                                                value={
                                                        <HighlightText
                                                                text={`${request.department_name}`}
                                                                search={searchTerm}
                                                        />
                                                }
                                        />
                                        <InfoItem
                                                label="Solicitant"
                                                value={
                                                        <HighlightText
                                                                text={`${request.assignee_first_name} ${request.assignee_last_name}`}
                                                                search={searchTerm}
                                                        />
                                                }
                                        />
                                        <InfoItem
                                                label="Email"
                                                value={
                                                        <HighlightText
                                                                text={request.assignee_email}
                                                                search={searchTerm}
                                                        />
                                                }
                                        />
                                        {request.due_date && (
                                                <InfoItem
                                                        label="Termen"
                                                        value={formatDate(request.due_date)}
                                                />
                                        )}
                                </InfoList>
                        </BaseDashboardCard>

                        <Modal
                                isOpen={isCloseModalOpen}
                                onClose={() => setIsCloseModalOpen(false)}
                                onConfirm={closeReq}
                                title="Arhivează dosarul"
                        >
                                <p className="modal-text">
                                        Eşti sigur că vrei să arhivezi dosarul{" "}
                                        <strong>{request.title}</strong>?
                                </p>
                                <p className="modal-subtext">
                                        Această acțiune va marca dosarul ca arhivat.
                                        Solicitantul nu va mai putea adăuga documente. Acțiunea
                                        este reversibilă.
                                </p>
                        </Modal>
                </>
        );
}

interface RequestFooterProps {
        archived?: boolean;
        cancelled?: boolean;
        canManage: boolean;
        onView: () => void;
        onClose: () => void;
        onReopen: () => void;
}

const RequestFooter: React.FC<RequestFooterProps> = ({
        archived,
        cancelled,
        canManage,
        onView,
        onClose,
        onReopen,
}) => {
        if (archived) {
                return canManage ? (
                        <ButtonPrimary
                                text="Redeschide dosar"
                                variant="ghost"
                                fullWidth
                                onClick={onReopen}
                        />
                ) : null;
        }

        return (
                <>
                        <ButtonPrimary
                                text="Vezi detalii"
                                variant="ghost"
                                fullWidth
                                onClick={onView}
                        />
                        {canManage && !cancelled && (
                                <ButtonPrimary
                                        text="Arhivează dosar"
                                        variant="ghost"
                                        fullWidth
                                        onClick={onClose}
                                />
                        )}
                </>
        );
};
