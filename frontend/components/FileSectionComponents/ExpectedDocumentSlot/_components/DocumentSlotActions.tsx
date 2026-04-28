"use client";
import React, { useRef, useState, useEffect } from "react";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import UploadDocumentButton from "@/components/Pages/RequestsComponents/UploadDocumentButton";
import RejectDocumentModal from "./RejectDocumentModal";
import { ExpectedDocumentStatus } from "@/types";
import { extractFileText, interpretFileText, speakFileText } from "@/lib/api/requests";
import toast from "react-hot-toast";
import "./DocumentSlotActions.css";

interface DocumentSlotActionsProps {
        canManage: boolean;
        status: ExpectedDocumentStatus;
        hasFiles: boolean;
        isLoading: boolean;
        requestId: number;
        expectedDocumentId: number;
        documentTitle: string;
        latestFileId: number | null;
        onApprove: () => void;
        onReject: (reason: string) => void;
        onReset: () => void;
        onExtractedText: (text: string) => void;
        onInterpretedText: (text: string) => void;
}

export default function DocumentSlotActions({
        canManage,
        status,
        hasFiles,
        isLoading,
        requestId,
        expectedDocumentId,
        documentTitle,
        latestFileId,
        onApprove,
        onReject,
        onReset,
        onExtractedText,
        onInterpretedText,
}: DocumentSlotActionsProps) {
        const [isRejectModalOpen, setIsRejectModalOpen] = useState(false);
        const [isExtracting, setIsExtracting] = useState(false);
        const [isInterpreting, setIsInterpreting] = useState(false);
        const [isMenuOpen, setIsMenuOpen] = useState(false);
        const [isPlaying, setIsPlaying] = useState(false);
        const menuRef = useRef<HTMLDivElement>(null);
        const audioRef = useRef<HTMLAudioElement | null>(null);

        useEffect(() => {
                const handleClickOutside = (e: MouseEvent) => {
                        if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
                                setIsMenuOpen(false);
                        }
                };
                document.addEventListener("mousedown", handleClickOutside);
                return () => document.removeEventListener("mousedown", handleClickOutside);
        }, []);

        const handleExtract = async () => {
                if (!latestFileId) return;
                setIsMenuOpen(false);
                setIsExtracting(true);
                const loadingToast = toast.loading("Se extrage textul...");
                const res = await extractFileText(requestId, latestFileId);
                setIsExtracting(false);
                if (!res.success || !res.data) {
                        toast.error(res.message ?? "Eroare la extragere.", {
                                id: loadingToast,
                        });
                        return;
                }
                toast.dismiss(loadingToast);
                onExtractedText(res.data.text);
        };

        const handleInterpret = async () => {
                if (!latestFileId) return;
                setIsMenuOpen(false);
                setIsInterpreting(true);
                const loadingToast = toast.loading("Se interpretează documentul cu AI...");
                const res = await interpretFileText(requestId, latestFileId, documentTitle);
                setIsInterpreting(false);
                if (!res.success || !res.data) {
                        toast.error(res.message ?? "Eroare la interpretare.", {
                                id: loadingToast,
                        });
                        return;
                }
                toast.dismiss(loadingToast);
                onInterpretedText(res.data.interpretation);
        };

        const handleSpeak = async () => {
                if (!latestFileId) return;
                setIsMenuOpen(false);
                if (audioRef.current) {
                        audioRef.current.pause();
                        audioRef.current = null;
                        setIsPlaying(false);
                        return;
                }
                const loadingToast = toast.loading("Se generează audio...");
                setIsPlaying(true);
                const res = await speakFileText(requestId, latestFileId);
                if (!res.success || !res.data) {
                        toast.error(res.message ?? "Eroare la generarea audio.", {
                                id: loadingToast,
                        });
                        setIsPlaying(false);
                        return;
                }
                toast.dismiss(loadingToast);
                const audio = new Audio(`data:audio/mpeg;base64,${res.data.audio}`);
                audioRef.current = audio;
                audio.play();
                audio.onended = () => {
                        setIsPlaying(false);
                        audioRef.current = null;
                };
        };

        if (!canManage) {
                if (status === "accepted") return null;
                return (
                        <UploadDocumentButton
                                requestId={requestId}
                                expectedDocumentId={expectedDocumentId}
                        />
                );
        }

        if (status === "accepted" || status === "rejected") {
                return (
                        <ButtonPrimary
                                text={isLoading ? "..." : "Anulează"}
                                variant="ghost"
                                onClick={onReset}
                        />
                );
        }

        if (!hasFiles) {
                return null;
        }
        const isBusy = isExtracting || isInterpreting || isLoading;

        return (
                <>
                        <div className="slot-actions-wrapper">
                                {/* three dot menu for AI tools */}
                                <div className="slot-actions" ref={menuRef}>
                                        <button
                                                className="slot-actions-trigger"
                                                disabled={isBusy}
                                                onClick={() => setIsMenuOpen((prev) => !prev)}
                                                title="Mai multe opțiuni"
                                        >
                                                {isBusy ? (
                                                        <span className="slot-actions-spinner" />
                                                ) : (
                                                        "⋯"
                                                )}
                                        </button>
                                        {isMenuOpen && (
                                                <div className="slot-actions-menu">
                                                        <button
                                                                className="slot-actions-item"
                                                                onClick={handleExtract}
                                                        >
                                                                Extrage text
                                                        </button>
                                                        <button
                                                                className="slot-actions-item"
                                                                onClick={handleInterpret}
                                                        >
                                                                Interpretează cu AI
                                                        </button>
                                                        <button
                                                                className="slot-actions-item"
                                                                onClick={handleSpeak}
                                                        >
                                                                {isPlaying
                                                                        ? "Oprește audio"
                                                                        : "Citește document"}
                                                        </button>
                                                </div>
                                        )}
                                </div>

                                {/* primary actions always visible */}
                                <ButtonPrimary
                                        text={isLoading ? "..." : "Refuză"}
                                        variant="ghost"
                                        disabled={isLoading}
                                        fullWidth
                                        onClick={() => setIsRejectModalOpen(true)}
                                />
                                <ButtonPrimary
                                        text={isLoading ? "..." : "Aprobă"}
                                        variant="primary"
                                        disabled={isLoading}
                                        fullWidth
                                        onClick={onApprove}
                                />
                        </div>

                        <RejectDocumentModal
                                isOpen={isRejectModalOpen}
                                onClose={() => setIsRejectModalOpen(false)}
                                onConfirm={onReject}
                                documentTitle={documentTitle}
                        />
                </>
        );
}
