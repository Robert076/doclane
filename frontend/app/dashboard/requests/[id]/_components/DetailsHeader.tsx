import StatusBadge from "@/components/Request/StatusBadge/StatusBadge";
import { DocumentRequest, RequestStatus } from "@/types";
import Link from "next/link";
import "./DetailsHeader.css";

export default function DetailsHeader({ data }: { data: DocumentRequest }) {
  return (
    <header className="details-header">
      <Link href="/dashboard" className="back-link">
        ← Înapoi la Cereri
      </Link>
      <div className="header-main">
        <h1>{data.title}</h1>
        <StatusBadge status={data.status as RequestStatus} />
      </div>
    </header>
  );
}
