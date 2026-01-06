"use client";

import { useState } from "react";
import FileItem from "../FileItem/FileItem";
import "./FileSection.css";
import SectionTitle from "../SectionTitle/SectionTitle";

interface FileSectionProps {
  files: any[];
}

export default function FileSection({ files }: FileSectionProps) {
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 4;

  const totalPages = Math.ceil(files.length / itemsPerPage);

  const startIndex = (currentPage - 1) * itemsPerPage;
  const currentFiles = files.slice(startIndex, startIndex + itemsPerPage);

  if (files.length === 0) {
    return (
      <section className="details-card files-section">
        <h2 className="section-title">Documents</h2>
        <div className="empty-files">
          <p>No documents have been uploaded yet.</p>
        </div>
      </section>
    );
  }

  return (
    <section className="details-card files-section">
      <div className="section-header">
        <SectionTitle text="Documents" />
        {totalPages > 1 && (
          <span className="pagination-counter">
            {startIndex + 1}-{Math.min(startIndex + itemsPerPage, files.length)} din{" "}
            {files.length}
          </span>
        )}
      </div>

      <div className="files-container">
        <div className="files-stack">
          {currentFiles.map((file) => (
            <FileItem key={file.id} file={file} />
          ))}
        </div>
      </div>

      {totalPages > 1 && (
        <div className="pagination-footer">
          <button
            className="pag-button"
            onClick={() => setCurrentPage((p) => Math.max(1, p - 1))}
            disabled={currentPage === 1}
          >
            Anterior
          </button>

          <div className="page-dots">
            {Array.from({ length: totalPages }).map((_, i) => (
              <div key={i} className={`dot ${currentPage === i + 1 ? "active" : ""}`} />
            ))}
          </div>

          <button
            className="pag-button"
            onClick={() => setCurrentPage((p) => Math.min(totalPages, p + 1))}
            disabled={currentPage === totalPages}
          >
            Următor
          </button>
        </div>
      )}
    </section>
  );
}
