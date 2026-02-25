"use client";

import "./DetailCardsActionSidebar.css";
import SectionTitle from "@/app/dashboard/requests/[id]/_components/SectionTitle/SectionTitle";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import toast from "react-hot-toast";
import { sendEmail } from "@/lib/api/api";
import { UI_TEXT } from "@/locales/ro";

export default function DetailCardsActionSidebar({ id }: { id: string }) {
        return (
                <aside className="details-card actions-sidebar">
                        <SectionTitle text={UI_TEXT.request.details.actions} />

                        <div className="action-buttons">
                                <ButtonPrimary
                                        text={UI_TEXT.buttons.sendNotification.normal}
                                        fullWidth={true}
                                        onClick={async () => {
                                                const response = await sendEmail(+id);
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
