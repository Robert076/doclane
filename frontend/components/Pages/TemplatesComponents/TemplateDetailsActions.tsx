"use client";
import "./TemplateDetailsActions.css";
import SectionTitle from "@/components/Pages/RequestsComponents/SectionTitle";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import toast from "react-hot-toast";
import { Template, User } from "@/types";
import { useState } from "react";
import { useRouter } from "next/navigation";
import AssignClientModal from "./AssignClientModal";
import EditTemplateModal from "./EditTemplateModal";
import { instantiateTemplate, patchTemplate } from "@/lib/api/templates";

export default function TemplateActions({
        id,
        clients,
        template,
}: {
        id: string;
        clients: User[];
        template: Template;
}) {
        const [isInstantiateModalOpen, setIsInstantiateModalOpen] = useState(false);
        const [isEditModalOpen, setIsEditModalOpen] = useState(false);
        const router = useRouter();

        const handleInstantiateConfirm = async (client: User) => {
                toast.promise(
                        instantiateTemplate(+id, {
                                client_id: +client.id,
                                is_scheduled: false,
                        }),
                        {
                                loading: "Se generează dosarul...",
                                success: (res) => {
                                        if (!res.success) throw new Error(res.error);
                                        router.push("/dashboard/templates");
                                        return "Dosar generat cu succes!";
                                },
                                error: (err) => `Eroare: ${err.message}`,
                        },
                );
        };

        const handleEditConfirm = async (data: { title?: string; description?: string }) => {
                toast.promise(patchTemplate(+id, data), {
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
                                        <ButtonPrimary
                                                text="Generează dosar"
                                                fullWidth={true}
                                                onClick={() => setIsInstantiateModalOpen(true)}
                                        />
                                        <ButtonPrimary
                                                text="Editează șablon"
                                                fullWidth={true}
                                                variant="ghost"
                                                onClick={() => setIsEditModalOpen(true)}
                                        />
                                </div>
                        </aside>

                        <AssignClientModal
                                isOpen={isInstantiateModalOpen}
                                onClose={() => setIsInstantiateModalOpen(false)}
                                onConfirm={handleInstantiateConfirm}
                                clients={clients}
                        />
                        <EditTemplateModal
                                isOpen={isEditModalOpen}
                                onClose={() => setIsEditModalOpen(false)}
                                onConfirm={handleEditConfirm}
                                template={template}
                        />
                </>
        );
}
