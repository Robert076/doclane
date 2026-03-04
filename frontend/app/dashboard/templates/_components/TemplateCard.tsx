"use client";
import React from "react";
import { DocumentRequestTemplate } from "@/types";
import "./TemplateCard.css";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import HighlightText from "@/components/OtherComponents/HighlightText/HighlightText";
import { useRouter } from "next/navigation";
import { formatDate } from "@/lib/client/formatDate";
import { UI_TEXT } from "@/locales/ro";
import { archiveTemplate } from "@/lib/api/api";
import toast from "react-hot-toast";

interface TemplateCardProps {
        template: DocumentRequestTemplate;
        searchTerm?: string;
}

const TemplateCard: React.FC<TemplateCardProps> = ({ template, searchTerm }) => {
        const router = useRouter();

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

        return (
                <div className="template-card">
                        <h3 className="template-name">
                                <HighlightText text={template.title} search={searchTerm} />
                        </h3>
                        <div className="template-body">
                                <div className="template-info">
                                        <div className="template-info-item">
                                                <span className="template-label">
                                                        {UI_TEXT.common.createdAt}
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
                        </div>
                        <div className="template-footer">
                                <ButtonPrimary
                                        text="Vezi şablon"
                                        variant="ghost"
                                        fullWidth={true}
                                        onClick={handleView}
                                />
                                <ButtonPrimary
                                        text="Arhivează şablon"
                                        variant="ghost"
                                        fullWidth={true}
                                        onClick={handleArchive}
                                />
                        </div>
                </div>
        );
};

export default TemplateCard;
