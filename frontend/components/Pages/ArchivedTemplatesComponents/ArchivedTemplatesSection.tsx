"use client";
import { useState } from "react";
import { Template } from "@/types";
import SearchBar from "@/components/OtherComponents/SearchBar/SearchBar";
import TemplateCard from "@/components/Pages/TemplatesComponents/TemplateCard";
import NotFound from "@/components/OtherComponents/NotFound/NotFound";
import "./ArchivedTemplatesSection.css";

type Props = {
        templates: Template[];
};

export default function ArchivedTemplatesSection({ templates }: Props) {
        const [searchInput, setSearchInput] = useState("");

        const filteredTemplates = templates.filter((t) =>
                t.title?.toLowerCase().includes(searchInput.toLowerCase()),
        );

        if (templates.length === 0) {
                return (
                        <NotFound
                                text="Nu există șabloane arhivate."
                                subtext="Șabloanele arhivate vor apărea aici."
                                background="white"
                        />
                );
        }

        return (
                <div className="archived-templates">
                        <SearchBar
                                value={searchInput}
                                onChange={setSearchInput}
                                placeholder="Caută șablon arhivat..."
                        />
                        {filteredTemplates.length === 0 ? (
                                <NotFound
                                        text="Nu am găsit niciun șablon"
                                        subtext="Nu există niciun rezultat care să corespundă căutării tale."
                                        background="white"
                                />
                        ) : (
                                <div className="archived-grid">
                                        {filteredTemplates.map((t) => (
                                                <TemplateCard
                                                        key={t.id}
                                                        searchTerm={searchInput}
                                                        template={t}
                                                        archived={true}
                                                />
                                        ))}
                                </div>
                        )}
                </div>
        );
}
