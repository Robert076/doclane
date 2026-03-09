"use client";

import { useState } from "react";
import toast from "react-hot-toast";
import { MdFileOpen } from "react-icons/md";

// Tipurile și Contextul
import { DocumentFile, ExpectedDocument, RequestStatus } from "@/types";
import { useUser } from "@/context/UserContext";
import { presignExampleURL } from "@/lib/api/requests";
import { usePagination } from "@/hooks/usePagination";

// Sub-componente
import FileItem from "../FileItem/FileItem";
import PaginationFooter from "@/components/FileSectionComponents/FileSection/_components/PaginationFooter";
import StatusBadge from "@/components/Pages/RequestsComponents/StatusBadge";
import DocumentSlotActions from "./_components/DocumentSlotActions";
import { useDocumentStatus } from "./_hooks/useDocumentStatus";

import "./ExpectedDocumentSlot.css";

interface ExpectedDocumentSlotProps {
        expectedDocument: ExpectedDocument;
        requestId: string;
        uploadedFiles: DocumentFile[];
}

const ITEMS_PER_PAGE = 3;
const DECISION_STATUSES = ["approved", "rejected"];

export default function ExpectedDocumentSlot({
        expectedDocument,
        requestId,
        uploadedFiles,
}: ExpectedDocumentSlotProps) {
        const [isLoadingExample, setIsLoadingExample] = useState(false);
        const user = useUser();
        const isProfessional = user?.role === "PROFESSIONAL";

        // 1. Folosim hook-ul nostru standard pentru paginare
        const { currentPage, setCurrentPage, totalPages, paginatedItems } = usePagination(
                uploadedFiles,
                ITEMS_PER_PAGE,
        );

        const { approve, reject, reset, isLoading } = useDocumentStatus(
                expectedDocument.id.toString(),
                requestId,
                uploadedFiles.length > 0,
        );

        const hasDecision = DECISION_STATUSES.includes(expectedDocument.status);

        // 2. Standardizăm request-ul cu Try/Catch și Toast Loading
        const handleViewExample = async () => {
                setIsLoadingExample(true);
                const loadingToast = toast.loading("Se deschide exemplul...");

                try {
                        const res = await presignExampleURL(expectedDocument.id);

                        if (!res.success || !res.data) {
                                throw new Error(
                                        res.error ||
                                                res.message ||
                                                "Nu s-a putut încărca exemplul.",
                                );
                        }

                        // Ștergem toast-ul de loading pentru că se deschide un tab nou
                        toast.dismiss(loadingToast);
                        window.open(res.data.url, "_blank");
                } catch (error: any) {
                        toast.error(error.message, { id: loadingToast });
                } finally {
                        setIsLoadingExample(false);
                }
        };

        return (
                <div className="expected-document-slot">
                        <div className="expected-document-slot-top">
                                <div className="expected-document-slot-status-row">
                                        <StatusBadge
                                                status={
                                                        expectedDocument.status as RequestStatus
                                                }
                                        />
                                </div>

                                <div className="expected-document-slot-title-row">
                                        <div className="expected-document-slot-info">
                                                <span className="expected-document-slot-title">
                                                        {expectedDocument.title}
                                                </span>

                                                {expectedDocument.description && (
                                                        <span className="expected-document-slot-description">
                                                                {expectedDocument.description}
                                                        </span>
                                                )}

                                                {/* Am grupat condițiile pentru lizibilitate */}
                                                {expectedDocument.status === "rejected" &&
                                                        expectedDocument.rejection_reason && (
                                                                <span className="expected-document-slot-rejection-reason">
                                                                        <strong>
                                                                                Motiv refuz:
                                                                        </strong>{" "}
                                                                        {
                                                                                expectedDocument.rejection_reason
                                                                        }
                                                                </span>
                                                        )}

                                                {expectedDocument.example_file_path && (
                                                        <button
                                                                type="button"
                                                                className="expected-document-slot-example-btn"
                                                                onClick={handleViewExample}
                                                                disabled={isLoadingExample}
                                                        >
                                                                <MdFileOpen />
                                                                {isLoadingExample
                                                                        ? "Se încarcă..."
                                                                        : "Vezi exemplu"}
                                                        </button>
                                                )}
                                        </div>

                                        <div
                                                className={
                                                        hasDecision
                                                                ? "upload-button-wrapper"
                                                                : "upload-button-wrapper--actions"
                                                }
                                                style={{ display: "flex", gap: "8px" }}
                                        >
                                                <DocumentSlotActions
                                                        isProfessional={isProfessional}
                                                        status={expectedDocument.status}
                                                        hasFiles={uploadedFiles.length > 0}
                                                        isLoading={isLoading}
                                                        requestId={parseInt(requestId, 10)}
                                                        expectedDocumentId={
                                                                expectedDocument.id
                                                        }
                                                        documentTitle={expectedDocument.title}
                                                        onApprove={approve}
                                                        onReject={reject}
                                                        onReset={reset}
                                                />
                                        </div>
                                </div>
                        </div>

                        {uploadedFiles.length > 0 && (
                                <div className="expected-document-slot-files">
                                        {/* 3. Randăm itemii direct din paginarea automată */}
                                        {paginatedItems.map((file) => (
                                                <FileItem key={file.id} file={file} />
                                        ))}

                                        {totalPages > 1 && (
                                                <PaginationFooter
                                                        currentPage={currentPage}
                                                        totalPages={totalPages}
                                                        setCurrentPage={setCurrentPage}
                                                />
                                        )}
                                </div>
                        )}
                </div>
        );
}
