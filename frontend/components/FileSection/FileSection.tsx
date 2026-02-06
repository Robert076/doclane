"use client";
import { useState } from "react";
import FileItem from "../FileItem/FileItem";
import NotFound from "@/components/NotFound/NotFound";
import "./FileSection.css";
import SectionTitle from "../SectionTitle/SectionTitle";
import PaginationCounter from "./_components/PaginationCounter";
import PaginationFooter from "./_components/PaginationFooter";
import { DocumentFile } from "@/types";

interface FileSectionProps {
  files: DocumentFile[];
}

export default function FileSection({ files }: FileSectionProps) {
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 3;
  const totalPages = Math.ceil(files.length / itemsPerPage);
  const startIndex = (currentPage - 1) * itemsPerPage;
  const currentFiles = files.slice(startIndex, startIndex + itemsPerPage);

  if (files.length === 0) {
    return <NotFound text="No documents found." subtext="Upload documents to get started." />;
  }

  return (
    <section className="details-card files-section">
      <div className="section-header">
        <SectionTitle text="Documents" />
        {totalPages > 1 && (
          <PaginationCounter
            startIndex={startIndex}
            itemsPerPage={itemsPerPage}
            length={files.length}
          />
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
        <PaginationFooter
          currentPage={currentPage}
          totalPages={totalPages}
          setCurrentPage={setCurrentPage}
        />
      )}
    </section>
  );
}
