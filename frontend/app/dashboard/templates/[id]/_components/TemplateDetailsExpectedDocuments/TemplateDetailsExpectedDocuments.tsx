"use client";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import { presignTemplateExample } from "@/lib/api/api";
import { ExpectedDocumentTemplate } from "@/types";
import { useState } from "react";
import toast from "react-hot-toast";

export default function TemplateDetailsExpectedDocuments({
        documents,
        templateID,
}: {
        documents: ExpectedDocumentTemplate[];
        templateID: number;
}) {
        const [isLoadingExample, setIsLoadingExample] = useState<boolean>(false);
        const handleViewExample = async (id: number) => {
                setIsLoadingExample(true);
                try {
                        const res = await presignTemplateExample(templateID, id);
                        console.log(res);
                        if (!res.success) {
                                toast.error("Nu s-a putut deschide exemplul.");
                                return;
                        }
                        window.open(res.data, "_blank");
                } catch {
                        toast.error("Nu s-a putut deschide exemplul.");
                } finally {
                        setIsLoadingExample(false);
                }
        };

        return (
                <div className="template-details-expected-documents">
                        {documents.map((document: ExpectedDocumentTemplate) => (
                                <div key={document.id}>
                                        <div className="title">{document.title}</div>
                                        {document.example_file_path && (
                                                <ButtonPrimary
                                                        variant="primary"
                                                        text="Deschide"
                                                        onClick={() => {
                                                                handleViewExample(document.id);
                                                        }}
                                                />
                                        )}
                                </div>
                        ))}
                </div>
        );
}
