"use client";
import React from "react";
import ButtonPrimary from "@/components/Buttons/ButtonPrimary/ButtonPrimary";
import "./DeactivateClientModal.css";

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
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h3>Deactivate Client Account</h3>
          <button className="modal-close" onClick={onClose}>
            ×
          </button>
        </div>
        <div className="modal-body">
          <p className="deactivate-warning">
            Are you sure you want to deactivate <strong>{clientName}</strong>'s account?
          </p>
          <p className="deactivate-description">
            This action will prevent them from accessing their account. This can be reversed
            later if needed.
          </p>
        </div>
        <div className="modal-footer">
          <ButtonPrimary text="Cancel" variant="ghost" onClick={onClose} />
          <ButtonPrimary text="Deactivate Account" variant="primary" onClick={handleConfirm} />
        </div>
      </div>
    </div>
  );
};

export default DeactivateClientModal;
