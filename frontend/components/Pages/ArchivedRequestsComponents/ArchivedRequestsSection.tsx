"use client";
import { useCallback, useMemo, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { Request, User } from "@/types";
import SearchBar from "@/components/OtherComponents/SearchBar/SearchBar";
import FilterTabs from "@/components/InputComponents/FilterTabs";
import RequestCard from "@/components/CardComponents/RequestCard/RequestCard";
import NotFound from "@/components/OtherComponents/NotFound/NotFound";
import PaginationFooter from "@/components/FileSectionComponents/FileSection/_components/PaginationFooter";
import { useSearch } from "@/hooks/useSearch";
import { usePagination } from "@/hooks/usePagination";
import "@/components/Pages/ArchivedRequestsComponents/ArchivedRequestsSection.css";

const ITEMS_PER_PAGE = 8;

const CLAIM_FILTERS: { label: string; value: "all" | "claimed" | "unclaimed" }[] = [
        { label: "Toate", value: "all" },
        { label: "Preluate", value: "claimed" },
        { label: "Nepreluate", value: "unclaimed" },
];

type Props = {
        requests: Request[];
        user: User;
};

export default function ArchivedRequestsSection({ requests, user }: Props) {
        const router = useRouter();
        const searchParams = useSearchParams();
        const [currentPage, setCurrentPage] = useState(1);

        const canManage = user.role === "admin" || user.department_id !== null;

        const activeClaim = (searchParams.get("claim") ?? "all") as
                | "all"
                | "claimed"
                | "unclaimed";

        const claimFiltered = useMemo(() => {
                if (activeClaim === "claimed")
                        return requests.filter(
                                (r) => r.claimed_by !== null && r.claimed_by !== undefined,
                        );
                if (activeClaim === "unclaimed") return requests.filter((r) => !r.claimed_by);
                return requests;
        }, [requests, activeClaim]);

        const searchFn = useCallback((req: Request, search: string) => {
                return [
                        req.title,
                        req.description,
                        req.department_name,
                        req.assignee_first_name,
                        req.assignee_last_name,
                        req.assignee_email,
                ]
                        .filter(Boolean)
                        .join(" ")
                        .toLowerCase()
                        .includes(search.toLowerCase());
        }, []);

        const { searchInput, setSearchInput, filteredItems } = useSearch(
                claimFiltered,
                searchFn,
        );
        const { totalPages, paginatedItems } = usePagination(filteredItems, ITEMS_PER_PAGE);

        const handleClaimChange = (claim: "all" | "claimed" | "unclaimed") => {
                const params = new URLSearchParams(searchParams.toString());
                claim === "all" ? params.delete("claim") : params.set("claim", claim);
                router.push(`?${params.toString()}`);
                setCurrentPage(1);
        };

        if (requests.length === 0) {
                return (
                        <NotFound
                                text="Nu există dosare arhivate."
                                subtext="Dosarele arhivate vor apărea aici."
                        />
                );
        }

        return (
                <div className="section">
                        <div className="requests-filters">
                                <SearchBar
                                        value={searchInput}
                                        onChange={(value) => {
                                                setSearchInput(value);
                                                setCurrentPage(1);
                                        }}
                                        placeholder="Caută dosar arhivat..."
                                />
                                {canManage && (
                                        <FilterTabs
                                                tabs={CLAIM_FILTERS.map((f) => ({
                                                        ...f,
                                                        count:
                                                                f.value === "all"
                                                                        ? requests.length
                                                                        : f.value === "claimed"
                                                                          ? requests.filter(
                                                                                    (r) =>
                                                                                            r.claimed_by !==
                                                                                                    null &&
                                                                                            r.claimed_by !==
                                                                                                    undefined,
                                                                            ).length
                                                                          : requests.filter(
                                                                                    (r) =>
                                                                                            !r.claimed_by,
                                                                            ).length,
                                                }))}
                                                active={activeClaim}
                                                onChange={handleClaimChange}
                                        />
                                )}
                        </div>
                        {filteredItems.length === 0 ? (
                                <NotFound
                                        text="Nu am găsit niciun dosar"
                                        subtext="Nu există niciun rezultat care să corespundă căutării tale."
                                />
                        ) : (
                                <>
                                        <div className="objects-grid">
                                                {paginatedItems.map((r) => (
                                                        <RequestCard
                                                                key={r.id}
                                                                user={user}
                                                                searchTerm={searchInput}
                                                                request={r}
                                                                archived={true}
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
