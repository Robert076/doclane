"use client";

import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import UploadDocumentButton from "@/components/Buttons/UploadDocumentButton/UploadDocumentButton";
import "./DetailCardsActionSidebar.css";
import SectionTitle from "@/components/SectionTitle/SectionTitle";
import toast from "react-hot-toast";
import { useRouter } from "next/navigation";

export default function DetailCardsActionSidebar({ id }: { id: string }) {
        const router = useRouter();

        const handleMarkCompleted = async () => {
                const updateStatusPromise = fetch(
                        `/api/backend/document-requests/${id}/status`,
                        {
                                method: "PUT",
                                credentials: "include",
                                headers: {
                                        "Content-Type": "application/json",
                                },
                                body: JSON.stringify({
                                        status: "uploaded",
                                }),
                        },
                ).then(async (res) => {
                        if (!res.ok) {
                                const errorData = await res.json();
                                throw new Error(errorData.error || "Failed to update status");
                        }

                        return res.json();
                });

                toast.promise(updateStatusPromise, {
                        loading: "Updating status...",
                        success: "Request marked as completed",
                        error: (err) => `Failed: ${err.message}`,
                });

                updateStatusPromise.then(() => {
                        router.refresh();
                });
        };

        return (
                <aside className="details-card actions-sidebar">
                        <SectionTitle text="Actions" />

                        <div className="action-buttons">
                                <UploadDocumentButton requestId={id} />

                                {/* <ButtonPrimary
          text="Mark as Completed"
          variant="secondary"
          fullWidth
          onClick={handleMarkCompleted}
        /> */}
                        </div>
                </aside>
        );
}
