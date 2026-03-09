"use client";
import "./TemplateDetailsActions.css";
import SectionTitle from "@/components/Pages/RequestsComponents/SectionTitle";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import toast from "react-hot-toast";
import { User } from "@/types";
import { useState } from "react";
import { useRouter } from "next/navigation";
import AssignClientModal from "./AssignClientModal";
import { instantiateTemplate } from "@/lib/api/templates";

export default function TemplateActions({ id, clients }: { id: string; clients: User[] }) {
        const [isModalOpen, setIsModalOpen] = useState(false);
        const router = useRouter();

        const handleConfirm = async (client: User) => {
                toast.promise(
                        instantiateTemplate(+id, {
                                client_id: +client.id,
                                is_scheduled: false,
                        }),
                        {
                                loading: "Se generează dosarul...",
                                success: (res) => {
                                        if (!res.success) throw new Error(res.error);
                                        router.push("/dashboard");
                                        return "Dosar generat cu succes!";
                                },
                                error: (err) => `Eroare: ${err.message}`,
                        },
                );
        };

        return (
                <>
                        <aside className="template-actions">
                                <SectionTitle text="Acțiuni" />
                                <div className="template-action-buttons">
                                        <ButtonPrimary
                                                text="Generează dosar"
                                                fullWidth={true}
                                                onClick={() => setIsModalOpen(true)}
                                        />
                                </div>
                        </aside>
                        <AssignClientModal
                                isOpen={isModalOpen}
                                onClose={() => setIsModalOpen(false)}
                                onConfirm={handleConfirm}
                                clients={clients}
                        />
                </>
        );
}
