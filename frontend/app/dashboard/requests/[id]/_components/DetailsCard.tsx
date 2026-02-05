import { DocumentRequest, RequestStatus } from "@/types";
import "./DetailsCard.css";
import SectionTitle from "@/components/SectionTitle/SectionTitle";
import StatusBadge from "@/components/Request/StatusBadge/StatusBadge";
import { formatDate } from "@/lib/formatDate";

export default function DetailsCard({ data }: { data: DocumentRequest }) {
  return (
    <section className="details-card">
      <SectionTitle text="Request details" />
      <div className="info-row">
        <strong>Client:</strong>
        <span>{data.client_email}</span>
      </div>
      <div className="info-row">
        <strong>Created At:</strong>
        <span>{formatDate(data.created_at)}</span>
      </div>
      <div className="info-row">
        <strong>Status:</strong>
        <span>
          <StatusBadge status={data.status as RequestStatus} />
        </span>
      </div>
      {data.due_date && (
        <div className="info-row">
          <strong>Due date:</strong>
          <p>{formatDate(data.due_date)}</p>
        </div>
      )}
      {data.description && (
        <div className="info-description">
          <strong>Description:</strong>
          <p>{data.description}</p>
        </div>
      )}
    </section>
  );
}
