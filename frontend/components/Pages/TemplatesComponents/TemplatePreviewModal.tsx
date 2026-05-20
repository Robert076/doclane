"use client";
import { useState, useEffect } from "react";
import { ExpectedDocumentTemplate } from "@/types";
import { getTemplatePreview } from "@/lib/api/templates";
import Modal from "@/components/Modals/Modal";
import "./TemplatePreviewModal.css";

interface Props {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: () => void;
  templateTitle: string;
  templateDescription?: string | null;
  templateId: number;
  isSubmitting: boolean;
}

export default function TemplatePreviewModal({
  isOpen,
  onClose,
  onConfirm,
  templateTitle,
  templateDescription,
  templateId,
  isSubmitting,
}: Props) {
  const [docs, setDocs] = useState<ExpectedDocumentTemplate[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    if (!isOpen) return;
    setIsLoading(true);
    getTemplatePreview(templateId).then((res) => {
        setDocs(res.data ?? []);
        setIsLoading(false);
        });
  }, [isOpen, templateId]);

  return (
    <Modal
  isOpen={isOpen}
  onClose={onClose}
  onConfirm={onConfirm}
  title="Previzualizare cerere"
  closeOnConfirm={false}
>
      <div className="tpl-preview">
        <div className="tpl-preview-info">
          <h4 className="tpl-preview-title">{templateTitle}</h4>
          {templateDescription && (
            <p className="tpl-preview-desc">{templateDescription}</p>
          )}
        </div>

        <div className="tpl-preview-docs">
          <span className="tpl-preview-docs-label">
            Documente necesare ({docs.length}):
          </span>

          {isLoading ? (
            <p className="tpl-preview-loading">Se încarcă...</p>
          ) : docs.length === 0 ? (
            <p className="tpl-preview-empty">
              Niciun document necesar pentru acest șablon.
            </p>
          ) : (
            <ul className="tpl-preview-list">
              {docs.map((doc) => (
                <li key={doc.id} className="tpl-preview-doc">
                  <span className="tpl-preview-doc-title">{doc.title}</span>
                  {doc.description && (
                    <span className="tpl-preview-doc-desc">
                      {doc.description}
                    </span>
                  )}
                </li>
              ))}
            </ul>
          )}
        </div>

        <p className="tpl-preview-note">
          Asigură-te că ai pregătite toate documentele înainte de a depune
          cererea. Le vei putea încărca ulterior din pagina cererii.
        </p>
      </div>
    </Modal>
  );
}