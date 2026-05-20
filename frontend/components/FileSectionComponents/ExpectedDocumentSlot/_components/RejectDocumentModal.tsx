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

const PREDEFINED_REASONS = [
  "Scan ilizibil",
  "Document expirat",
  "Semnătură lipsă",
  "Document incomplet",
  "Format incorect",
  "Document deteriorat",
  "Altceva",
];

const RejectDocumentModal: React.FC<RejectDocumentModalProps> = ({
  isOpen,
  onClose,
  onConfirm,
  documentTitle,
}) => {
  const [selected, setSelected] = useState<string | null>(null);
  const [customReason, setCustomReason] = useState("");
  const [error, setError] = useState("");

  if (!isOpen) return null;

  const isCustom = selected === "Altceva";

  const getFinalReason = (): string => {
    if (!selected) return "";
    if (isCustom) return customReason.trim();
    return selected;
  };

  const handleConfirm = () => {
    const reason = getFinalReason();
    if (!selected) {
      setError("Selectați un motiv de respingere.");
      return;
    }
    if (isCustom && !customReason.trim()) {
      setError("Introduceți motivul respingerii.");
      return;
    }
    onConfirm(reason);
    handleClose();
  };

  const handleClose = () => {
    setSelected(null);
    setCustomReason("");
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
          <button className="reject-document-modal-close" onClick={handleClose}>
            ×
          </button>
        </div>

        <div className="reject-document-modal-body">
          <p className="reject-document-warning">
            {UI_TEXT.modals.rejectDocument.subtitle1(documentTitle)}
          </p>

          <div className="reject-reasons-grid">
            {PREDEFINED_REASONS.map((reason) => (
              <button
                key={reason}
                className={`reject-reason-chip ${selected === reason ? "reject-reason-chip--active" : ""}`}
                onClick={() => {
                  setSelected(reason);
                  setError("");
                }}
              >
                {reason}
              </button>
            ))}
          </div>

          {isCustom && (
            <div className="reject-custom-input">
              <TextArea
                label="Motivul respingerii"
                value={customReason}
                onChange={(e) => {
                  setCustomReason(e.target.value);
                  setError("");
                }}
                placeholder="Descrieți motivul respingerii..."
                minHeight={80}
                maxHeight={150}
                fullWidth
              />
            </div>
          )}

          {error && <span className="reject-error">{error}</span>}
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