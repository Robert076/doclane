"use client";
import { useState } from "react";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import "./InvitationCodeGenerator.css";
import { MdContentCopy, MdCheck, MdClose } from "react-icons/md";
import { UI_TEXT } from "@/locales/ro";

const InvitationCodeGenerator = () => {
        const [isModalOpen, setIsModalOpen] = useState(false);
        const [code, setCode] = useState<string | null>(null);
        const [isLoading, setIsLoading] = useState(false);
        const [isCopied, setIsCopied] = useState(false);
        const [error, setError] = useState<string | null>(null);

        const generateCode = async () => {
                setIsLoading(true);
                setError(null);

                try {
                        const response = await fetch("/api/backend/invitations/generate", {
                                method: "POST",
                                credentials: "include",
                                headers: {
                                        "Content-Type": "application/json",
                                },
                                body: JSON.stringify({
                                        expires_in_days: 7,
                                }),
                        });

                        if (!response.ok) {
                                throw new Error(UI_TEXT.modals.generateCode.errorMaxCodes);
                        }

                        const data = await response.json();
                        setCode(data.data.code);
                } catch (err) {
                        setError(UI_TEXT.modals.generateCode.errorMaxCodes);
                        console.error(err);
                } finally {
                        setIsLoading(false);
                }
        };

        const copyToClipboard = async () => {
                if (!code) return;

                try {
                        await navigator.clipboard.writeText(code);
                        setIsCopied(true);
                        setTimeout(() => setIsCopied(false), 2000);
                } catch (err) {
                        console.error("Failed to copy:", err);
                }
        };

        const closeModal = () => {
                setIsModalOpen(false);
                setCode(null);
                setIsCopied(false);
                setError(null);
        };

        return (
                <>
                        <ButtonPrimary
                                text={UI_TEXT.buttons.addClient.normal}
                                onClick={() => setIsModalOpen(true)}
                        />

                        {isModalOpen && (
                                <div className="modal-overlay" onClick={closeModal}>
                                        <div
                                                className="modal-content"
                                                onClick={(e) => e.stopPropagation()}
                                        >
                                                <button
                                                        className="modal-close"
                                                        onClick={closeModal}
                                                        aria-label="Close"
                                                >
                                                        <MdClose size={20} />
                                                </button>

                                                {!code ? (
                                                        <>
                                                                <h2 className="modal-title">
                                                                        {
                                                                                UI_TEXT.modals
                                                                                        .generateCode
                                                                                        .title
                                                                        }
                                                                </h2>
                                                                <p className="modal-description">
                                                                        {
                                                                                UI_TEXT.modals
                                                                                        .generateCode
                                                                                        .subtitle1
                                                                        }
                                                                </p>
                                                                <p className="modal-note">
                                                                        {
                                                                                UI_TEXT.modals
                                                                                        .generateCode
                                                                                        .subtitle2
                                                                        }
                                                                </p>
                                                                <div className="modal-actions">
                                                                        <ButtonPrimary
                                                                                text={
                                                                                        isLoading
                                                                                                ? UI_TEXT
                                                                                                          .common
                                                                                                          .loading
                                                                                                : UI_TEXT
                                                                                                          .common
                                                                                                          .continue
                                                                                }
                                                                                onClick={
                                                                                        generateCode
                                                                                }
                                                                                disabled={
                                                                                        isLoading
                                                                                }
                                                                        />
                                                                        <ButtonPrimary
                                                                                text={
                                                                                        UI_TEXT
                                                                                                .common
                                                                                                .cancel
                                                                                }
                                                                                variant="ghost"
                                                                                onClick={
                                                                                        closeModal
                                                                                }
                                                                        />
                                                                </div>
                                                        </>
                                                ) : (
                                                        <>
                                                                <h2 className="modal-title">
                                                                        {
                                                                                UI_TEXT.modals
                                                                                        .generateCode
                                                                                        .title
                                                                        }
                                                                </h2>
                                                                <p className="modal-description">
                                                                        {
                                                                                UI_TEXT.modals
                                                                                        .generateCode
                                                                                        .subtitle3
                                                                        }
                                                                </p>
                                                                <div className="code-box">
                                                                        <span className="code-text">
                                                                                {code}
                                                                        </span>
                                                                        <button
                                                                                className="copy-button"
                                                                                onClick={
                                                                                        copyToClipboard
                                                                                }
                                                                                aria-label="Copy code"
                                                                        >
                                                                                {isCopied ? (
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
                                                                </div>
                                                                <p className="code-expiry">
                                                                        {
                                                                                UI_TEXT.modals
                                                                                        .generateCode
                                                                                        .expiryNotice
                                                                        }
                                                                </p>
                                                                <div className="modal-actions">
                                                                        <ButtonPrimary
                                                                                text={
                                                                                        UI_TEXT
                                                                                                .common
                                                                                                .continue
                                                                                }
                                                                                onClick={
                                                                                        closeModal
                                                                                }
                                                                        />
                                                                </div>
                                                        </>
                                                )}

                                                {error && (
                                                        <p className="error-message">
                                                                {error}
                                                        </p>
                                                )}
                                        </div>
                                </div>
                        )}
                </>
        );
};

export default InvitationCodeGenerator;
