import { DocumentRequest, RequestStatus } from "@/types";
import "./DetailsCard.css";
import SectionTitle from "@/components/SectionTitle/SectionTitle";
import StatusBadge from "@/components/Request/StatusBadge/StatusBadge";

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
        <span>
          {new Date(data.created_at).toLocaleDateString("us-US", {
            day: "2-digit",
            month: "long",
            year: "numeric",
          })}
        </span>
      </div>
      <div className="info-row">
        <strong>Status:</strong>
        <span>
          <StatusBadge status={data.status as RequestStatus} />
        </span>
      </div>
      {data.description && (
        <div className="info-description">
          <strong>Description:</strong>
          <p>{data.description}</p>
        </div>
      )}
    </section>
  );
}
