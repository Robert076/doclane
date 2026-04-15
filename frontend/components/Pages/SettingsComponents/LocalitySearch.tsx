"use client";
import { useState, useEffect, useRef } from "react";
import "./LocalitySearch.css";

interface NominatimResult {
        place_id: number;
        display_name: string;
        address: {
                village?: string;
                town?: string;
                city?: string;
                municipality?: string;
                county?: string;
        };
}

interface Props {
        value: string;
        onChange: (value: string) => void;
}

function extractLocality(result: NominatimResult): string {
        const a = result.address;
        return a.village ?? a.town ?? a.city ?? a.municipality ?? result.display_name;
}

export default function LocalitySearch({ value, onChange }: Props) {
        const [input, setInput] = useState(value);
        const [results, setResults] = useState<NominatimResult[]>([]);
        const [isOpen, setIsOpen] = useState(false);
        const [isLoading, setIsLoading] = useState(false);
        const debounceRef = useRef<NodeJS.Timeout | null>(null);
        const wrapperRef = useRef<HTMLDivElement>(null);

        useEffect(() => {
                setInput(value);
        }, [value]);

        useEffect(() => {
                const handleClickOutside = (e: MouseEvent) => {
                        if (
                                wrapperRef.current &&
                                !wrapperRef.current.contains(e.target as Node)
                        ) {
                                setIsOpen(false);
                        }
                };
                document.addEventListener("mousedown", handleClickOutside);
                return () => document.removeEventListener("mousedown", handleClickOutside);
        }, []);

        const search = async (query: string) => {
                if (query.length < 2) {
                        setResults([]);
                        setIsOpen(false);
                        return;
                }

                setIsLoading(true);
                try {
                        const res = await fetch(
                                `https://nominatim.openstreetmap.org/search?q=${encodeURIComponent(query)}&countrycodes=ro&format=json&addressdetails=1&limit=5`,
                                { headers: { "Accept-Language": "ro" } },
                        );
                        const data: NominatimResult[] = await res.json();
                        setResults(data);
                        setIsOpen(data.length > 0);
                } catch {
                        setResults([]);
                } finally {
                        setIsLoading(false);
                }
        };

        const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
                const val = e.target.value;
                setInput(val);
                onChange(val);

                if (debounceRef.current) clearTimeout(debounceRef.current);
                debounceRef.current = setTimeout(() => search(val), 350);
        };

        const handleSelect = (result: NominatimResult) => {
                const locality = extractLocality(result);
                setInput(locality);
                onChange(locality);
                setIsOpen(false);
                setResults([]);
        };

        return (
                <div className="locality-search" ref={wrapperRef}>
                        <div className="input-wrapper input-wrapper--full-width">
                                <label>Localitate</label>
                                <div className="input-with-icon">
                                        <input
                                                type="text"
                                                value={input}
                                                onChange={handleInputChange}
                                                placeholder="ex: Cluj-Napoca"
                                                className="locality-input"
                                                autoComplete="off"
                                        />
                                </div>
                        </div>
                        {isOpen && (
                                <div className="locality-dropdown">
                                        {isLoading ? (
                                                <div className="locality-dropdown-item locality-dropdown-loading">
                                                        Se caută...
                                                </div>
                                        ) : (
                                                results.map((r) => (
                                                        <div
                                                                key={r.place_id}
                                                                className="locality-dropdown-item"
                                                                onClick={() => handleSelect(r)}
                                                        >
                                                                <span className="locality-dropdown-name">
                                                                        {extractLocality(r)}
                                                                </span>
                                                                {r.address.county && (
                                                                        <span className="locality-dropdown-county">
                                                                                {
                                                                                        r
                                                                                                .address
                                                                                                .county
                                                                                }
                                                                        </span>
                                                                )}
                                                        </div>
                                                ))
                                        )}
                                </div>
                        )}
                </div>
        );
}
