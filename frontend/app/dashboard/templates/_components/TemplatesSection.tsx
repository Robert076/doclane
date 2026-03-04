"use client";
import { useEffect, useState } from "react";
import { DocumentRequestTemplate } from "@/types";
import NotFound from "@/components/OtherComponents/NotFound/NotFound";
import SearchBar from "@/components/OtherComponents/SearchBar/SearchBar";
import PaginationFooter from "@/components/ClientComponents/ClientsSection/_components/PaginationFooter";

import { UI_TEXT } from "@/locales/ro";
import "./TemplatesSection.css";
import TemplateCard from "./TemplateCard";

interface TemplatesSectionProps {
        templates: DocumentRequestTemplate[];
}

const ITEMS_PER_PAGE = 12;

export default function TemplatesSection({ templates }: TemplatesSectionProps) {
        const [currentPage, setCurrentPage] = useState<number>(1);
        const [searchInput, setSearchInput] = useState<string>("");

        const filteredTemplates = templates.filter((template) => {
                if (template.is_closed) return false;
                if (!searchInput) return true;
                const searchLower = searchInput.toLowerCase().trim();
                return (
                        template.title.toLowerCase().includes(searchLower) ||
                        template.description?.toLowerCase().includes(searchLower)
                );
        });

        useEffect(() => {
                setCurrentPage(1);
        }, [searchInput]);

        const totalPages = Math.ceil(filteredTemplates.length / ITEMS_PER_PAGE);
        const startIndex = (currentPage - 1) * ITEMS_PER_PAGE;
        const currentTemplates = filteredTemplates.slice(
                startIndex,
                startIndex + ITEMS_PER_PAGE,
        );

        return (
                <div className="templates-section">
                        {templates.length > 0 && (
                                <SearchBar
                                        value={searchInput}
                                        onChange={setSearchInput}
                                        placeholder={UI_TEXT.common.search}
                                />
                        )}
                        {templates.length === 0 && (
                                <NotFound
                                        text="Nu ai niciun şablon încă."
                                        subtext="Începe prin a crea primul şablon."
                                        background="#fff"
                                />
                        )}
                        {filteredTemplates.length === 0 && templates.length > 0 && (
                                <NotFound
                                        text={UI_TEXT.common.searchNotFound}
                                        subtext=""
                                        background="#fff"
                                />
                        )}
                        <div className="templates-grid">
                                {currentTemplates.map((template) => (
                                        <TemplateCard
                                                key={template.id}
                                                template={template}
                                                searchTerm={searchInput}
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
                </div>
        );
}
