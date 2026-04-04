"use client";
import React from "react";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import "./Modal.css";

interface ModalProps {
        isOpen: boolean;
        onClose: () => void;
        onConfirm?: () => void;
        title: string;
        children: React.ReactNode;
        closeOnConfirm?: boolean;
        hideFooter?: boolean;
}

const Modal: React.FC<ModalProps> = ({
        isOpen,
        onClose,
        onConfirm,
        title,
        children,
        closeOnConfirm = true,
        hideFooter = false,
}) => {
        if (!isOpen) return null;

        const handleConfirm = () => {
                onConfirm?.();
                if (closeOnConfirm) onClose();
        };

        return (
                <div className="modal-overlay" onClick={onClose}>
                        <div className="modal-content" onClick={(e) => e.stopPropagation()}>
                                <div className="modal-header">
                                        <h3>{title}</h3>
                                        <button className="modal-close" onClick={onClose}>
                                                ×
                                        </button>
                                </div>
                                <div className="modal-body">{children}</div>
                                {!hideFooter && (
                                        <div className="modal-footer">
                                                <ButtonPrimary
                                                        text="Anulează"
                                                        variant="ghost"
                                                        onClick={onClose}
                                                />
                                                <ButtonPrimary
                                                        text="Continuă"
                                                        variant="primary"
                                                        onClick={handleConfirm}
                                                />
                                        </div>
                                )}
                        </div>
                </div>
        );
};

export default Modal;
