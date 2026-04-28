"use client";
import { useState } from "react";
import { Department } from "@/types";
import NotFound from "@/components/OtherComponents/NotFound/NotFound";
import SearchBar from "@/components/OtherComponents/SearchBar/SearchBar";
import PaginationFooter from "@/components/FileSectionComponents/FileSection/_components/PaginationFooter";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import { useSearch } from "@/hooks/useSearch";
import { usePagination } from "@/hooks/usePagination";
import { useRouter } from "next/navigation";
import CreateDepartmentModal from "./CreateDepartmentModal";
import DepartmentCard from "./DepartmentCard";
import "./DepartmentsSection.css";

interface DepartmentsSectionProps {
        departments: Department[];
}

const ITEMS_PER_PAGE = 12;

export default function DepartmentsSection({ departments }: DepartmentsSectionProps) {
        const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
        const router = useRouter();

        const { searchInput, setSearchInput, filteredItems } = useSearch(
                departments,
                (dept, search) => dept.name.toLowerCase().includes(search.toLowerCase()),
        );

        const { currentPage, setCurrentPage, totalPages, paginatedItems } = usePagination(
                filteredItems,
                ITEMS_PER_PAGE,
        );

        return (
                <div className="departments-section">
                        <div className="departments-toolbar">
                                <div className="departments-toolbar-action">
                                        <ButtonPrimary
                                                text="Departament nou"
                                                fullWidth
                                                variant="primary"
                                                onClick={() => setIsCreateModalOpen(true)}
                                        />
                                </div>
                                <SearchBar
                                        value={searchInput}
                                        onChange={(value) => {
                                                setSearchInput(value);
                                                setCurrentPage(1);
                                        }}
                                        placeholder="Caută departament..."
                                />
                        </div>
                        {filteredItems.length === 0 ? (
                                <NotFound
                                        text="Nu există niciun departament încă."
                                        subtext="Creează primul departament pentru a putea atribui șabloane și utilizatori."
                                        background="#fff"
                                />
                        ) : (
                                <>
                                        <div className="objects-grid">
                                                {paginatedItems.map((dept) => (
                                                        <DepartmentCard
                                                                key={dept.id}
                                                                department={dept}
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
                        <CreateDepartmentModal
                                isOpen={isCreateModalOpen}
                                onClose={() => setIsCreateModalOpen(false)}
                                onCreated={() => router.refresh()}
                        />
                </div>
        );
}
