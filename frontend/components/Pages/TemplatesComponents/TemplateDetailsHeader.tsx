"use client";
import Link from "next/link";
import { Template } from "@/types";
import "./TemplateDetailsHeader.css";

export default function TemplateDetailsHeader({ data }: { data: Template }) {
        return (
                <header className="template-details-header">
                        <Link
                                href="/dashboard/templates"
                                className="back-link"
                                title="Mergi înapoi"
                        >
                                ← Înapoi la Șabloane
                        </Link>
                        <div
                                className="template-details-header-main"
                                title="Titlul șablonului"
                        >
                                <h1>{data.title}</h1>
                        </div>
                </header>
        );
}
