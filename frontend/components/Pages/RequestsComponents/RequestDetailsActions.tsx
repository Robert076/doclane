"use client";
import { useState } from "react";
import SectionTitle from "@/components/Pages/RequestsComponents/SectionTitle";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import toast from "react-hot-toast";
import { UI_TEXT } from "@/locales/ro";
import "./RequestDetailsActions.css";
import { sendEmail } from "@/lib/api/users";
import { addComment } from "@/lib/api/requests"; // you'll add this

export default function RequestDetailsActions({ id }: { id: string }) {
        const [comment, setComment] = useState("");
        const [submitting, setSubmitting] = useState(false);

        const handleAddComment = async () => {
                const trimmed = comment.trim();
                if (trimmed.length < 3) {
                        toast.error("Comentariul trebuie să aibă cel puțin 3 caractere.");
                        return;
                }
                setSubmitting(true);
                const response = await addComment(+id, trimmed);
                setSubmitting(false);
                if (response.success) {
                        toast.success("Comentariu adăugat.");
                        setComment("");
                } else {
                        toast.error(response.message);
                }
        };

        return (
                <aside className="actions-sidebar details-card">
                        <SectionTitle text="Acțiuni" />
                        <div className="action-buttons">
                                <ButtonPrimary
                                        text={UI_TEXT.buttons.sendNotification.normal}
                                        fullWidth={true}
                                        onClick={() => handleNotification(+id)}
                                />
                        </div>
                </aside>
        );
}

const handleNotification = async (id: number) => {
        const response = await sendEmail(id);
        if (response.success === true) {
                toast.success(response.message, { duration: 4000 });
        } else {
                toast.error(response.message, { duration: 4000 });
        }
};
