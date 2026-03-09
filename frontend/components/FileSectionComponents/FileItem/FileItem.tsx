"use client";

import { useState } from "react";
import toast from "react-hot-toast";

// Tipurile și Utilitarele
import { DocumentFile } from "@/types";
import { formatDate } from "@/lib/client/formatDate";
import { presignDocumentURL } from "@/lib/api/requests";
import "./FileItem.css";

interface FileItemProps {
        file: DocumentFile;
}

// 1. Optimizare: Funcția scoasă în afara componentei pentru a nu fi recreată la fiecare randare
const formatFileSize = (bytes: number) => {
        if (!bytes || bytes === 0) return "0 Bytes";
        const k = 1024;
        const sizes = ["Bytes", "KB", "MB", "GB"];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
};

export default function FileItem({ file }: FileItemProps) {
        const [isRequesting, setIsRequesting] = useState(false);

        // 2. Request standardizat cu Toast de tip loading
        const handleViewFile = async () => {
                if (isRequesting) return;

                setIsRequesting(true);
                const loadingToast = toast.loading("Se deschide fișierul...");

                try {
                        const res = await presignDocumentURL(
                                file.document_request_id,
                                file.id,
                        );

                        if (!res.success || !res.data) {
                                throw new Error(
                                        res.error ||
                                                res.message ||
                                                "Nu s-a putut genera link-ul fișierului.",
                                );
                        }

                        // Închidem toast-ul de loading și deschidem tab-ul
                        toast.dismiss(loadingToast);
                        window.open(res.data.url, "_blank", "noopener,noreferrer");
                } catch (error: any) {
                        toast.error(error.message, { id: loadingToast });
                } finally {
                        setIsRequesting(false);
                }
        };

        return (
                <div className="file-item">
                        <div className="file-icon-wrapper">
                                <FileIcon mimeType={file.mime_type} />
                        </div>

                        <div className="file-info-wrapper">
                                <p className="file-name" title={file.file_name}>
                                        {file.file_name}
                                </p>
                                <div className="file-metadata">
                                        <span>{formatFileSize(file.file_size)}</span>
                                        <span className="metadata-separator">•</span>
                                        <span>{formatDate(file.uploaded_at)}</span>

                                        {file.uploaded_by && (
                                                <>
                                                        <span className="metadata-separator">
                                                                •
                                                        </span>
                                                        <span>
                                                                Încărcat de{" "}
                                                                {file.uploaded_by_first_name}{" "}
                                                                {file.uploaded_by_last_name}
                                                        </span>
                                                </>
                                        )}
                                </div>
                        </div>

                        <div className="file-actions">
                                <button
                                        onClick={handleViewFile}
                                        className="view-button"
                                        disabled={isRequesting}
                                >
                                        {isRequesting ? "Se deschide..." : "Vezi fișierul"}
                                </button>
                        </div>
                </div>
        );
}

// Sub-componenta lăsată la fel, își face treaba foarte bine
function FileIcon({ mimeType }: { mimeType: string }) {
        const isImage = mimeType?.toLowerCase().includes("image");
        const isPDF = mimeType?.toLowerCase().includes("pdf");

        const badgeClass = isPDF ? "pdf" : isImage ? "img" : "doc";
        const label = isPDF ? "PDF" : isImage ? "IMG" : "DOC";

        return <div className={`file-type-badge ${badgeClass}`}>{label}</div>;
}
