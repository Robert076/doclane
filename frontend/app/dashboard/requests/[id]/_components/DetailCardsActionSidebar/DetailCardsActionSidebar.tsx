"use client";

import "./DetailCardsActionSidebar.css";
import SectionTitle from "@/app/dashboard/requests/[id]/_components/SectionTitle/SectionTitle";
import { useRouter } from "next/navigation";
import UploadDocumentButton from "@/app/dashboard/requests/[id]/_components/UploadDocumentButton/UploadDocumentButton";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import { sendEmail } from "@/lib/api/api";
import toast from "react-hot-toast";

export default function DetailCardsActionSidebar({ id }: { id: string }) {
        const router = useRouter();

        return (
                <aside className="details-card actions-sidebar">
                        <SectionTitle text="Actions" />

                        <div className="action-buttons">
                                <UploadDocumentButton requestId={id} />
                                <ButtonPrimary
                                        text="Send Notification"
                                        fullWidth={true}
                                        onClick={async () => {
                                                const response = await sendEmail(id);
                                                if (response.success === true) {
                                                        toast.success(response.message, {
                                                                duration: 4000,
                                                        });
                                                } else {
                                                        toast.error(response.message, {
                                                                duration: 4000,
                                                        });
                                                }
                                        }}
                                />
                        </div>
                </aside>
        );
}
