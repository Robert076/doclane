"use client";
import { useState } from "react";
import "./FileItem.css";
import { DocumentFile } from "@/types";
import { presignDocumentURL } from "@/lib/api/api";
import toast from "react-hot-toast";
import { formatDate } from "@/lib/client/formatDate";

interface FileItemProps {
        file: DocumentFile;
}

export default function FileItem({ file }: FileItemProps) {
        const [isRequesting, setIsRequesting] = useState(false);

        const formatFileSize = (bytes: number) => {
                if (!bytes || bytes === 0) return "0 Bytes";
                const k = 1024;
                const sizes = ["Bytes", "KB", "MB", "GB"];
                const i = Math.floor(Math.log(bytes) / Math.log(k));
                return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
        };

        const handleViewFile = async () => {
                if (isRequesting) return;

                setIsRequesting(true);
                try {
                        const result = await presignDocumentURL(
                                file.document_request_id,
                                file.id,
                        );

                        if (result.success && result.data?.url) {
                                window.open(result.data.url, "_blank", "noopener,noreferrer");
                                toast.success(result.message);
                        } else {
                                toast.error(result.message);
                                throw new Error(result.message || "Error generating the link");
                        }
                } catch (error) {
                        console.error("Error fetching file URL:", error);
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
                                                                {"Uploaded by " +
                                                                        file.uploaded_by_first_name +
                                                                        " " +
                                                                        file.uploaded_by_last_name}
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
                                        {isRequesting ? "Loading..." : "View file"}
                                </button>
                        </div>
                </div>
        );
}

function FileIcon({ mimeType }: { mimeType: string }) {
        const isImage = mimeType?.toLowerCase().includes("image");
        const isPDF = mimeType?.toLowerCase().includes("pdf");

        const badgeClass = isPDF ? "pdf" : isImage ? "img" : "doc";
        const label = isPDF ? "PDF" : isImage ? "IMG" : "DOC";

        return <div className={`file-type-badge ${badgeClass}`}>{label}</div>;
}
