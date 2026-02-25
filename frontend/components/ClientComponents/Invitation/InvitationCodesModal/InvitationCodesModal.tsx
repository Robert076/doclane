"use client";
import { useState } from "react";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import NotFound from "@/components/OtherComponents/NotFound/NotFound";
import "./InvitationCodesModal.css";
import { MdClose, MdContentCopy, MdCheck, MdDelete } from "react-icons/md";
import { UI_TEXT } from "@/locales/ro";

interface InvitationCode {
        id: number;
        code: string;
        used_by_user_id: number | null;
        expires_at: string | null;
        created_at: string;
}

const InvitationCodesModal = () => {
        const [isModalOpen, setIsModalOpen] = useState(false);
        const [codes, setCodes] = useState<InvitationCode[]>([]);
        const [isLoading, setIsLoading] = useState(false);
        const [copiedId, setCopiedId] = useState<number | null>(null);
        const [error, setError] = useState<string | null>(null);

        const fetchCodes = async () => {
                setIsLoading(true);
                setError(null);

                try {
                        const response = await fetch("/api/backend/invitations/my-codes", {
                                method: "GET",
                                credentials: "include",
                        });

                        if (!response.ok) {
                                throw new Error("Failed to fetch codes");
                        }

                        const data = await response.json();
                        setCodes(data.data || []);
                } catch (err) {
                        setError("Failed to load invitation codes.");
                        console.error(err);
                } finally {
                        setIsLoading(false);
                }
        };

        const deleteCode = async (id: number) => {
                try {
                        const response = await fetch(`/api/backend/invitations/${id}`, {
                                method: "DELETE",
                                credentials: "include",
                        });

                        if (!response.ok) {
                                throw new Error("Failed to delete code");
                        }

                        setCodes(codes.filter((code) => code.id !== id));
                } catch (err) {
                        setError("Failed to delete code.");
                        console.error(err);
                }
        };

        const copyToClipboard = async (code: string, id: number) => {
                try {
                        await navigator.clipboard.writeText(code);
                        setCopiedId(id);
                        setTimeout(() => setCopiedId(null), 2000);
                } catch (err) {
                        console.error("Failed to copy:", err);
                }
        };

        const openModal = () => {
                setIsModalOpen(true);
                fetchCodes();
        };

        const closeModal = () => {
                setIsModalOpen(false);
                setError(null);
        };

        const formatDate = (dateString: string) => {
                const date = new Date(dateString);
                return date.toLocaleDateString("en-US", {
                        month: "short",
                        day: "numeric",
                        year: "numeric",
                });
        };

        const isExpired = (expiresAt: string | null) => {
                if (!expiresAt) return false;
                return new Date(expiresAt) < new Date();
        };

        const unusedCodes = codes.filter((code) => !code.used_by_user_id);

        return (
                <>
                        <ButtonPrimary
                                text={UI_TEXT.buttons.viewInvitationCodes.normal}
                                variant="primary"
                                onClick={openModal}
                        />

                        {isModalOpen && (
                                <div className="modal-overlay" onClick={closeModal}>
                                        <div
                                                className="modal-content codes-modal"
                                                onClick={(e) => e.stopPropagation()}
                                        >
                                                <button
                                                        className="modal-close"
                                                        onClick={closeModal}
                                                        aria-label="Close"
                                                >
                                                        <MdClose size={24} />
                                                </button>

                                                <h2 className="modal-title">
                                                        {UI_TEXT.modals.codesModal.title}
                                                </h2>
                                                <p className="modal-description">
                                                        {UI_TEXT.modals.codesModal.subtitle}
                                                </p>

                                                {isLoading ? (
                                                        <div className="codes-loading">
                                                                <p>Loading codes...</p>
                                                        </div>
                                                ) : unusedCodes.length === 0 ? (
                                                        <NotFound
                                                                text="No active invitation codes."
                                                                subtext="Generate one to invite clients."
                                                        />
                                                ) : (
                                                        <div className="codes-list">
                                                                {unusedCodes.map((code) => (
                                                                        <div
                                                                                key={code.id}
                                                                                className={`code-item ${isExpired(code.expires_at) ? "expired" : ""}`}
                                                                        >
                                                                                <div className="code-info">
                                                                                        <div className="code-value-section">
                                                                                                <span className="code-value">
                                                                                                        {
                                                                                                                code.code
                                                                                                        }
                                                                                                </span>
                                                                                                {isExpired(
                                                                                                        code.expires_at,
                                                                                                ) && (
                                                                                                        <span className="code-status expired-badge">
                                                                                                                Expired
                                                                                                        </span>
                                                                                                )}
                                                                                        </div>
                                                                                        <div className="code-meta">
                                                                                                <span className="code-date">
                                                                                                        {
                                                                                                                UI_TEXT
                                                                                                                        .modals
                                                                                                                        .codesModal
                                                                                                                        .createdAt
                                                                                                        }
                                                                                                        {formatDate(
                                                                                                                code.created_at,
                                                                                                        )}
                                                                                                </span>
                                                                                                {code.expires_at && (
                                                                                                        <span className="code-expiry">
                                                                                                                {
                                                                                                                        UI_TEXT
                                                                                                                                .modals
                                                                                                                                .codesModal
                                                                                                                                .expiresAt
                                                                                                                }
                                                                                                                {formatDate(
                                                                                                                        code.expires_at,
                                                                                                                )}
                                                                                                        </span>
                                                                                                )}
                                                                                        </div>
                                                                                </div>
                                                                                <div className="code-actions">
                                                                                        <button
                                                                                                className="icon-button copy-btn"
                                                                                                onClick={() =>
                                                                                                        copyToClipboard(
                                                                                                                code.code,
                                                                                                                code.id,
                                                                                                        )
                                                                                                }
                                                                                                aria-label="Copy code"
                                                                                                disabled={isExpired(
                                                                                                        code.expires_at,
                                                                                                )}
                                                                                        >
                                                                                                {copiedId ===
                                                                                                code.id ? (
                                                                                                        <MdCheck
                                                                                                                size={
                                                                                                                        20
                                                                                                                }
                                                                                                        />
                                                                                                ) : (
                                                                                                        <MdContentCopy
                                                                                                                size={
                                                                                                                        20
                                                                                                                }
                                                                                                        />
                                                                                                )}
                                                                                        </button>
                                                                                        <button
                                                                                                className="icon-button delete-btn"
                                                                                                onClick={() =>
                                                                                                        deleteCode(
                                                                                                                code.id,
                                                                                                        )
                                                                                                }
                                                                                                aria-label="Delete code"
                                                                                        >
                                                                                                <MdDelete
                                                                                                        size={
                                                                                                                20
                                                                                                        }
                                                                                                />
                                                                                        </button>
                                                                                </div>
                                                                        </div>
                                                                ))}
                                                        </div>
                                                )}

                                                {error && (
                                                        <p className="error-message">
                                                                {error}
                                                        </p>
                                                )}

                                                <div className="modal-actions">
                                                        <ButtonPrimary
                                                                text={UI_TEXT.common.close}
                                                                onClick={closeModal}
                                                        />
                                                </div>
                                        </div>
                                </div>
                        )}
                </>
        );
};

export default InvitationCodesModal;
