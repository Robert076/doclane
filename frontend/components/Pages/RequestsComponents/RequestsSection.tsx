"use client";

import { DocumentRequest, User } from "@/types";
import RequestCard from "@/components/CardComponents/RequestCard/RequestCard";
import NotFound from "@/components/OtherComponents/NotFound/NotFound";
import SearchBar from "@/components/OtherComponents/SearchBar/SearchBar";
import PaginationFooter from "@/components/ClientComponents/ClientsSection/_components/PaginationFooter";
import { UI_TEXT } from "@/locales/ro";
import { useSearch } from "@/hooks/useSearch";
import { usePagination } from "@/hooks/usePagination";

interface RequestsSectionProps {
        requests: DocumentRequest[];
        user: User;
}

const ITEMS_PER_PAGE = 8;

export default function RequestsSection({ requests, user }: RequestsSectionProps) {
        const openRequests = requests.filter((r) => {
                if (r.is_closed === true) return false;
                return true;
        });
        const searchFn = (req: DocumentRequest, search: string) => {
                if (req.is_closed) return false;

                const fullName =
                        `${req.client_first_name ?? ""} ${req.client_last_name ?? ""}`.toLowerCase();

                return [
                        req.title,
                        req.description,
                        req.client_first_name,
                        req.client_last_name,
                        req.client_email,
                        req.status,
                        fullName,
                ]
                        .filter(Boolean)
                        .some((field) => field!.toLowerCase().includes(search));
        };

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
                                text={UI_TEXT.common.notFoundTitleRequests}
                                subtext={
                                        user.role === "PROFESSIONAL"
                                                ? UI_TEXT.common
                                                          .notFoundSubtitleRequestsProfessional
                                                : UI_TEXT.common.notFoundSubtitleRequestsClient
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
                                placeholder={UI_TEXT.common.search}
                        />

                        {filteredItems.length === 0 ? (
                                <NotFound
                                        text={UI_TEXT.common.notFoundTitleRequests}
                                        subtext={UI_TEXT.common.searchNotFound}
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
