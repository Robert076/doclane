"use client";

import { useState } from "react";

import { DocumentRequestTemplate } from "@/types";
import SearchBar from "@/components/OtherComponents/SearchBar/SearchBar";
import { UI_TEXT } from "@/locales/ro";
import "./ArchivedTemplatesSection.css";
import TemplateCard from "@/components/Pages/TemplatesComponents/TemplateCard";
import NotFound from "@/components/OtherComponents/NotFound/NotFound";

type Props = {
        templates: DocumentRequestTemplate[];
};

const ArchivedTemplatesSection = ({ templates }: Props) => {
        const [searchInput, setSearchInput] = useState("");

        console.log(templates);
        const filteredTemplates = templates.filter((t) => {
                if (t.is_closed === false) return false;
                t.title?.toLowerCase().includes(searchInput.toLowerCase());
                return true;
        });

        console.log(filteredTemplates);
        return (
                <div className="archived-templates">
                        <SearchBar
                                value={searchInput}
                                onChange={setSearchInput}
                                placeholder={UI_TEXT.common.search}
                        />

                        {filteredTemplates.length === 0 && (
                                <NotFound
                                        text="Nu ai niciun şablon arhivat."
                                        background="white"
                                />
                        )}
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
                </div>
        );
};

export default ArchivedTemplatesSection;
