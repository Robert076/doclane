"use client";

import "./DetailCardsActionSidebar.css";
import SectionTitle from "@/app/dashboard/requests/[id]/_components/SectionTitle/SectionTitle";
import { useRouter } from "next/navigation";
import UploadDocumentButton from "@/components/ButtonComponents/UploadDocumentButton/UploadDocumentButton";

export default function DetailCardsActionSidebar({ id }: { id: string }) {
        const router = useRouter();

        return (
                <aside className="details-card actions-sidebar">
                        <SectionTitle text="Actions" />

                        <div className="action-buttons">
                                <UploadDocumentButton requestId={id} />
                        </div>
                </aside>
        );
}
