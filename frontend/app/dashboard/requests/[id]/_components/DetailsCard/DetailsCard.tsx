import { DocumentRequest, RequestStatus } from "@/types";
import "./DetailsCard.css";
import SectionTitle from "@/app/dashboard/requests/[id]/_components/SectionTitle/SectionTitle";
import StatusBadge from "@/components/RequestComponents/StatusBadge/StatusBadge";
import { formatDate } from "@/lib/formatDate";
import { UI_TEXT } from "@/locales/ro";

export default function DetailsCard({ data }: { data: DocumentRequest }) {
        return (
                <section className="details-card">
                        <SectionTitle text={UI_TEXT.request.details.title} />
                        <div className="info-row">
                                <strong>{UI_TEXT.roles.clientSingular}</strong>
                                <span>{data.client_email}</span>
                        </div>
                        <div className="info-row">
                                <strong>{UI_TEXT.common.createdAt}:</strong>
                                <span>{formatDate(data.created_at)}</span>
                        </div>
                        <div className="info-row">
                                <strong>{UI_TEXT.common.status}</strong>
                                <span>
                                        <StatusBadge status={data.status as RequestStatus} />
                                </span>
                        </div>
                        {data.due_date && (
                                <div className="info-row">
                                        <strong>{UI_TEXT.common.dueDate}</strong>
                                        <p>{formatDate(data.due_date)}</p>
                                </div>
                        )}
                        {data.next_due_at && !data.due_date && (
                                <div className="info-row">
                                        <strong>{UI_TEXT.common.nextDueAt}</strong>
                                        <p>{formatDate(data.next_due_at)}</p>
                                </div>
                        )}
                        {data.description && (
                                <div className="info-description">
                                        <strong>{UI_TEXT.common.description}</strong>
                                        <p>{data.description}</p>
                                </div>
                        )}
                </section>
        );
}
