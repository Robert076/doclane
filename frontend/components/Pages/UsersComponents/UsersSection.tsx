"use client";
import { useState } from "react";
import { User } from "@/types";
import BaseDashboardCard from "@/components/CardComponents/BaseDashboardCard/BaseDashboardCard";
import InfoList from "@/components/CardComponents/InfoList/InfoList";
import InfoItem from "@/components/CardComponents/InfoItem/InfoItem";
import SearchBar from "@/components/OtherComponents/SearchBar/SearchBar";
import NotFound from "@/components/OtherComponents/NotFound/NotFound";
import PaginationFooter from "@/components/FileSectionComponents/FileSection/_components/PaginationFooter";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import Modal from "@/components/Modals/Modal";
import { useSearch } from "@/hooks/useSearch";
import { usePagination } from "@/hooks/usePagination";
import { formatDate } from "@/lib/client/formatDate";
import UserCard from "./UserCard";

interface Props {
        users: User[];
}

const ITEMS_PER_PAGE = 12;

export default function UsersSection({ users }: Props) {
        const [currentPage, setCurrentPage] = useState(1);

        const { searchInput, setSearchInput, filteredItems } = useSearch(
                users,
                (user, search) =>
                        user.email.toLowerCase().includes(search) ||
                        user.first_name.toLowerCase().includes(search) ||
                        user.last_name.toLowerCase().includes(search) ||
                        `${user.first_name} ${user.last_name}`.toLowerCase().includes(search),
        );

        const { totalPages, paginatedItems } = usePagination(filteredItems, ITEMS_PER_PAGE);

        if (users.length === 0) {
                return (
                        <NotFound
                                text="Niciun utilizator înregistrat."
                                subtext="Utilizatorii înregistrați vor apărea aici."
                                background="#fff"
                        />
                );
        }

        return (
                <div className="section">
                        <SearchBar
                                value={searchInput}
                                onChange={(val) => {
                                        setSearchInput(val);
                                        setCurrentPage(1);
                                }}
                                placeholder="Caută utilizator..."
                        />
                        {filteredItems.length === 0 ? (
                                <NotFound
                                        text="Niciun rezultat găsit."
                                        subtext="Încearcă o altă căutare."
                                        background="#fff"
                                />
                        ) : (
                                <>
                                        <div className="objects-grid">
                                                {paginatedItems.map((user) => (
                                                        <UserCard key={user.id} user={user} />
                                                ))}
                                        </div>
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
