"use client";
import { useState } from "react";
import { Request, User } from "@/types";
import SearchBar from "@/components/OtherComponents/SearchBar/SearchBar";
import RequestCard from "@/components/CardComponents/RequestCard/RequestCard";
import NotFound from "@/components/OtherComponents/NotFound/NotFound";

type Props = {
        requests: Request[];
        user: User;
};

export default function ArchivedRequestsSection({ requests, user }: Props) {
        const [searchInput, setSearchInput] = useState("");

        const filteredRequests = requests.filter((r) =>
                r.title?.toLowerCase().includes(searchInput.toLowerCase()),
        );

        if (requests.length === 0) {
                return (
                        <NotFound
                                text="Nu există dosare arhivate."
                                subtext="Dosarele arhivate vor apărea aici."
                                background="white"
                        />
                );
        }

        return (
                <div className="archived-templates">
                        <SearchBar
                                value={searchInput}
                                onChange={setSearchInput}
                                placeholder="Caută dosar arhivat..."
                        />
                        {filteredRequests.length === 0 ? (
                                <NotFound
                                        text="Nu am găsit niciun dosar"
                                        subtext="Nu există niciun rezultat care să corespundă căutării tale."
                                        background="white"
                                />
                        ) : (
                                <div className="archived-grid">
                                        {filteredRequests.map((r) => (
                                                <RequestCard
                                                        key={r.id}
                                                        user={user}
                                                        searchTerm={searchInput}
                                                        request={r}
                                                        archived={true}
                                                />
                                        ))}
                                </div>
                        )}
                </div>
        );
}
