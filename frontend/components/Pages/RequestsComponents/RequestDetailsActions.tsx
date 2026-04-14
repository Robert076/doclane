"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import SectionTitle from "@/components/Pages/RequestsComponents/SectionTitle";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import Modal from "@/components/Modals/Modal";
import toast from "react-hot-toast";
import { notifyUser } from "@/lib/api/users";
import { cancelRequest } from "@/lib/api/requests";
import { useUser } from "@/context/UserContext";
import { Request } from "@/types";
import "./RequestDetailsActions.css";

interface Props {
        assignee: number;
        request: Request;
}

export default function RequestDetailsActions({ assignee, request }: Props) {
        const user = useUser();
        const router = useRouter();
        const [isCancelModalOpen, setIsCancelModalOpen] = useState(false);

        const canManage = user.role === "admin" || user.department_id !== null;

        const handleNotify = async () => {
                const response = await notifyUser(assignee);
                response.success
                        ? toast.success(response.message, { duration: 4000 })
                        : toast.error(response.message, { duration: 4000 });
        };

        const handleCancel = async () => {
                const response = await cancelRequest(request.id);
                if (response.success) {
                        toast.success("Cerere retrasă cu succes.");
                        router.push("/dashboard/requests");
                } else {
                        toast.error(response.message);
                }
        };

        if (request.is_cancelled) {
                return null;
        }

        return (
                <aside className="actions-sidebar details-card">
                        <SectionTitle text="Acțiuni" />
                        <div className="action-buttons">
                                {canManage && (
                                        <ButtonPrimary
                                                text="Trimite notificare"
                                                fullWidth
                                                onClick={handleNotify}
                                        />
                                )}
                                <ButtonPrimary
                                        text="Retrage cerere"
                                        fullWidth
                                        onClick={() => setIsCancelModalOpen(true)}
                                />
                        </div>

                        <Modal
                                isOpen={isCancelModalOpen}
                                onClose={() => setIsCancelModalOpen(false)}
                                onConfirm={handleCancel}
                                title="Retrage cererea"
                        >
                                <p className="modal-text">
                                        Ești sigur că vrei să retragi cererea{" "}
                                        <strong>{request.title}</strong>?
                                </p>
                                <p className="modal-subtext">
                                        Această acțiune este ireversibilă. Cererea va fi
                                        marcată ca retrasă și nu vei mai putea adăuga
                                        documente.
                                </p>
                        </Modal>
                </aside>
        );
}
