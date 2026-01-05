"use client";
import { DocumentRequest, RequestStatus } from "@/types";
import StatusBadge from "./StatusBadge/StatusBadge";

import "./Request.css";
import ButtonPrimary from "../Buttons/ButtonPrimary/ButtonPrimary";

interface RequestProps {
  request: DocumentRequest;
}

const Request: React.FC<RequestProps> = ({ request }) => {
  const isOverdue = request.status === "overdue";

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
          onClick={() => console.log("Details for:", request.id)}
        />
      </div>
    </div>
  );
};

export default Request;
