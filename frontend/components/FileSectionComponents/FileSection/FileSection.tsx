"use client";
import { useState } from "react";
import "./FileSection.css";
import SectionTitle from "../../../app/dashboard/requests/[id]/_components/SectionTitle/SectionTitle";
import { DocumentFile, ExpectedDocument } from "@/types";
import ExpectedDocumentSlot from "../ExpectedDocumentSlot/ExpectedDocumentSlot";
import NotFound from "@/components/OtherComponents/NotFound/NotFound";
import PaginationFooter from "./_components/PaginationFooter";
import { UI_TEXT } from "@/locales/ro";

interface FileSectionProps {
        files: DocumentFile[];
        expectedDocuments: ExpectedDocument[];
        requestId: string;
}

const ITEMS_PER_PAGE = 2;

export default function FileSection({
        files,
        expectedDocuments,
        requestId,
}: FileSectionProps) {
        const [currentPage, setCurrentPage] = useState(1);
        const totalPages = Math.ceil(expectedDocuments.length / ITEMS_PER_PAGE);
        const startIndex = (currentPage - 1) * ITEMS_PER_PAGE;
        const currentDocs = expectedDocuments.slice(startIndex, startIndex + ITEMS_PER_PAGE);

        if (expectedDocuments.length === 0) {
                return (
                        <NotFound
                                text="No expected documents found."
                                subtext="The professional has not added any expected documents yet."
                        />
                );
        }

        return (
                <section className="details-card files-section">
                        <SectionTitle text={UI_TEXT.request.details.files} />
                        <div className="files-stack">
                                {currentDocs.map((ed) => {
                                        const uploadedFiles = files.filter(
                                                (f) => f.expected_document_id === ed.id,
                                        );
                                        return (
                                                <ExpectedDocumentSlot
                                                        key={ed.id}
                                                        expectedDocument={ed}
                                                        requestId={requestId}
                                                        uploadedFiles={uploadedFiles}
                                                />
                                        );
                                })}
                        </div>
                        {totalPages > 1 && (
                                <PaginationFooter
                                        currentPage={currentPage}
                                        totalPages={totalPages}
                                        setCurrentPage={setCurrentPage}
                                />
                        )}
                </section>
        );
}
