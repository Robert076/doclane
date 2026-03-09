import { useState, useMemo } from "react";

export function usePagination<T>(items: T[], itemsPerPage: number) {
        const [currentPage, setCurrentPage] = useState(1);

        const totalPages = Math.ceil(items.length / itemsPerPage);

        const paginatedItems = useMemo(() => {
                const start = (currentPage - 1) * itemsPerPage;
                return items.slice(start, start + itemsPerPage);
        }, [items, currentPage, itemsPerPage]);

        return {
                currentPage,
                setCurrentPage,
                totalPages,
                paginatedItems,
        };
}
