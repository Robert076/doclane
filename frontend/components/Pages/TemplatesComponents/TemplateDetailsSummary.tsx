import { DocumentRequestTemplate } from "@/types";
import "./TemplateDetailsSummary.css";
import SectionTitle from "@/components/Pages/RequestsComponents/SectionTitle";
import { formatDate } from "@/lib/formatDate";
import { UI_TEXT } from "@/locales/ro";

export default function DetailsCard({ data }: { data: DocumentRequestTemplate }) {
        return (
                <section className="template-details-summary">
                        <SectionTitle text={UI_TEXT.request.details.title} />
                        <div className="info-row" title="Titlul dosarului din şablon">
                                <strong>Titlu dosar</strong>
                                <span>{data.title}</span>
                        </div>
                        <div className="info-row" title="Şablon creat la data">
                                <strong>{UI_TEXT.common.createdAt}:</strong>
                                <span>{formatDate(data.created_at)}</span>
                        </div>
                        {data.description && (
                                <div
                                        className="info-description"
                                        title="Descrierea şablonului"
                                >
                                        <strong>{UI_TEXT.common.description}</strong>
                                        <p>{data.description}</p>
                                </div>
                        )}
                </section>
        );
}
