import { Template } from "@/types";
import "./TemplateDetailsSummary.css";
import SectionTitle from "@/components/Pages/RequestsComponents/SectionTitle";
import { formatDate } from "@/lib/client/formatDate";

export default function TemplateDetailsSummary({ data }: { data: Template }) {
        console.log(data);
        return (
                <section className="template-details-summary">
                        <SectionTitle text="Detalii șablon" />
                        <div className="info-row" title="Titlul dosarului din șablon">
                                <strong>Titlu</strong>
                                <span>{data.title}</span>
                        </div>
                        <div className="info-row" title="Departamentul șablonului">
                                <strong>Departament:</strong>
                                <span>{data.department_name}</span>
                        </div>
                        <div className="info-row" title="Șablon creat la data">
                                <strong>Creat la:</strong>
                                <span>{formatDate(data.created_at)}</span>
                        </div>
                        {data.description && (
                                <div
                                        className="info-description"
                                        title="Descrierea șablonului"
                                >
                                        <strong>Descriere</strong>
                                        <p>{data.description}</p>
                                </div>
                        )}
                </section>
        );
}
