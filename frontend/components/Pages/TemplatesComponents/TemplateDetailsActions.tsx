"use client";
import "./TemplateDetailsActions.css";
import SectionTitle from "@/components/Pages/RequestsComponents/SectionTitle";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import toast from "react-hot-toast";
import { Template } from "@/types";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { useUser } from "@/context/UserContext";
import EditTemplateModal from "./EditTemplateModal";
import { patchTemplate } from "@/lib/api/templates";
import { createRequest } from "@/lib/api/requests";

export default function TemplateDetailsActions({
        id,
        template,
}: {
        id: number;
        template: Template;
}) {
        const [isEditModalOpen, setIsEditModalOpen] = useState(false);
        const [isSubmitting, setIsSubmitting] = useState(false);
        const router = useRouter();
        const user = useUser();
        const canManage = user.role === "admin" || user.department_id !== null;

        const handleSubmitRequest = async () => {
                setIsSubmitting(true);
                const response = await createRequest({ template_id: id });
                setIsSubmitting(false);
                if (response.success) {
                        toast.success("Cerere depusă cu succes!");
                        router.push("/dashboard/requests");
                } else {
                        toast.error(response.message);
                }
        };

        const handleEditConfirm = async (data: { title?: string; description?: string }) => {
                toast.promise(patchTemplate(id, data), {
                        loading: "Se salvează...",
                        success: (res) => {
                                if (!res.success) throw new Error(res.error);
                                router.refresh();
                                return "Șablon actualizat!";
                        },
                        error: (err) => `Eroare: ${err.message}`,
                });
        };

        return (
                <>
                        <aside className="template-actions">
                                <SectionTitle text="Acțiuni" />
                                <div className="template-action-buttons">
                                        {canManage ? (
                                                <ButtonPrimary
                                                        text="Editează șablon"
                                                        fullWidth
                                                        variant="ghost"
                                                        onClick={() =>
                                                                setIsEditModalOpen(true)
                                                        }
                                                />
                                        ) : (
                                                <ButtonPrimary
                                                        text={
                                                                isSubmitting
                                                                        ? "Se trimite..."
                                                                        : "Depune cerere"
                                                        }
                                                        fullWidth
                                                        disabled={isSubmitting}
                                                        onClick={handleSubmitRequest}
                                                />
                                        )}
                                </div>
                        </aside>
                        <EditTemplateModal
                                isOpen={isEditModalOpen}
                                onClose={() => setIsEditModalOpen(false)}
                                onConfirm={handleEditConfirm}
                                template={template}
                        />
                </>
        );
}
