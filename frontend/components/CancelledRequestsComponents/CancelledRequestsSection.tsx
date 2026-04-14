"use client";
import { useState } from "react";
import { Request, User } from "@/types";
import SearchBar from "@/components/OtherComponents/SearchBar/SearchBar";
import RequestCard from "@/components/CardComponents/RequestCard/RequestCard";
import NotFound from "@/components/OtherComponents/NotFound/NotFound";
import "@/components/Pages/ArchivedRequestsComponents/ArchivedRequestsSection.css";

interface Props {
        requests: Request[];
        user: User;
}

export default function CancelledRequestsSection({ requests, user }: Props) {
        const [searchInput, setSearchInput] = useState("");

        const filtered = requests.filter((r) =>
                r.title?.toLowerCase().includes(searchInput.toLowerCase()),
        );

        if (requests.length === 0) {
                return (
                        <NotFound
                                text="Nu există dosare anulate."
                                subtext="Dosarele retrase de utilizatori vor apărea aici."
                                background="white"
                        />
                );
        }

        return (
                <div className="archived-templates">
                        <SearchBar
                                value={searchInput}
                                onChange={setSearchInput}
                                placeholder="Caută dosar retras..."
                        />
                        {filtered.length === 0 ? (
                                <NotFound
                                        text="Nu am găsit niciun dosar"
                                        subtext="Nu există niciun rezultat care să corespundă căutării tale."
                                        background="white"
                                />
                        ) : (
                                <div className="archived-grid">
                                        {filtered.map((r) => (
                                                <RequestCard
                                                        key={r.id}
                                                        user={user}
                                                        searchTerm={searchInput}
                                                        request={r}
                                                        archived={false}
                                                        cancelled={true}
                                                />
                                        ))}
                                </div>
                        )}
                </div>
        );
}
