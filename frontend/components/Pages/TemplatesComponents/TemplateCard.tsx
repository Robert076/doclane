"use client";
import React, { useState } from "react";
import { Template } from "@/types";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import HighlightText from "@/components/OtherComponents/HighlightText/HighlightText";
import { useRouter } from "next/navigation";
import { formatDate } from "@/lib/client/formatDate";
import toast from "react-hot-toast";
import BaseDashboardCard from "@/components/CardComponents/BaseDashboardCard/BaseDashboardCard";
import InfoList from "@/components/CardComponents/InfoList/InfoList";
import InfoItem from "@/components/CardComponents/InfoItem/InfoItem";
import { archiveTemplate, deleteTemplate, unarchiveTemplate } from "@/lib/api/templates";
import { createRequest } from "@/lib/api/requests";
import DeleteTemplateModal from "./DeleteTemplateModal";
import Modal from "@/components/Modals/Modal";
import { useUser } from "@/context/UserContext";

interface TemplateCardProps {
        template: Template;
        searchTerm?: string;
        archived?: boolean;
}

const TemplateCard: React.FC<TemplateCardProps> = ({ template, searchTerm, archived }) => {
        const router = useRouter();
        const user = useUser();
        const canManage = user.role === "admin" || user.department_id !== null;

        const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
        const [isArchiveModalOpen, setIsArchiveModalOpen] = useState(false);
        const [isSubmitting, setIsSubmitting] = useState(false);

        const handleView = () => router.push(`/dashboard/templates/${template.id}`);

        const handleSubmitRequest = async () => {
                setIsSubmitting(true);
                const response = await createRequest({ template_id: template.id });
                setIsSubmitting(false);
                if (response.success) {
                        toast.success("Cerere depusă cu succes!");
                        router.push("/dashboard/requests");
                } else {
                        toast.error(response.message);
                }
        };

        const handleArchive = async () => {
                const id = toast.loading("Se arhivează...");
                const response = await archiveTemplate(template.id);
                toast.dismiss(id);
                response.success
                        ? toast.success("Șablon arhivat cu succes!")
                        : toast.error(response.message);
        };

        const handleUnarchive = async () => {
                const id = toast.loading("Se restaurează...");
                const response = await unarchiveTemplate(template.id);
                toast.dismiss(id);
                response.success
                        ? toast.success("Șablon restaurat cu succes!")
                        : toast.error(response.message);
        };

        const handleDelete = async () => {
                const id = toast.loading("Se șterge...");
                const response = await deleteTemplate(template.id);
                toast.dismiss(id);
                response.success
                        ? toast.success("Șablon șters definitiv cu succes!")
                        : toast.error(response.message);
        };

        const footer =
                archived === false ? (
                        <>
                                {canManage && (
                                        <ButtonPrimary
                                                text="Vezi șablon"
                                                variant="ghost"
                                                fullWidth
                                                onClick={handleView}
                                        />
                                )}
                                {!canManage && (
                                        <ButtonPrimary
                                                text={
                                                        isSubmitting
                                                                ? "Se trimite..."
                                                                : "Depune cerere"
                                                }
                                                variant="ghost"
                                                fullWidth
                                                disabled={isSubmitting}
                                                onClick={handleSubmitRequest}
                                        />
                                )}
                                {canManage && (
                                        <ButtonPrimary
                                                text="Arhivează șablon"
                                                variant="ghost"
                                                fullWidth
                                                onClick={() => setIsArchiveModalOpen(true)}
                                        />
                                )}
                        </>
                ) : archived === true ? (
                        <>
                                <ButtonPrimary
                                        text="Restaurare șablon"
                                        variant="ghost"
                                        fullWidth
                                        onClick={handleUnarchive}
                                />
                                <ButtonPrimary
                                        text="Șterge definitiv"
                                        variant="ghost"
                                        fullWidth
                                        onClick={() => setIsDeleteModalOpen(true)}
                                />
                        </>
                ) : null;

        return (
                <>
                        <BaseDashboardCard
                                title={
                                        <HighlightText
                                                text={template.title}
                                                search={searchTerm}
                                        />
                                }
                                footer={footer}
                        >
                                <InfoList>
                                        <InfoItem
                                                label="Departament"
                                                value={
                                                        <HighlightText
                                                                text={template.department_name}
                                                                search={searchTerm}
                                                        />
                                                }
                                        />
                                        <InfoItem
                                                label="Creat pe data"
                                                value={formatDate(template.created_at)}
                                        />
                                        {template.description && (
                                                <InfoItem
                                                        label="Descriere"
                                                        value={
                                                                <HighlightText
                                                                        text={
                                                                                template.description
                                                                        }
                                                                        search={searchTerm}
                                                                />
                                                        }
                                                />
                                        )}
                                </InfoList>
                        </BaseDashboardCard>

                        <Modal
                                isOpen={isArchiveModalOpen}
                                onClose={() => setIsArchiveModalOpen(false)}
                                onConfirm={handleArchive}
                                title="Arhivează șablon"
                        >
                                <p>
                                        Ești sigur că vrei să arhivezi șablonul{" "}
                                        <strong>„{template.title}"</strong>? Îl vei putea
                                        restaura ulterior din secțiunea de arhivă.
                                </p>
                        </Modal>

                        <DeleteTemplateModal
                                isOpen={isDeleteModalOpen}
                                onClose={() => setIsDeleteModalOpen(false)}
                                onConfirm={handleDelete}
                                templateTitle={template.title}
                        />
                </>
        );
};

export default TemplateCard;
