"use client";
import { useState } from "react";
import DetailsCard from "./DetailsCard";
import FileSection from "@/components/FileSectionComponents/FileSection/FileSection";
import RequestComments from "./RequestComments";
import { Request, DocumentFile, RequestComment } from "@/types";
import "./RequestTabs.css";

type Tab = "details" | "files" | "comments";

interface RequestTabsProps {
        data: Request;
        files: DocumentFile[];
        comments: RequestComment[];
        requestId: number;
}

export default function RequestTabs({ data, files, comments, requestId }: RequestTabsProps) {
        const [active, setActive] = useState<Tab>("details");

        return (
                <div className="request-tabs">
                        <div className="tab-bar">
                                <button
                                        className={`tab-btn ${active === "details" ? "tab-btn--active" : ""}`}
                                        onClick={() => setActive("details")}
                                >
                                        Detalii
                                </button>
                                <button
                                        className={`tab-btn ${active === "files" ? "tab-btn--active" : ""}`}
                                        onClick={() => setActive("files")}
                                >
                                        Fișiere
                                </button>
                                <button
                                        className={`tab-btn ${active === "comments" ? "tab-btn--active" : ""}`}
                                        onClick={() => setActive("comments")}
                                >
                                        Comentarii
                                </button>
                        </div>
                        <div className="tab-content">
                                {active === "details" && <DetailsCard data={data} />}
                                {active === "files" && (
                                        <FileSection
                                                files={files}
                                                expectedDocuments={
                                                        data.expected_documents ?? []
                                                }
                                                requestId={requestId}
                                        />
                                )}
                                {active === "comments" && (
                                        <RequestComments
                                                comments={comments}
                                                requestId={requestId}
                                        />
                                )}
                        </div>
                </div>
        );
}
