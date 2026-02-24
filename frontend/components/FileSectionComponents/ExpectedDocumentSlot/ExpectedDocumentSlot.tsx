"use client";
import { useState } from "react";
import "./ExpectedDocumentSlot.css";
import UploadDocumentButton from "@/app/dashboard/requests/[id]/_components/UploadDocumentButton/UploadDocumentButton";
import { DocumentFile, ExpectedDocument } from "@/types";
import FileItem from "../FileItem/FileItem";
import PaginationFooter from "@/components/FileSectionComponents/FileSection/_components/PaginationFooter";

interface ExpectedDocumentSlotProps {
        expectedDocument: ExpectedDocument;
        requestId: string;
        uploadedFiles: DocumentFile[];
}

const ITEMS_PER_PAGE = 3;

export default function ExpectedDocumentSlot({
        expectedDocument,
        requestId,
        uploadedFiles,
}: ExpectedDocumentSlotProps) {
        const [currentPage, setCurrentPage] = useState(1);
        const totalPages = Math.ceil(uploadedFiles.length / ITEMS_PER_PAGE);
        const startIndex = (currentPage - 1) * ITEMS_PER_PAGE;
        const currentFiles = uploadedFiles.slice(startIndex, startIndex + ITEMS_PER_PAGE);

        return (
                <div className="expected-document-slot">
                        <div className="expected-document-slot-top">
                                <div className="expected-document-slot-info">
                                        <span className="expected-document-slot-title">
                                                {expectedDocument.title}
                                        </span>
                                        {expectedDocument.description && (
                                                <span className="expected-document-slot-description">
                                                        {expectedDocument.description}
                                                </span>
                                        )}
                                        {expectedDocument.notes && (
                                                <span className="expected-document-slot-notes">
                                                        📋 {expectedDocument.notes}
                                                </span>
                                        )}
                                </div>
                                <div className="upload-button-wrapper">
                                        <UploadDocumentButton
                                                requestId={requestId}
                                                expectedDocumentId={expectedDocument.id}
                                        />
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
