"use client";

import { useState } from "react";

import { Request } from "@/types";
import SearchBar from "@/components/OtherComponents/SearchBar/SearchBar";
import { UI_TEXT } from "@/locales/ro";
import "./ArchivedRequestsSection.css";
import RequestCard from "@/components/CardComponents/RequestCard/RequestCard";
import { useUser } from "@/context/UserContext";
import NotFound from "@/components/OtherComponents/NotFound/NotFound";

type Props = {
        requests: Request[];
};

const ArchivedRequestsSection = ({ requests }: Props) => {
        const [searchInput, setSearchInput] = useState("");
        const user = useUser();

        const filteredRequests = requests.filter((r) => {
                if (r.is_closed === false) return false;
                r.title?.toLowerCase().includes(searchInput.toLowerCase());
                return true;
        });

        if (filteredRequests.length === 0) {
                return (
                        <NotFound
                                text="Nu ai niciun dosar arhivat încă."
                                subtext="Aici vor apărea dosarele pe care le arhivezi."
                                background="white"
                        />
                );
        }

        return (
                <div className="archived-templates">
                        {requests.length > 0 && (
                                <SearchBar
                                        value={searchInput}
                                        onChange={setSearchInput}
                                        placeholder={UI_TEXT.common.search}
                                />
                        )}
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
                </div>
        );
};

export default ArchivedRequestsSection;
