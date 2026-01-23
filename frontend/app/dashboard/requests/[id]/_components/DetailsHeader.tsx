import { DocumentRequest } from "@/types";
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
      </div>
    </header>
  );
}
