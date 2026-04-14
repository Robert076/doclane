"use client";
import { useMemo, useCallback } from "react";
import { Request, User } from "@/types";
import RequestCard from "@/components/CardComponents/RequestCard/RequestCard";
import NotFound from "@/components/OtherComponents/NotFound/NotFound";
import SearchBar from "@/components/OtherComponents/SearchBar/SearchBar";
import { useSearch } from "@/hooks/useSearch";
import { usePagination } from "@/hooks/usePagination";
import PaginationFooter from "@/components/FileSectionComponents/FileSection/_components/PaginationFooter";

interface RequestsSectionProps {
        requests: Request[];
        user: User;
}

const ITEMS_PER_PAGE = 8;

export default function RequestsSection({ requests, user }: RequestsSectionProps) {
        const openRequests = useMemo(
                () => requests.filter((r) => !r.is_closed && !r.is_cancelled),
                [requests],
        );

        const searchFn = useCallback((req: Request, search: string) => {
                const searchLower = search.toLowerCase();
                return [
                        req.title,
                        req.description,
                        req.department_name,
                        req.assignee_first_name,
                        req.assignee_last_name,
                        req.assignee_email,
                        req.status,
                ]
                        .filter(Boolean)
                        .join(" ")
                        .toLowerCase()
                        .includes(searchLower);
        }, []);

        const { searchInput, setSearchInput, filteredItems } = useSearch(
                openRequests,
                searchFn,
        );
        const { currentPage, setCurrentPage, totalPages, paginatedItems } = usePagination(
                filteredItems,
                ITEMS_PER_PAGE,
        );

        if (requests.length === 0) {
                return (
                        <NotFound
                                text="Nu ai niciun dosar încă."
                                subtext={
                                        user.role === "admin" || user.department_id !== null
                                                ? "Aici vor apărea dosarele pe care le gestionezi."
                                                : "Aici vor apărea dosarele tale. Accesează șabloanele pentru a deschide un dosar."
                                }
                                background="#fff"
                        />
                );
        }

        return (
                <div className="section">
                        <SearchBar
                                value={searchInput}
                                onChange={(value) => {
                                        setSearchInput(value);
                                        setCurrentPage(1);
                                }}
                                placeholder="Caută..."
                        />
                        {filteredItems.length === 0 ? (
                                <NotFound
                                        text="Nu am găsit niciun dosar"
                                        subtext="Nu există niciun rezultat care să corespundă căutării tale."
                                        background="#fff"
                                />
                        ) : (
                                <>
                                        <div className="objects-grid">
                                                {paginatedItems.map((req) => (
                                                        <RequestCard
                                                                key={req.id}
                                                                request={req}
                                                                user={user}
                                                                searchTerm={searchInput}
                                                                archived={false}
                                                        />
                                                ))}
                                        </div>
                                        {filteredItems.length >= ITEMS_PER_PAGE && (
                                                <PaginationFooter
                                                        totalPages={totalPages}
                                                        currentPage={currentPage}
                                                        setCurrentPage={setCurrentPage}
                                                />
                                        )}
                                </>
                        )}
                </div>
        );
}
