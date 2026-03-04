"use client";
import React, { useState } from "react";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import "./RejectDocumentModal.css";
import { UI_TEXT } from "@/locales/ro";
import TextArea from "@/components/InputComponents/TextArea";

interface RejectDocumentModalProps {
        isOpen: boolean;
        onClose: () => void;
        onConfirm: (reason: string) => void;
        documentTitle: string;
}

const RejectDocumentModal: React.FC<RejectDocumentModalProps> = ({
        isOpen,
        onClose,
        onConfirm,
        documentTitle,
}) => {
        const [reason, setReason] = useState("");
        const [error, setError] = useState("");

        if (!isOpen) return null;

        const handleConfirm = () => {
                if (!reason.trim()) {
                        setError(UI_TEXT.modals.rejectDocument.reasonRequired);
                        return;
                }
                onConfirm(reason);
                setReason("");
                setError("");
                onClose();
        };

        const handleClose = () => {
                setReason("");
                setError("");
                onClose();
        };

        return (
                <div className="reject-document-modal-overlay" onClick={handleClose}>
                        <div
                                className="reject-document-modal-content"
                                onClick={(e) => e.stopPropagation()}
                        >
                                <div className="reject-document-modal-header">
                                        <h3>{UI_TEXT.modals.rejectDocument.title}</h3>
                                        <button
                                                className="reject-document-modal-close"
                                                onClick={handleClose}
                                        >
                                                ×
                                        </button>
                                </div>

                                <div className="reject-document-modal-body">
                                        <p className="reject-document-warning">
                                                {UI_TEXT.modals.rejectDocument.subtitle1(
                                                        documentTitle,
                                                )}
                                        </p>
                                        {/* <p className="reject-document-description">
                                                {UI_TEXT.modals.rejectDocument.subtitle2}
                                        </p> */}

                                        <div className="reject-document-input-group">
                                                <TextArea
                                                        label={
                                                                UI_TEXT.modals.rejectDocument
                                                                        .reasonLabel
                                                        }
                                                        value={reason}
                                                        onChange={(e) => {
                                                                setReason(e.target.value);
                                                                setError("");
                                                        }}
                                                        placeholder={
                                                                UI_TEXT.modals.rejectDocument
                                                                        .reasonPlaceholder
                                                        }
                                                        minHeight={80}
                                                        maxHeight={150}
                                                        fullWidth
                                                />
                                                {error && (
                                                        <span className="reject-error">
                                                                {error}
                                                        </span>
                                                )}
                                        </div>
                                </div>

                                <div className="reject-document-modal-footer">
                                        <ButtonPrimary
                                                text={UI_TEXT.modals.rejectDocument.confirm}
                                                variant="primary"
                                                onClick={handleConfirm}
                                        />
                                        <ButtonPrimary
                                                text={UI_TEXT.common.cancel}
                                                variant="ghost"
                                                onClick={handleClose}
                                        />
                                </div>
                        </div>
                </div>
        );
};

export default RejectDocumentModal;
