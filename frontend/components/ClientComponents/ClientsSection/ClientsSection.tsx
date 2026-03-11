"use client";

import { useCallback } from "react";
import { User } from "@/types";
import ClientCard from "../ClientCard/ClientCard";
import NotFound from "@/components/OtherComponents/NotFound/NotFound";
import PaginationFooter from "./_components/PaginationFooter";
import SearchBar from "@/components/OtherComponents/SearchBar/SearchBar";
import { useSearch } from "@/hooks/useSearch";
import { usePagination } from "@/hooks/usePagination";
import "./ClientsSection.css";

interface ClientsSectionProps {
        clients: User[];
}

const ITEMS_PER_PAGE = 12;

export default function ClientsSection({ clients }: ClientsSectionProps) {
        // 1. Funcția de căutare optimizată (exact ca la dosare)
        const searchFn = useCallback((client: User, search: string) => {
                const searchLower = search.toLowerCase();

                const searchableText = [client.first_name, client.last_name, client.email]
                        .filter(Boolean)
                        .join(" ")
                        .toLowerCase();

                return searchableText.includes(searchLower);
        }, []);

        // 2. Aplicăm Hook-urile (asta rezolvă bug-urile de randare și paginare)
        const { searchInput, setSearchInput, filteredItems } = useSearch(clients, searchFn);

        const { currentPage, setCurrentPage, totalPages, paginatedItems } = usePagination(
                filteredItems,
                ITEMS_PER_PAGE,
        );

        // Caz 1: Utilizatorul nu are niciun client în baza de date
        if (clients.length === 0) {
                return (
                        <NotFound
                                text="Nu ai niciun client momentan."
                                subtext="Începe prin a adăuga primul tău client folosind butonul de mai sus."
                                background="#fff"
                        />
                );
        }

        return (
                <div className="clients-section">
                        <SearchBar
                                value={searchInput}
                                onChange={(value) => {
                                        setSearchInput(value);
                                        setCurrentPage(1); // Resetăm pagina automat la căutare
                                }}
                                placeholder="Caută..."
                        />

                        {/* Caz 2: Are clienți, dar n-a găsit nimic la căutare */}
                        {filteredItems.length === 0 ? (
                                <NotFound
                                        text="Nu am găsit niciun client."
                                        subtext="Nu există niciun rezultat care să corespundă căutării tale."
                                        background="#fff"
                                />
                        ) : (
                                <>
                                        <div className="clients-grid">
                                                {/* Randăm DOAR elementele paginii curente */}
                                                {paginatedItems.map((client) => (
                                                        <ClientCard
                                                                key={client.id}
                                                                client={client}
                                                                searchTerm={searchInput}
                                                        />
                                                ))}
                                        </div>

                                        {/* Afișăm paginarea doar dacă avem mai mult de o pagină */}
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
