"use client";
import { useEffect, useState } from "react";
import { Template } from "@/types";
import NotFound from "@/components/OtherComponents/NotFound/NotFound";
import SearchBar from "@/components/OtherComponents/SearchBar/SearchBar";
import PaginationFooter from "@/components/ClientComponents/ClientsSection/_components/PaginationFooter";
import TemplateCard from "./TemplateCard";
import { UI_TEXT } from "@/locales/ro";
import "./TemplatesSection.css";
import { useSearch } from "@/hooks/useSearch";

interface TemplatesSectionProps {
        templates: Template[];
}

const ITEMS_PER_PAGE = 12;

export default function TemplatesSection({ templates }: TemplatesSectionProps) {
        const [currentPage, setCurrentPage] = useState(1);

        const openTemplates = templates.filter((t) => !t.is_closed);

        const {
                searchInput,
                setSearchInput,
                filteredItems: filteredTemplates,
        } = useSearch(openTemplates, (template, search) => {
                return (
                        template.title.toLowerCase().includes(search) ||
                        (template.description ?? "").toLowerCase().includes(search)
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
        if (openTemplates.length === 0) {
                return (
                        <NotFound
                                text="Nu ai niciun şablon încă."
                                subtext="Aici vor apărea şabloanele pe care le salvezi pentru a fi refolosite."
                                background="#fff"
                        />
                );
        }

        return (
                <div className="templates-section">
                        {openTemplates.length > 0 && (
                                <SearchBar
                                        value={searchInput}
                                        onChange={setSearchInput}
                                        placeholder={UI_TEXT.common.search}
                                />
                        )}

                        {filteredTemplates.length === 0 && openTemplates.length > 0 && (
                                <NotFound
                                        text="Nu am găsit niciun șablon"
                                        subtext="Nu există niciun rezultat care să corespundă căutarii tale."
                                        background="#fff"
                                />
                        )}

                        <div className="templates-grid">
                                {currentTemplates.map((template) => (
                                        <TemplateCard
                                                key={template.id}
                                                template={template}
                                                searchTerm={searchInput}
                                                archived={false}
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
