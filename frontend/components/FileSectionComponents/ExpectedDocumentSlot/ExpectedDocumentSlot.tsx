"use client";
import { useState } from "react";
import toast from "react-hot-toast";
import { MdFileOpen } from "react-icons/md";
import { DocumentFile, ExpectedDocument, RequestStatus } from "@/types";
import { useUser } from "@/context/UserContext";
import { presignExampleURL } from "@/lib/api/requests";
import { usePagination } from "@/hooks/usePagination";
import FileItem from "../FileItem/FileItem";
import PaginationFooter from "@/components/FileSectionComponents/FileSection/_components/PaginationFooter";
import StatusBadge from "@/components/Pages/RequestsComponents/StatusBadge";
import DocumentSlotActions from "./_components/DocumentSlotActions";
import { useDocumentStatus } from "./_hooks/useDocumentStatus";
import "./ExpectedDocumentSlot.css";

interface ExpectedDocumentSlotProps {
        expectedDocument: ExpectedDocument;
        requestId: number;
        uploadedFiles: DocumentFile[];
}

const ITEMS_PER_PAGE = 3;
const DECISION_STATUSES = ["accepted", "rejected"];

export default function ExpectedDocumentSlot({
        expectedDocument,
        requestId,
        uploadedFiles,
}: ExpectedDocumentSlotProps) {
        const [isLoadingExample, setIsLoadingExample] = useState(false);
        const [extractedText, setExtractedText] = useState<string | null>(null);
        const [interpretedText, setInterpretedText] = useState<string | null>(null);

        const user = useUser();
        const canManage = user.role === "admin" || user.department_id !== null;

        const { currentPage, setCurrentPage, totalPages, paginatedItems } = usePagination(
                uploadedFiles,
                ITEMS_PER_PAGE,
        );

        const { approve, reject, reset, isLoading } = useDocumentStatus(
                expectedDocument.id,
                requestId,
                uploadedFiles.length > 0,
        );

        const hasDecision = DECISION_STATUSES.includes(expectedDocument.status);

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
                        toast.dismiss(loadingToast);
                        window.open(res.data, "_blank", "noopener,noreferrer");
                } catch (error: unknown) {
                        toast.error(
                                error instanceof Error ? error.message : "Eroare necunoscută",
                                { id: loadingToast },
                        );
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
                                                        canManage={canManage}
                                                        status={expectedDocument.status}
                                                        hasFiles={uploadedFiles.length > 0}
                                                        isLoading={isLoading}
                                                        requestId={requestId}
                                                        expectedDocumentId={
                                                                expectedDocument.id
                                                        }
                                                        documentTitle={expectedDocument.title}
                                                        latestFileId={
                                                                uploadedFiles[
                                                                        uploadedFiles.length -
                                                                                1
                                                                ]?.id ?? null
                                                        }
                                                        onApprove={approve}
                                                        onReject={reject}
                                                        onReset={reset}
                                                        onExtractedText={(text) =>
                                                                setExtractedText(text)
                                                        }
                                                        onInterpretedText={(text) =>
                                                                setInterpretedText(text)
                                                        }
                                                />
                                        </div>
                                </div>
                        </div>

                        {uploadedFiles.length > 0 && (
                                <div className="expected-document-slot-files">
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

                        {extractedText && (
                                <div className="extracted-text">
                                        <div className="extracted-text-header">
                                                <span>Text extras</span>
                                                <button
                                                        type="button"
                                                        className="extracted-text-close"
                                                        onClick={() => setExtractedText(null)}
                                                >
                                                        ×
                                                </button>
                                        </div>
                                        <pre className="extracted-text-content">
                                                {extractedText}
                                        </pre>
                                </div>
                        )}

                        {interpretedText && (
                                <div className="extracted-text extracted-text--interpreted">
                                        <div className="extracted-text-header">
                                                <span>Interpretare AI</span>
                                                <button
                                                        type="button"
                                                        className="extracted-text-close"
                                                        onClick={() =>
                                                                setInterpretedText(null)
                                                        }
                                                >
                                                        ×
                                                </button>
                                        </div>
                                        <pre className="extracted-text-content">
                                                {interpretedText}
                                        </pre>
                                </div>
                        )}
                </div>
        );
}
