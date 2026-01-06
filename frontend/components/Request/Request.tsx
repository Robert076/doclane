"use client";

import { useRouter } from "next/navigation"; // Importăm hook-ul de navigare
import { DocumentRequest, RequestStatus } from "@/types";
import StatusBadge from "./StatusBadge/StatusBadge";
import ButtonPrimary from "../Buttons/ButtonPrimary/ButtonPrimary";
import "./Request.css";

interface RequestProps {
  request: DocumentRequest;
}

const Request: React.FC<RequestProps> = ({ request }) => {
  const router = useRouter(); // Inițializăm router-ul
  const isOverdue = request.status === "overdue";

  const handleViewDetails = () => {
    router.push(`/dashboard/requests/${request.id}`);
  };

  return (
    <div className={`document-request-card ${isOverdue ? "is-overdue" : ""}`}>
      <div className="request-header">
        <StatusBadge status={request.status as RequestStatus} />
      </div>

      <h3 className="request-title">{request.title}</h3>

      <div className="request-body">
        <p className="request-info">
          <strong>Client:</strong> {request.client_email}
        </p>
        {request.description && <p className="request-desc">{request.description}</p>}
      </div>

      <div className="request-footer">
        <ButtonPrimary
          text="View Details"
          variant="ghost"
          fullWidth={true}
          onClick={handleViewDetails}
        />
      </div>
    </div>
  );
};

export default Request;
