"use client";
import React from "react";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import "./CloseRequestModal.css";
import { UI_TEXT } from "@/locales/ro";

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
                                        <h3>{UI_TEXT.modals.closeRequest.title}</h3>
                                        <button
                                                className="close-request-modal-close"
                                                onClick={onClose}
                                        >
                                                ×
                                        </button>
                                </div>
                                <div className="close-request-modal-body">
                                        <p className="close-request-warning">
                                                {UI_TEXT.modals.closeRequest.subtitle1}
                                                <strong>{requestTitle}</strong>?
                                        </p>
                                        <p className="close-request-description">
                                                {UI_TEXT.modals.closeRequest.subtitle2}
                                        </p>
                                </div>
                                <div className="close-request-modal-footer">
                                        <ButtonPrimary
                                                text={UI_TEXT.common.cancel}
                                                variant="ghost"
                                                onClick={onClose}
                                        />
                                        <ButtonPrimary
                                                text={UI_TEXT.common.save}
                                                variant="primary"
                                                onClick={handleConfirm}
                                        />
                                </div>
                        </div>
                </div>
        );
};

export default CloseRequestModal;
