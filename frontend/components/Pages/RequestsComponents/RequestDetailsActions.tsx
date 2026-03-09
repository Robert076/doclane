"use client";

import SectionTitle from "@/components/Pages/RequestsComponents/SectionTitle";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import toast from "react-hot-toast";
import { UI_TEXT } from "@/locales/ro";
import "./RequestDetailsActions.css";
import { sendEmail } from "@/lib/api/users";

export default function RequestDetailsActions({ id }: { id: string }) {
        return (
                <aside className="actions-sidebar details-card">
                        <SectionTitle text="Acțiuni" />

                        <div className="action-buttons">
                                <ButtonPrimary
                                        text={UI_TEXT.buttons.sendNotification.normal}
                                        fullWidth={true}
                                        onClick={() => {
                                                handleNotification(+id);
                                        }}
                                />
                        </div>
                </aside>
        );
}

const handleNotification = async (id: number) => {
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
};
