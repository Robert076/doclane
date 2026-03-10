"use client";
import React, { useState } from "react";
import { DocumentRequestTemplate } from "@/types";

import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import HighlightText from "@/components/OtherComponents/HighlightText/HighlightText";
import { useRouter } from "next/navigation";
import { formatDate } from "@/lib/client/formatDate";
import { UI_TEXT } from "@/locales/ro";
import toast from "react-hot-toast";
import BaseDashboardCard from "@/components/CardComponents/BaseDashboardCard/BaseDashboardCard";
import { archiveTemplate, deleteTemplate, unarchiveTemplate } from "@/lib/api/templates";
import DeleteTemplateModal from "./DeleteTemplateModal";

interface TemplateCardProps {
        template: DocumentRequestTemplate;
        searchTerm?: string;
        archived?: boolean;
}

const TemplateCard: React.FC<TemplateCardProps> = ({ template, searchTerm, archived }) => {
        const router = useRouter();
        const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);

        const handleView = () => {
                router.push(`/dashboard/templates/${template.id}`);
        };

        const handleArchive = async () => {
                const loadingToastID = toast.loading(UI_TEXT.common.loading);
                const response = await archiveTemplate(template.id);
                toast.dismiss(loadingToastID);

                if (response.success) {
                        toast.success("Şablon arhivat cu success!");
                } else {
                        toast.error(response.message);
                }
        };

        const handleUnarchive = async () => {
                const loadingToastID = toast.loading(UI_TEXT.common.loading);
                const response = await unarchiveTemplate(template.id);
                toast.dismiss(loadingToastID);

                if (response.success) {
                        toast.success("Şablon restaurat cu success!");
                } else {
                        toast.error(response.message);
                }
        };

        const handleDelete = async () => {
                const loadingToastID = toast.loading(UI_TEXT.common.loading);
                const response = await deleteTemplate(template.id);
                toast.dismiss(loadingToastID);

                if (response.success) {
                        toast.success("Şablon şters definitiv cu success!");
                } else {
                        toast.error(response.message);
                }
        };

        const footer =
                archived === false ? (
                        <>
                                <ButtonPrimary
                                        text="Vezi şablon"
                                        variant="ghost"
                                        fullWidth
                                        onClick={handleView}
                                />
                                <ButtonPrimary
                                        text="Arhivează şablon"
                                        variant="ghost"
                                        fullWidth
                                        onClick={handleArchive}
                                />
                        </>
                ) : archived === true ? (
                        <>
                                <ButtonPrimary
                                        text="Restaurare şablon"
                                        variant="ghost"
                                        fullWidth
                                        onClick={handleUnarchive}
                                />
                                <ButtonPrimary
                                        text="Şterge definitiv"
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
                                <div className="template-info">
                                        <div className="template-info-item">
                                                <span className="template-label">
                                                        Creat pe data
                                                </span>
                                                <span className="template-value">
                                                        {formatDate(template.created_at)}
                                                </span>
                                        </div>

                                        {template.description && (
                                                <div className="template-info-item">
                                                        <span className="template-label">
                                                                {UI_TEXT.common.description}
                                                        </span>
                                                        <span className="template-value">
                                                                <HighlightText
                                                                        text={
                                                                                template.description
                                                                        }
                                                                        search={searchTerm}
                                                                />
                                                        </span>
                                                </div>
                                        )}

                                        {template.is_recurring && template.recurrence_cron && (
                                                <div className="template-info-item">
                                                        <span className="template-label">
                                                                Recurenţă
                                                        </span>
                                                        <span className="template-value">
                                                                {template.recurrence_cron}
                                                        </span>
                                                </div>
                                        )}
                                </div>
                        </BaseDashboardCard>

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
