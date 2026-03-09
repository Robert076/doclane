"use client";
import React, { useState } from "react";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import UploadDocumentButton from "@/components/Pages/RequestsComponents/UploadDocumentButton";
import RejectDocumentModal from "./RejectDocumentModal";
import { ExpectedDocumentStatus } from "@/types";

interface DocumentSlotActionsProps {
        isProfessional: boolean;
        status: ExpectedDocumentStatus;
        hasFiles: boolean;
        isLoading: boolean;
        requestId: number;
        expectedDocumentId: number;
        documentTitle: string;
        onApprove: () => void;
        onReject: (reason: string) => void;
        onReset: () => void;
}

export default function DocumentSlotActions({
        isProfessional,
        status,
        hasFiles,
        isLoading,
        requestId,
        expectedDocumentId,
        documentTitle,
        onApprove,
        onReject,
        onReset,
}: DocumentSlotActionsProps) {
        const [isRejectModalOpen, setIsRejectModalOpen] = useState(false);

        if (!isProfessional) {
                if (status === "approved") return null;
                return (
                        <UploadDocumentButton
                                requestId={requestId.toString()}
                                expectedDocumentId={expectedDocumentId}
                        />
                );
        }

        // Professional views
        if (status === "approved" || status === "rejected") {
                return (
                        <ButtonPrimary
                                text={isLoading ? "..." : "Anulează"}
                                variant="ghost"
                                onClick={onReset}
                        />
                );
        }

        if (hasFiles) {
                return (
                        <>
                                <ButtonPrimary
                                        text={isLoading ? "..." : "Refuză"}
                                        variant="ghost"
                                        onClick={() => setIsRejectModalOpen(true)}
                                />
                                <ButtonPrimary
                                        text={isLoading ? "..." : "Aprobă"}
                                        variant="primary"
                                        onClick={onApprove}
                                />
                                <RejectDocumentModal
                                        isOpen={isRejectModalOpen}
                                        onClose={() => setIsRejectModalOpen(false)}
                                        onConfirm={onReject}
                                        documentTitle={documentTitle}
                                />
                        </>
                );
        }

        return null;
}
