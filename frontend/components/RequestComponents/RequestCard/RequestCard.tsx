"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { DocumentRequest, RequestStatus, User } from "@/types";
import StatusBadge from "../StatusBadge/StatusBadge";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import HighlightText from "../../OtherComponents/HighlightText/HighlightText";
import "./RequestCard.css";
import RequestBodyProfessional from "../_components/RequestBodyProfessional";
import RequestBodyClient from "../_components/RequestBodyClient";
import { formatDate } from "@/lib/client/formatDate";
import { closeRequest } from "@/lib/api/api";
import toast from "react-hot-toast";
import CloseRequestModal from "../_components/CloseRequestModal";
import { UI_TEXT } from "@/locales/ro";

interface RequestProps {
        request: DocumentRequest;
        searchTerm?: string;
        user: User;
}

const Request: React.FC<RequestProps> = ({ request, searchTerm, user }) => {
        const router = useRouter();
        const [isCloseModalOpen, setIsCloseModalOpen] = useState(false);

        const isOverdue = request.status === "overdue";
        const isScheduledFuture =
                request.is_scheduled &&
                request.scheduled_for &&
                new Date(request.scheduled_for) > new Date();

        const handleViewDetails = () => {
                router.push(`/dashboard/requests/${request.id}`);
        };

        const handleCloseRequest = async () => {
                toast.promise(closeRequest(request.id), {
                        loading: "Closing request...",
                        success: (response) => {
                                if (!response.success) throw new Error(response.message);
                                return response.message || "Request closed successfully";
                        },
                        error: (err) => err.message || "Something went wrong",
                });
        };

        return (
                <>
                        <div
                                className={`document-request-card ${isOverdue ? "is-overdue" : ""}`}
                        >
                                <div className="request-header">
                                        <StatusBadge
                                                status={request.status as RequestStatus}
                                        />
                                        {isScheduledFuture && (
                                                <span
                                                        className="scheduled-badge"
                                                        title={`Programată pentru ${formatDate(request.scheduled_for!)}`}
                                                >
                                                        SCHEDULED
                                                </span>
                                        )}
                                </div>
                                <h3 className="request-title">
                                        <HighlightText
                                                text={request.title}
                                                search={searchTerm}
                                        />
                                </h3>
                                {user.role === "PROFESSIONAL" && (
                                        <RequestBodyProfessional
                                                request={request}
                                                searchTerm={searchTerm}
                                        />
                                )}
                                {user.role === "CLIENT" && (
                                        <RequestBodyClient
                                                request={request}
                                                searchTerm={searchTerm}
                                        />
                                )}
                                <div className="request-footer">
                                        <ButtonPrimary
                                                text={UI_TEXT.request.actions.viewDetails}
                                                variant="ghost"
                                                fullWidth={true}
                                                onClick={handleViewDetails}
                                        />
                                        <ButtonPrimary
                                                text={UI_TEXT.request.actions.closeRequest}
                                                variant="ghost"
                                                fullWidth={true}
                                                onClick={() => setIsCloseModalOpen(true)}
                                        />
                                </div>
                        </div>

                        <CloseRequestModal
                                isOpen={isCloseModalOpen}
                                onClose={() => setIsCloseModalOpen(false)}
                                onConfirm={handleCloseRequest}
                                requestTitle={request.title}
                        />
                </>
        );
};

export default Request;
