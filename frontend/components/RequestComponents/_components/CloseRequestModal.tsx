"use client";
import React from "react";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import "./CloseRequestModal.css";

interface CloseRequestModalProps {
        isOpen: boolean;
        onClose: () => void;
        onConfirm: () => void;
        requestTitle: string;
}

const CloseRequestModal: React.FC<CloseRequestModalProps> = ({
        isOpen,
        onClose,
        onConfirm,
        requestTitle,
}) => {
        if (!isOpen) return null;

        const handleConfirm = () => {
                onConfirm();
                onClose();
        };

        return (
                <div className="close-request-modal-overlay" onClick={onClose}>
                        <div
                                className="close-request-modal-content"
                                onClick={(e) => e.stopPropagation()}
                        >
                                <div className="close-request-modal-header">
                                        <h3>Close Request</h3>
                                        <button
                                                className="close-request-modal-close"
                                                onClick={onClose}
                                        >
                                                ×
                                        </button>
                                </div>
                                <div className="close-request-modal-body">
                                        <p className="close-request-warning">
                                                Are you sure you want to close{" "}
                                                <strong>{requestTitle}</strong>?
                                        </p>
                                        <p className="close-request-description">
                                                This action will mark the request as closed.
                                                You will still be able to view the documents,
                                                however your client cannot upload anymore.
                                        </p>
                                </div>
                                <div className="close-request-modal-footer">
                                        <ButtonPrimary
                                                text="Cancel"
                                                variant="ghost"
                                                onClick={onClose}
                                        />
                                        <ButtonPrimary
                                                text="Close Request"
                                                variant="primary"
                                                onClick={handleConfirm}
                                        />
                                </div>
                        </div>
                </div>
        );
};

export default CloseRequestModal;
