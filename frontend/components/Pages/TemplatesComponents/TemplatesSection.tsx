"use client";
import { useEffect, useState } from "react";
import { useSearchParams } from "next/navigation";
import { Template } from "@/types";
import NotFound from "@/components/OtherComponents/NotFound/NotFound";
import SearchBar from "@/components/OtherComponents/SearchBar/SearchBar";
import TemplateCard from "./TemplateCard";
import PaginationFooter from "@/components/FileSectionComponents/FileSection/_components/PaginationFooter";
import { useSearch } from "@/hooks/useSearch";
import "./TemplatesSection.css";

interface TemplatesSectionProps {
        templates: Template[];
        isAdmin: boolean;
}

const ITEMS_PER_PAGE = 12;

export default function TemplatesSection({ templates, isAdmin }: TemplatesSectionProps) {
        const searchParams = useSearchParams();
        const [currentPage, setCurrentPage] = useState(1);

        const departmentParam = searchParams.get("department");
        const selectedDepartmentId = departmentParam ? Number(departmentParam) : null;

        const openTemplates = templates.filter((t) => !t.is_closed);

        const departmentFiltered = selectedDepartmentId
                ? openTemplates.filter((t) => t.department_id === selectedDepartmentId)
                : openTemplates;

        const { searchInput, setSearchInput, filteredItems } = useSearch(
                departmentFiltered,
                (template, search) =>
                        template.title.toLowerCase().includes(search) ||
                        (template.description ?? "").toLowerCase().includes(search),
        );

        useEffect(() => {
                setCurrentPage(1);
        }, [searchInput, selectedDepartmentId]);

        const totalPages = Math.ceil(filteredItems.length / ITEMS_PER_PAGE);
        const startIndex = (currentPage - 1) * ITEMS_PER_PAGE;
        const currentTemplates = filteredItems.slice(startIndex, startIndex + ITEMS_PER_PAGE);

        if (openTemplates.length === 0) {
                return (
                        <NotFound
                                text="Nu există niciun șablon încă."
                                subtext={
                                        isAdmin
                                                ? "Creează primul șablon pentru a permite cetățenilor să depună cereri."
                                                : "Nu există șabloane disponibile momentan."
                                }
                                background="#fff"
                        />
                );
        }

        if (departmentFiltered.length === 0) {
                return (
                        <NotFound
                                text="Niciun șablon în acest departament."
                                subtext="Nu există șabloane asociate departamentului selectat."
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
