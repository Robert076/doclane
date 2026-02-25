"use client";
import React from "react";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import "./DeactivateClientModal.css";
import { UI_TEXT } from "@/locales/ro";

interface DeactivateClientModalProps {
        isOpen: boolean;
        onClose: () => void;
        onConfirm: () => void;
        clientName: string;
}

const DeactivateClientModal: React.FC<DeactivateClientModalProps> = ({
        isOpen,
        onClose,
        onConfirm,
        clientName,
}) => {
        if (!isOpen) return null;

        const handleConfirm = () => {
                onConfirm();
                onClose();
        };

        return (
                <div className="deactivate-modal-overlay" onClick={onClose}>
                        <div
                                className="deactivate-modal-content"
                                onClick={(e) => e.stopPropagation()}
                        >
                                <div className="deactivate-modal-header">
                                        <h3>{UI_TEXT.modals.deactivateClient.title}</h3>
                                        <button
                                                className="deactivate-modal-close"
                                                onClick={onClose}
                                        >
                                                ×
                                        </button>
                                </div>
                                <div className="deactivate-modal-body">
                                        <p className="deactivate-warning">
                                                {UI_TEXT.modals.deactivateClient.subtitle1(
                                                        clientName,
                                                )}
                                        </p>
                                        <p className="deactivate-description">
                                                {UI_TEXT.modals.deactivateClient.subtitle2}
                                        </p>
                                </div>
                                <div className="deactivate-modal-footer">
                                        <ButtonPrimary
                                                text={UI_TEXT.common.cancel}
                                                variant="ghost"
                                                onClick={onClose}
                                        />
                                        <ButtonPrimary
                                                text={UI_TEXT.modals.deactivateClient.confirm}
                                                variant="primary"
                                                onClick={handleConfirm}
                                        />
                                </div>
                        </div>
                </div>
        );
};

export default DeactivateClientModal;
