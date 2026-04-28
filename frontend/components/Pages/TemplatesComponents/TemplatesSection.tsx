"use client";
import { useEffect, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { Template, Department, Tag } from "@/types";
import NotFound from "@/components/OtherComponents/NotFound/NotFound";
import SearchBar from "@/components/OtherComponents/SearchBar/SearchBar";
import FilterTabs from "@/components/InputComponents/FilterTabs";
import TemplateCard from "./TemplateCard";
import PaginationFooter from "@/components/FileSectionComponents/FileSection/_components/PaginationFooter";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import ManageTagsModal from "../TagsComponents/ManageTagsModal";
import { useSearch } from "@/hooks/useSearch";
import "./TemplatesSection.css";

interface TemplatesSectionProps {
        templates: Template[];
        isAdmin: boolean;
        userDepartmentId: number | null;
        departments: Department[];
        tags: Tag[];
}

const ITEMS_PER_PAGE = 12;

export default function TemplatesSection({
        templates,
        isAdmin,
        userDepartmentId,
        departments,
        tags,
}: TemplatesSectionProps) {
        const router = useRouter();
        const searchParams = useSearchParams();
        const [currentPage, setCurrentPage] = useState(1);
        const [isTagsModalOpen, setIsTagsModalOpen] = useState(false);

        const departmentParam = searchParams.get("department");
        const selectedDepartmentId = departmentParam ? Number(departmentParam) : null;

        const tagParam = searchParams.get("tag");
        const selectedTagId = tagParam ? Number(tagParam) : null;

        const openTemplates = templates.filter((t) => {
                if (t.is_closed) return false;
                if (!isAdmin && userDepartmentId !== null) {
                        return t.department_id === userDepartmentId;
                }
                return true;
        });

        const departmentFiltered = selectedDepartmentId
                ? openTemplates.filter((t) => t.department_id === selectedDepartmentId)
                : openTemplates;

        const tagFiltered = selectedTagId
                ? departmentFiltered.filter((t) =>
                          (t.tags ?? []).some((tag) => tag.id === selectedTagId),
                  )
                : departmentFiltered;

        const { searchInput, setSearchInput, filteredItems } = useSearch(
                tagFiltered,
                (template, search) =>
                        template.title.toLowerCase().includes(search) ||
                        (template.description ?? "").toLowerCase().includes(search) ||
                        template.department_name.toLowerCase().includes(search) ||
                        (template.tags ?? []).some((tag) =>
                                tag.name.toLowerCase().includes(search),
                        ),
        );

        useEffect(() => {
                setCurrentPage(1);
        }, [searchInput, selectedDepartmentId, selectedTagId]);

        const totalPages = Math.ceil(filteredItems.length / ITEMS_PER_PAGE);
        const startIndex = (currentPage - 1) * ITEMS_PER_PAGE;
        const currentTemplates = filteredItems.slice(startIndex, startIndex + ITEMS_PER_PAGE);

        const handleDepartmentChange = (deptId: string) => {
                const params = new URLSearchParams(searchParams.toString());
                if (deptId === "all") {
                        params.delete("department");
                } else {
                        params.set("department", deptId);
                }
                router.push(`?${params.toString()}`);
        };

        const handleTagChange = (tagId: string) => {
                const params = new URLSearchParams(searchParams.toString());
                if (tagId === "all") {
                        params.delete("tag");
                } else {
                        params.set("tag", tagId);
                }
                router.push(`?${params.toString()}`);
        };

        const departmentTabs = [
                { label: "Toate", value: "all", count: openTemplates.length },
                ...departments.map((d) => ({
                        label: d.name,
                        value: String(d.id),
                        count: openTemplates.filter((t) => t.department_id === d.id).length,
                })),
        ];

        const tagTabs = [
                { label: "Toate", value: "all", count: departmentFiltered.length },
                ...tags.map((t) => ({
                        label: t.name,
                        value: String(t.id),
                        count: departmentFiltered.filter((tmpl) =>
                                (tmpl.tags ?? []).some((tag) => tag.id === t.id),
                        ).length,
                })),
        ];

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

        const renderContent = () => {
                if (departmentFiltered.length === 0) {
                        return (
                                <NotFound
                                        text="Niciun șablon în acest departament."
                                        subtext="Nu există șabloane asociate departamentului selectat."
                                        background="#fff"
                                />
                        );
                }
                if (tagFiltered.length === 0) {
                        return (
                                <NotFound
                                        text="Niciun șablon cu acest tag."
                                        subtext="Nu există șabloane asociate tagului selectat."
                                        background="#fff"
                                />
                        );
                }
                if (filteredItems.length === 0) {
                        return (
                                <NotFound
                                        text="Nu am găsit niciun șablon"
                                        subtext="Nu există niciun rezultat care să corespundă căutării tale."
                                        background="#fff"
                                />
                        );
                }
                return (
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
                );
        };

        return (
                <div className="templates-section">
                        <div className="templates-toolbar">
                                {isAdmin && (
                                        <div className="templates-toolbar-actions">
                                                <ButtonPrimary
                                                        text="Șablon nou"
                                                        fullWidth
                                                        variant="primary"
                                                        onClick={() =>
                                                                router.push(
                                                                        "/dashboard/templates/create",
                                                                )
                                                        }
                                                />
                                                <ButtonPrimary
                                                        text="Gestionează taguri"
                                                        fullWidth
                                                        variant="ghost"
                                                        onClick={() =>
                                                                setIsTagsModalOpen(true)
                                                        }
                                                />
                                        </div>
                                )}
                                <SearchBar
                                        value={searchInput}
                                        onChange={setSearchInput}
                                        placeholder="Caută șablon..."
                                />
                        </div>
                        {isAdmin && (
                                <FilterTabs
                                        tabs={departmentTabs}
                                        active={
                                                selectedDepartmentId
                                                        ? String(selectedDepartmentId)
                                                        : "all"
                                        }
                                        onChange={handleDepartmentChange}
                                />
                        )}
                        {tags.length > 0 && (
                                <FilterTabs
                                        tabs={tagTabs}
                                        active={selectedTagId ? String(selectedTagId) : "all"}
                                        onChange={handleTagChange}
                                />
                        )}
                        {renderContent()}
                        {isAdmin && (
                                <ManageTagsModal
                                        isOpen={isTagsModalOpen}
                                        onClose={() => setIsTagsModalOpen(false)}
                                        tags={tags}
                                />
                        )}
                </div>
        );
}
