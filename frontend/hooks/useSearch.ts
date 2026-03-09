import { useMemo, useState } from "react";

export function useSearch<T>(items: T[], searchFn: (item: T, search: string) => boolean) {
        const [searchInput, setSearchInput] = useState("");

        const searchLower = searchInput.toLowerCase().trim();

        const filteredItems = useMemo(() => {
                if (!searchLower) return items;
                return items.filter((item) => searchFn(item, searchLower));
        }, [items, searchLower, searchFn]);

        return {
                searchInput,
                setSearchInput,
                filteredItems,
        };
}
