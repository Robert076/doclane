"use client";

import { useMemo, useCallback } from "react";
import { Request, User } from "@/types";
import RequestCard from "@/components/CardComponents/RequestCard/RequestCard";
import NotFound from "@/components/OtherComponents/NotFound/NotFound";
import SearchBar from "@/components/OtherComponents/SearchBar/SearchBar";
import PaginationFooter from "@/components/ClientComponents/ClientsSection/_components/PaginationFooter";
import { useSearch } from "@/hooks/useSearch";
import { usePagination } from "@/hooks/usePagination";

interface RequestsSectionProps {
        requests: Request[];
        user: User;
}

const ITEMS_PER_PAGE = 8;

export default function RequestsSection({ requests, user }: RequestsSectionProps) {
        // 1. Memorăm lista dosarelor deschise
        const openRequests = useMemo(() => {
                return requests.filter((r) => r.is_closed === false);
        }, [requests]);

        // 2. Optimizăm funcția de căutare
        const searchFn = useCallback((req: Request, search: string) => {
                const searchLower = search.toLowerCase();

                const searchableText = [
                        req.title,
                        req.description,
                        req.client_first_name,
                        req.client_last_name,
                        req.client_email,
                        req.status,
                ]
                        .filter(Boolean)
                        .join(" ")
                        .toLowerCase();

                return searchableText.includes(searchLower);
        }, []);

        const { searchInput, setSearchInput, filteredItems } = useSearch(
                openRequests,
                searchFn,
        );

        const { currentPage, setCurrentPage, totalPages, paginatedItems } = usePagination(
                filteredItems,
                ITEMS_PER_PAGE,
        );

        // Caz 1: Utilizatorul nu are niciun dosar asociat
        if (requests.length === 0) {
                const isProfessional = user.role === "PROFESSIONAL";
                return (
                        <NotFound
                                text="Nu ai niciun dosar încă."
                                subtext={
                                        isProfessional
                                                ? "Aici vor apărea dosarele pe care le gestionezi pentru clienții tăi."
                                                : "Aici vor apărea dosarele pe care profesioniștii le-au deschis pentru tine."
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

                        {/* Caz 2: Are dosare, dar căutarea nu returnează nimic */}
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

                                        <PaginationFooter
                                                totalPages={totalPages}
                                                currentPage={currentPage}
                                                setCurrentPage={setCurrentPage}
                                        />
                                </>
                        )}
                </div>
        );
}
