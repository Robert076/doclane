"use client";
import { useState, useEffect } from "react";
import { DocumentRequest, User } from "@/types";
import Request from "../Request/Request";
import NotFound from "@/components/NotFound/NotFound";
import "./RequestsSection.css";
import SearchBar from "../SearchBar/SearchBar";
import PaginationFooter from "../ClientComponents/ClientsSection/_components/PaginationFooter";

interface RequestsSectionProps {
        requests: DocumentRequest[];
        user: User;
}

const ITEMS_PER_PAGE = 8;

const RequestsSection: React.FC<RequestsSectionProps> = ({ requests, user }) => {
        const [currentPage, setCurrentPage] = useState<number>(1);
        const [searchInput, setSearchInput] = useState<string>("");

        const filteredRequests = requests.filter((req) => {
                if (!searchInput) return true;
                const searchLower = searchInput.toLowerCase().trim();

                const fullName =
                        `${req.client_first_name || ""} ${req.client_last_name || ""}`.toLowerCase();

                return (
                        req.title?.toLowerCase().includes(searchLower) ||
                        req.description?.toLowerCase().includes(searchLower) ||
                        req.client_first_name?.toLowerCase().includes(searchLower) ||
                        req.client_last_name?.toLowerCase().includes(searchLower) ||
                        req.client_email?.toLowerCase().includes(searchLower) ||
                        fullName.includes(searchLower) ||
                        req.status?.toLowerCase().includes(searchLower)
                );
        });

        useEffect(() => {
                setCurrentPage(1);
        }, [searchInput]);

        if (requests.length === 0) {
                return user.role === "PROFESSIONAL" ? (
                        <NotFound
                                text="No document requests found."
                                subtext="Create a request to get started."
                                background="#fff"
                        />
                ) : (
                        <NotFound
                                text="No document requests found."
                                subtext="Seems like it's your lucky day."
                                background="#fff"
                        />
                );
        }

        const totalPages = Math.ceil(filteredRequests.length / ITEMS_PER_PAGE);
        const startIndex = (currentPage - 1) * ITEMS_PER_PAGE;
        const endIndex = startIndex + ITEMS_PER_PAGE;
        const currentRequests = filteredRequests.slice(startIndex, endIndex);

        return (
                <div className="requests-section">
                        <SearchBar
                                value={searchInput}
                                onChange={setSearchInput}
                                placeholder="Search requests..."
                        />

                        {filteredRequests.length === 0 ? (
                                <NotFound
                                        text="No requests found."
                                        subtext="Try adjusting your search terms."
                                        background="#fff"
                                />
                        ) : (
                                <>
                                        <div className="requests-grid">
                                                {currentRequests.map((req) => (
                                                        <Request
                                                                key={req.id}
                                                                request={req}
                                                                user={user}
                                                                searchTerm={searchInput}
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
};

export default RequestsSection;
