"use client";

import { MdAdd, MdClose, MdUploadFile } from "react-icons/md";

// Componente UI
import Input from "@/components/InputComponents/Input";
import TextArea from "@/components/InputComponents/TextArea";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import "./ExpectedDocumentsList.css";

// Reutilizăm interfața din requests.ts (dacă ai mutat-o acolo, sau o declari aici dacă vrei)
// Am extins denumirea să se potrivească
export interface ExpectedDocumentInput {
        title: string;
        description: string;
        exampleFile?: File;
}

interface ExpectedDocumentsListProps {
        documents: ExpectedDocumentInput[];
        onChange: (documents: ExpectedDocumentInput[]) => void;
}

export default function ExpectedDocumentsList({
        documents,
        onChange,
}: ExpectedDocumentsListProps) {
        // --- Handlere pentru gestionarea listei ---

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

        const handleExampleFileChange = (index: number, file: File | undefined) => {
                const updated = documents.map((doc, i) =>
                        i === index ? { ...doc, exampleFile: file } : doc,
                );
                onChange(updated);
        };

        return (
                <div className="expected-documents-list">
                        <div className="expected-documents-header">
                                <label className="expected-documents-label">
                                        Documente necesare
                                </label>
                                <ButtonPrimary
                                        text="Adaugă un document necesar"
                                        variant="ghost"
                                        icon={MdAdd}
                                        type="button"
                                        onClick={handleAdd}
                                />
                        </div>

                        {/* Stare Goală */}
                        {documents.length === 0 && (
                                <p className="expected-documents-empty">
                                        Nu ai adăugat niciun document necesar. Te rugăm să
                                        adaugi cel puțin unul.
                                </p>
                        )}

                        {/* Lista de documente */}
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
                                                        title="Șterge documentul"
                                                >
                                                        <MdClose />
                                                </button>
                                        </div>

                                        <Input
                                                label="Titlu document"
                                                placeholder="Ex: Carte de identitate, Extras de cont..."
                                                value={doc.title}
                                                onChange={(
                                                        e: React.ChangeEvent<HTMLInputElement>,
                                                ) =>
                                                        handleChange(
                                                                index,
                                                                "title",
                                                                e.target.value,
                                                        )
                                                }
                                        />

                                        <TextArea
                                                label="Descriere (opțional)"
                                                placeholder="Adaugă detalii despre ce informații trebuie să conțină documentul..."
                                                value={doc.description}
                                                onChange={(
                                                        e: React.ChangeEvent<HTMLTextAreaElement>,
                                                ) =>
                                                        handleChange(
                                                                index,
                                                                "description",
                                                                e.target.value,
                                                        )
                                                }
                                                minHeight={120} // Scăzut puțin pentru a nu lungi formularul aiurea
                                                maxHeight={200}
                                        />

                                        {/* Secțiunea de Upload Exemplu */}
                                        <div className="expected-document-example">
                                                <label className="expected-document-example-label">
                                                        Exemplu de completare (Opțional)
                                                </label>

                                                {doc.exampleFile ? (
                                                        <div className="expected-document-example-preview">
                                                                <span className="expected-document-example-name">
                                                                        {doc.exampleFile.name}
                                                                </span>
                                                                <button
                                                                        type="button"
                                                                        className="expected-document-remove"
                                                                        onClick={() =>
                                                                                handleExampleFileChange(
                                                                                        index,
                                                                                        undefined,
                                                                                )
                                                                        }
                                                                        title="Șterge exemplul"
                                                                >
                                                                        <MdClose />
                                                                </button>
                                                        </div>
                                                ) : (
                                                        <label className="expected-document-example-upload">
                                                                <MdUploadFile />
                                                                <span>Încarcă un exemplu</span>
                                                                <input
                                                                        type="file"
                                                                        accept=".pdf,.jpg,.jpeg,.png,.doc,.docx"
                                                                        style={{
                                                                                display: "none",
                                                                        }}
                                                                        onChange={(
                                                                                e: React.ChangeEvent<HTMLInputElement>,
                                                                        ) => {
                                                                                const file =
                                                                                        e
                                                                                                .target
                                                                                                .files?.[0];
                                                                                if (file)
                                                                                        handleExampleFileChange(
                                                                                                index,
                                                                                                file,
                                                                                        );
                                                                        }}
                                                                />
                                                        </label>
                                                )}
                                        </div>
                                </div>
                        ))}
                </div>
        );
}
