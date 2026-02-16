"use client";
import { useRouter } from "next/navigation";
import { DocumentRequest, RequestStatus, User } from "@/types";
import StatusBadge from "../StatusBadge/StatusBadge";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import HighlightText from "../../HighlightText/HighlightText";
import "./RequestCard.css";
import RequestBodyProfessional from "../_components/RequestBodyProfessional";
import RequestBodyClient from "../_components/RequestBodyClient";
import { formatDate } from "@/lib/client/formatDate";

interface RequestProps {
        request: DocumentRequest;
        searchTerm?: string;
        user: User;
}

const Request: React.FC<RequestProps> = ({ request, searchTerm, user }) => {
        const router = useRouter();
        const isOverdue = request.status === "overdue";

        const isScheduledFuture =
                request.is_scheduled &&
                request.scheduled_for &&
                new Date(request.scheduled_for) > new Date();

        const handleViewDetails = () => {
                router.push(`/dashboard/requests/${request.id}`);
        };

        return (
                <div className={`document-request-card ${isOverdue ? "is-overdue" : ""}`}>
                        <div className="request-header">
                                <StatusBadge status={request.status as RequestStatus} />
                                {isScheduledFuture && (
                                        <span
                                                className="scheduled-badge"
                                                title={`Scheduled for ${formatDate(request.scheduled_for!)}`}
                                        >
                                                SCHEDULED
                                        </span>
                                )}
                        </div>
                        <h3 className="request-title">
                                <HighlightText text={request.title} search={searchTerm} />
                        </h3>
                        {user.role === "PROFESSIONAL" && (
                                <RequestBodyProfessional
                                        request={request}
                                        searchTerm={searchTerm}
                                />
                        )}
                        {user.role === "CLIENT" && (
                                <RequestBodyClient request={request} searchTerm={searchTerm} />
                        )}
                        <div className="request-footer">
                                <ButtonPrimary
                                        text="View Details"
                                        variant="ghost"
                                        fullWidth={true}
                                        onClick={handleViewDetails}
                                />
                        </div>
                </div>
        );
};

export default Request;
