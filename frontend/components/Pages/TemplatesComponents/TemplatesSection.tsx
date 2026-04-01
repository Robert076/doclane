"use client";
import { useEffect, useState } from "react";
import { Template } from "@/types";
import NotFound from "@/components/OtherComponents/NotFound/NotFound";
import SearchBar from "@/components/OtherComponents/SearchBar/SearchBar";
import TemplateCard from "./TemplateCard";
import { useSearch } from "@/hooks/useSearch";
import { useUser } from "@/context/UserContext";
import "./TemplatesSection.css";
import PaginationFooter from "@/components/FileSectionComponents/FileSection/_components/PaginationFooter";

interface TemplatesSectionProps {
        templates: Template[];
}

const ITEMS_PER_PAGE = 12;

export default function TemplatesSection({ templates }: TemplatesSectionProps) {
        const user = useUser();
        const [currentPage, setCurrentPage] = useState(1);
        const canManage = user.role === "admin" || user.department_id !== null;

        const openTemplates = templates.filter((t) => !t.is_closed);

        const { searchInput, setSearchInput, filteredItems } = useSearch(
                openTemplates,
                (template, search) =>
                        template.title.toLowerCase().includes(search) ||
                        (template.description ?? "").toLowerCase().includes(search),
        );

        useEffect(() => {
                setCurrentPage(1);
        }, [searchInput]);

        const totalPages = Math.ceil(filteredItems.length / ITEMS_PER_PAGE);
        const startIndex = (currentPage - 1) * ITEMS_PER_PAGE;
        const currentTemplates = filteredItems.slice(startIndex, startIndex + ITEMS_PER_PAGE);

        if (openTemplates.length === 0) {
                return (
                        <NotFound
                                text="Nu există niciun șablon încă."
                                subtext={
                                        canManage
                                                ? "Creează primul șablon pentru a permite cetățenilor să depună cereri."
                                                : "Nu există șabloane disponibile momentan."
                                }
                                background="#fff"
                        />
                );
        }

        return (
                <div className="templates-section">
                        <SearchBar
                                value={searchInput}
                                onChange={setSearchInput}
                                placeholder="Caută șablon..."
                        />
                        {filteredItems.length === 0 ? (
                                <NotFound
                                        text="Nu am găsit niciun șablon"
                                        subtext="Nu există niciun rezultat care să corespundă căutării tale."
                                        background="#fff"
                                />
                        ) : (
                                <>
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
                                </>
                        )}
                </div>
        );
}
