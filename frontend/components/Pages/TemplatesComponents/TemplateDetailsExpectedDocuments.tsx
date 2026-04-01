"use client";
import { ExpectedDocumentTemplate } from "@/types";
import { useState } from "react";
import toast from "react-hot-toast";
import "./TemplateDetailsExpectedDocuments.css";
import SectionTitle from "@/components/Pages/RequestsComponents/SectionTitle";
import { presignTemplateExample } from "@/lib/api/templates";
import { MdFileOpen } from "react-icons/md";

export default function TemplateDetailsExpectedDocuments({
        documents,
        templateID,
}: {
        documents: ExpectedDocumentTemplate[];
        templateID: number;
}) {
        const [loadingId, setLoadingId] = useState<number | null>(null);

        const handleViewExample = async (id: number) => {
                setLoadingId(id);
                const loadingToast = toast.loading("Se deschide exemplul...");

                try {
                        const res = await presignTemplateExample(templateID, id);

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
                        toast.error(error instanceof Error ? error.message : "Eroare necunoscuta", { id: loadingToast });
                } finally {
                        setLoadingId(null);
                }
        };

        return (
                <div className="template-details-expected-documents">
                        <SectionTitle text="Documente şablon" />
                        {documents.map((document: ExpectedDocumentTemplate) => (
                                <div className="template-expected-document" key={document.id}>
                                        <div className="template-expected-document-info">
                                                <span className="template-expected-document-title">
                                                        {document.title}
                                                </span>
                                                {document.description && (
                                                        <span className="template-expected-document-description">
                                                                {document.description}
                                                        </span>
                                                )}
                                                {document.example_file_path && (
                                                        <button
                                                                type="button"
                                                                className="expected-document-slot-example-btn"
                                                                onClick={() =>
                                                                        handleViewExample(
                                                                                document.id,
                                                                        )
                                                                }
                                                                disabled={
                                                                        loadingId ===
                                                                        document.id
                                                                }
                                                        >
                                                                <MdFileOpen />
                                                                {loadingId === document.id
                                                                        ? "Se încarcă..."
                                                                        : "Vezi exemplu"}
                                                        </button>
                                                )}
                                        </div>
                                </div>
                        ))}
                </div>
        );
}
