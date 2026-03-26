"use client";

import { useState } from "react";

import { Template } from "@/types";
import SearchBar from "@/components/OtherComponents/SearchBar/SearchBar";
import { UI_TEXT } from "@/locales/ro";
import "./ArchivedTemplatesSection.css";
import TemplateCard from "@/components/Pages/TemplatesComponents/TemplateCard";
import NotFound from "@/components/OtherComponents/NotFound/NotFound";

type Props = {
        templates: Template[];
};

const ArchivedTemplatesSection = ({ templates }: Props) => {
        const [searchInput, setSearchInput] = useState("");

        console.log(templates);
        const filteredTemplates = templates.filter((t) => {
                if (t.is_closed === false) return false;
                t.title?.toLowerCase().includes(searchInput.toLowerCase());
                return true;
        });

        if (filteredTemplates.length === 0) {
                return (
                        <NotFound
                                text="Nu ai niciun şablon arhivat."
                                subtext="Aici vor apărea şabloanele pe care le arhivezi."
                                background="white"
                        />
                );
        }

        return (
                <div className="archived-templates">
                        {templates.length > 0 && (
                                <SearchBar
                                        value={searchInput}
                                        onChange={setSearchInput}
                                        placeholder={UI_TEXT.common.search}
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
