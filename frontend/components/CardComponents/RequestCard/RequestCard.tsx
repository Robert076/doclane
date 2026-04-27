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
        isStale?: boolean;
        isDueSoon?: boolean;
}

export default function RequestCard({
        request,
        searchTerm,
        user,
        archived,
        isStale,
        isDueSoon,
}: RequestProps) {
        const router = useRouter();
        const { closeReq, reopenReq, claimReq, unclaimReq } = useRequestActions(request.id);
        const [isCloseModalOpen, setIsCloseModalOpen] = useState(false);

        const canManage = user.role === "admin" || user.department_id !== null;
        const isOverdue = request.status === "overdue";
        const isScheduledFuture =
                request.is_scheduled &&
                request.scheduled_for &&
                new Date(request.scheduled_for) > new Date();

        const isClaimed = request.claimed_by !== null && request.claimed_by !== undefined;
        const isClaimedByMe = request.claimed_by === user.id;

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
                                                        {isClaimed && (
                                                                <span className="claimed-badge">
                                                                        {isClaimedByMe
                                                                                ? "Preluat de tine"
                                                                                : `Preluat de ${request.claimed_by_first_name} ${request.claimed_by_last_name}`}
                                                                </span>
                                                        )}
                                                        {canManage && isStale && (
                                                                <span
                                                                        className="stale-badge"
                                                                        title="Dosar nepreluat de peste 7 zile"
                                                                >
                                                                        Nepreluat 7+ zile
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
                                                user={user}
                                                archived={archived}
                                                canManage={canManage}
                                                cancelled={request.is_cancelled}
                                                isClaimed={isClaimed}
                                                isClaimedByMe={isClaimedByMe}
                                                onView={() =>
                                                        router.push(
                                                                `/dashboard/requests/${request.id}`,
                                                        )
                                                }
                                                onClose={() => setIsCloseModalOpen(true)}
                                                onReopen={reopenReq}
                                                onClaim={claimReq}
                                                onUnclaim={unclaimReq}
                                        />
                                }
                                isHighlighted={isOverdue}
                        >
                                <InfoList>
                                        <InfoItem
                                                label="Departament"
                                                value={
                                                        <HighlightText
                                                                text={request.department_name}
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
        user: User;
        archived?: boolean;
        cancelled?: boolean;
        canManage: boolean;
        isClaimed: boolean;
        isClaimedByMe: boolean;
        onView: () => void;
        onClose: () => void;
        onReopen: () => void;
        onClaim: () => void;
        onUnclaim: () => void;
}

const RequestFooter: React.FC<RequestFooterProps> = ({
        user,
        archived,
        cancelled,
        canManage,
        isClaimed,
        isClaimedByMe,
        onView,
        onClose,
        onReopen,
        onClaim,
        onUnclaim,
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

        const isAdmin = user.role === "admin";
        const isMember = user.department_id !== null && !isAdmin;
        const canClaim = isMember && !isClaimed;
        const canUnclaim = isClaimedByMe;
        const canClose = isAdmin || isClaimedByMe;

        if (!canManage || cancelled) {
                return (
                        <ButtonPrimary
                                text="Vezi detalii"
                                variant="ghost"
                                fullWidth
                                onClick={onView}
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
                        {canClaim && (
                                <ButtonPrimary
                                        text="Preia dosar"
                                        variant="ghost"
                                        fullWidth
                                        onClick={onClaim}
                                />
                        )}
                        {canUnclaim && (
                                <ButtonPrimary
                                        text="Renunta"
                                        variant="ghost"
                                        fullWidth
                                        onClick={onUnclaim}
                                />
                        )}
                        {!canClaim && !canUnclaim && (
                                <ButtonPrimary
                                        text={
                                                isClaimed && !isClaimedByMe
                                                        ? "Preluat de altcineva"
                                                        : "Preia dosar"
                                        }
                                        variant="ghost"
                                        fullWidth
                                        disabled={isClaimed && !isClaimedByMe}
                                        onClick={!isClaimed ? onClaim : undefined}
                                />
                        )}
                        {canClose && (
                                <ButtonPrimary
                                        text="Arhiveaza dosar"
                                        variant="ghost"
                                        fullWidth
                                        onClick={onClose}
                                />
                        )}
                </>
        );
};
