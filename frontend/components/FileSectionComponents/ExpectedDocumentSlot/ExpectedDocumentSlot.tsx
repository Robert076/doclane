"use client";
import { useState } from "react";
import "./ExpectedDocumentSlot.css";
import { DocumentFile, ExpectedDocument, RequestStatus } from "@/types";
import FileItem from "../FileItem/FileItem";
import PaginationFooter from "@/components/FileSectionComponents/FileSection/_components/PaginationFooter";
import StatusBadge from "@/components/Pages/RequestsComponents/StatusBadge";
import { useUser } from "@/context/UserContext";
import { useDocumentStatus } from "./_hooks/useDocumentStatus";
import DocumentSlotActions from "./_components/DocumentSlotActions";
import { presignExampleURL } from "@/lib/api/api";
import { MdFileOpen } from "react-icons/md";
import toast from "react-hot-toast";

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
        const [currentPage, setCurrentPage] = useState(1);
        const [isLoadingExample, setIsLoadingExample] = useState(false);

        const totalPages = Math.ceil(uploadedFiles.length / ITEMS_PER_PAGE);
        const startIndex = (currentPage - 1) * ITEMS_PER_PAGE;
        const currentFiles = uploadedFiles.slice(startIndex, startIndex + ITEMS_PER_PAGE);

        const user = useUser();
        const isProfessional = user?.role === "PROFESSIONAL";

        const { approve, reject, reset, isLoading } = useDocumentStatus(
                expectedDocument.id.toString(),
                requestId,
                uploadedFiles.length > 0,
        );

        const hasDecision = DECISION_STATUSES.includes(expectedDocument.status);

        const handleViewExample = async () => {
                setIsLoadingExample(true);
                try {
                        const res = await presignExampleURL(expectedDocument.id);
                        if (!res.success) {
                                toast.error("Nu s-a putut deschide exemplul.");
                                return;
                        }
                        window.open(res.data.url, "_blank");
                } catch {
                        toast.error("Nu s-a putut deschide exemplul.");
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
                                                {expectedDocument.rejection_reason &&
                                                        expectedDocument.status ===
                                                                "rejected" && (
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
                                                        requestId={+requestId}
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
                                        {currentFiles.map((file) => (
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
