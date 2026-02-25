"use client";
import { useState } from "react";
import Input from "@/components/InputComponents/Input";
import TextArea from "@/components/InputComponents/TextArea";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import { MdAdd, MdClose } from "react-icons/md";
import "./ExpectedDocumentsList.css";
import { UI_TEXT } from "@/locales/ro";

export interface ExpectedDocumentInput {
        title: string;
        description: string;
}

interface ExpectedDocumentsListProps {
        documents: ExpectedDocumentInput[];
        onChange: (documents: ExpectedDocumentInput[]) => void;
}

const ExpectedDocumentsList: React.FC<ExpectedDocumentsListProps> = ({
        documents,
        onChange,
}) => {
        const handleAdd = () => {
                onChange([...documents, { title: "", description: "" }]);
        };

        const handleRemove = (index: number) => {
                onChange(documents.filter((_, i) => i !== index));
        };

        const handleChange = (
                index: number,
                field: keyof ExpectedDocumentInput,
                value: string,
        ) => {
                const updated = documents.map((doc, i) =>
                        i === index ? { ...doc, [field]: value } : doc,
                );
                onChange(updated);
        };

        return (
                <div className="expected-documents-list">
                        <div className="expected-documents-header">
                                <label className="expected-documents-label">
                                        {UI_TEXT.request.createForm.expectedDocuments}
                                </label>
                                <ButtonPrimary
                                        text={UI_TEXT.request.createForm.addExpectedDocument}
                                        variant="ghost"
                                        icon={MdAdd}
                                        type="button"
                                        onClick={handleAdd}
                                />
                        </div>

                        {documents.length === 0 && (
                                <p className="expected-documents-empty">
                                        {UI_TEXT.request.createForm.expectedDocumentsNotAdded}
                                </p>
                        )}

                        {documents.map((doc, index) => (
                                <div key={index} className="expected-document-item">
                                        <div className="expected-document-item-header">
                                                <span className="expected-document-number">
                                                        Document {index + 1}
                                                </span>
                                                <button
                                                        type="button"
                                                        className="expected-document-remove"
                                                        onClick={() => handleRemove(index)}
                                                >
                                                        <MdClose />
                                                </button>
                                        </div>
                                        <Input
                                                label={
                                                        UI_TEXT.request.createForm
                                                                .expectedDocumentTitle
                                                }
                                                placeholder={
                                                        UI_TEXT.request.createForm
                                                                .expectedDocumentTitlePlaceholder
                                                }
                                                value={doc.title}
                                                onChange={(e: any) =>
                                                        handleChange(
                                                                index,
                                                                "title",
                                                                e.target.value,
                                                        )
                                                }
                                        />
                                        <TextArea
                                                label={
                                                        UI_TEXT.request.createForm
                                                                .expectedDocumentDescription
                                                }
                                                placeholder={
                                                        UI_TEXT.request.createForm
                                                                .expectedDocumentDescriptionPlaceholder
                                                }
                                                value={doc.description}
                                                onChange={(e: any) =>
                                                        handleChange(
                                                                index,
                                                                "description",
                                                                e.target.value,
                                                        )
                                                }
                                                minHeight={200}
                                                maxHeight={200}
                                        />
                                </div>
                        ))}
                </div>
        );
};

export default ExpectedDocumentsList;
