"use client";
import { useMemo, useCallback, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { Request, User, RequestStatus } from "@/types";
import RequestCard from "@/components/CardComponents/RequestCard/RequestCard";
import NotFound from "@/components/OtherComponents/NotFound/NotFound";
import SearchBar from "@/components/OtherComponents/SearchBar/SearchBar";
import { useSearch } from "@/hooks/useSearch";
import { usePagination } from "@/hooks/usePagination";
import PaginationFooter from "@/components/FileSectionComponents/FileSection/_components/PaginationFooter";
import FilterTabs from "@/components/InputComponents/FilterTabs";
import "./RequestsSection.css";

interface RequestsSectionProps {
        requests: Request[];
        user: User;
}

const ITEMS_PER_PAGE = 8;

const STATUS_FILTERS: { label: string; value: RequestStatus | "all" }[] = [
        { label: "Orice status", value: "all" },
        { label: "In asteptare", value: "pending" },
        { label: "Incarcat", value: "uploaded" },
        { label: "Expirat", value: "overdue" },
];

const CLAIM_FILTERS: { label: string; value: "all" | "claimed" | "unclaimed" }[] = [
        { label: "Toate", value: "all" },
        { label: "Preluate", value: "claimed" },
        { label: "Nepreluate", value: "unclaimed" },
];

const STALE_DAYS = 7;

function isStale(req: Request): boolean {
        if (req.claimed_by !== null && req.claimed_by !== undefined) return false;
        const created = new Date(req.created_at);
        const diffMs = Date.now() - created.getTime();
        return diffMs > STALE_DAYS * 24 * 60 * 60 * 1000;
}

function isDueSoon(req: Request): boolean {
        if (!req.due_date) return false;
        if (req.is_closed || req.is_cancelled) return false;
        const diffMs = new Date(req.due_date).getTime() - Date.now();
        const diffDays = diffMs / (1000 * 60 * 60 * 24);
        return diffDays >= 0 && diffDays <= 3;
}

export default function RequestsSection({ requests, user }: RequestsSectionProps) {
        const router = useRouter();
        const searchParams = useSearchParams();
        const [currentPage, setCurrentPage] = useState(1);

        const activeStatus = (searchParams.get("status") ?? "all") as RequestStatus | "all";
        const activeClaim = (searchParams.get("claim") ?? "all") as
                | "all"
                | "claimed"
                | "unclaimed";

        const canManage = user.role === "admin" || user.department_id !== null;

        const openRequests = useMemo(
                () => requests.filter((r) => !r.is_closed && !r.is_cancelled),
                [requests],
        );

        const statusFiltered = useMemo(
                () =>
                        activeStatus === "all"
                                ? openRequests
                                : openRequests.filter((r) => r.status === activeStatus),
                [openRequests, activeStatus],
        );

        const claimFiltered = useMemo(() => {
                if (activeClaim === "claimed")
                        return statusFiltered.filter(
                                (r) => r.claimed_by !== null && r.claimed_by !== undefined,
                        );
                if (activeClaim === "unclaimed")
                        return statusFiltered.filter((r) => !r.claimed_by);
                return statusFiltered;
        }, [statusFiltered, activeClaim]);

        const searchFn = useCallback((req: Request, search: string) => {
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
                        .includes(search.toLowerCase());
        }, []);

        const { searchInput, setSearchInput, filteredItems } = useSearch(
                claimFiltered,
                searchFn,
        );
        const { totalPages, paginatedItems } = usePagination(filteredItems, ITEMS_PER_PAGE);

        const handleStatusChange = (status: RequestStatus | "all") => {
                const params = new URLSearchParams(searchParams.toString());
                status === "all" ? params.delete("status") : params.set("status", status);
                router.push(`?${params.toString()}`);
                setCurrentPage(1);
        };

        const handleClaimChange = (claim: "all" | "claimed" | "unclaimed") => {
                const params = new URLSearchParams(searchParams.toString());
                claim === "all" ? params.delete("claim") : params.set("claim", claim);
                router.push(`?${params.toString()}`);
                setCurrentPage(1);
        };

        if (requests.length === 0) {
                return (
                        <NotFound
                                text="Nu ai niciun dosar inca."
                                subtext={
                                        user.role === "admin" || user.department_id !== null
                                                ? "Aici vor aparea dosarele pe care le gestionezi."
                                                : "Aici vor aparea dosarele tale. Acceseaza sabloanele pentru a deschide un dosar."
                                }
                                background="#fff"
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
                                        placeholder="Cauta..."
                                />
                                <FilterTabs
                                        tabs={STATUS_FILTERS.map((f) => ({
                                                ...f,
                                                count:
                                                        f.value === "all"
                                                                ? openRequests.length
                                                                : openRequests.filter(
                                                                          (r) =>
                                                                                  r.status ===
                                                                                  f.value,
                                                                  ).length,
                                        }))}
                                        active={activeStatus}
                                        onChange={handleStatusChange}
                                />
                                {canManage && (
                                        <FilterTabs
                                                tabs={CLAIM_FILTERS.map((f) => ({
                                                        ...f,
                                                        count:
                                                                f.value === "all"
                                                                        ? openRequests.length
                                                                        : f.value === "claimed"
                                                                          ? openRequests.filter(
                                                                                    (r) =>
                                                                                            r.claimed_by !==
                                                                                                    null &&
                                                                                            r.claimed_by !==
                                                                                                    undefined,
                                                                            ).length
                                                                          : openRequests.filter(
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
                                        text="Nu am gasit niciun dosar"
                                        subtext="Nu exista niciun rezultat care sa corespunda cautarii tale."
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
                                                                isStale={isStale(req)}
                                                                isDueSoon={isDueSoon(req)}
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
