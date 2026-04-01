"use client";
import SectionTitle from "@/components/Pages/RequestsComponents/SectionTitle";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import toast from "react-hot-toast";
import { UI_TEXT } from "@/locales/ro";
import { notifyUser } from "@/lib/api/users";
import { useUser } from "@/context/UserContext";
import "./RequestDetailsActions.css";

export default function RequestDetailsActions({ assignee }: { assignee: number }) {
        const user = useUser();

        const handleNotification = async () => {
                const response = await notifyUser(assignee);
                if (response.success) {
                        toast.success(response.message, { duration: 4000 });
                } else {
                        toast.error(response.message, { duration: 4000 });
                }
        };

        return (
                <aside className="actions-sidebar details-card">
                        <SectionTitle text="Acțiuni" />
                        <div className="action-buttons">
                                {(user.role === "admin" || user.department_id !== null) && (
                                        <ButtonPrimary
                                                text={UI_TEXT.buttons.sendNotification.normal}
                                                fullWidth={true}
                                                onClick={handleNotification}
                                        />
                                )}
                        </div>
                </aside>
        );
}
