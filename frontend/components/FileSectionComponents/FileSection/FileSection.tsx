"use client";
import { useState } from "react";
import "./FileSection.css";
import SectionTitle from "@/components/Pages/RequestsComponents/SectionTitle";
import { DocumentFile, ExpectedDocument } from "@/types";
import ExpectedDocumentSlot from "../ExpectedDocumentSlot/ExpectedDocumentSlot";
import NotFound from "@/components/OtherComponents/NotFound/NotFound";
import PaginationFooter from "./_components/PaginationFooter";

interface FileSectionProps {
        files: DocumentFile[];
        expectedDocuments: ExpectedDocument[];
        requestId: number;
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
                                text="Nu există documente așteptate."
                                subtext="Dosarul nu are niciun document atașat."
                        />
                );
        }

        return (
                <section className="details-card files-section">
                        <SectionTitle text="Documente" />
                        <div className="files-stack">
                                {currentDocs.map((ed) => (
                                        <ExpectedDocumentSlot
                                                key={ed.id}
                                                expectedDocument={ed}
                                                requestId={requestId}
                                                uploadedFiles={files.filter(
                                                        (f) =>
                                                                f.expected_document_id ===
                                                                ed.id,
                                                )}
                                        />
                                ))}
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
