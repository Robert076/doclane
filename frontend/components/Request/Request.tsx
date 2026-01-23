"use client";

import { useRouter } from "next/navigation";
import { DocumentRequest, RequestStatus } from "@/types";
import StatusBadge from "./StatusBadge/StatusBadge";
import ButtonPrimary from "../Buttons/ButtonPrimary/ButtonPrimary";
import "./Request.css";
import RequestBody from "./_components/RequestBody";

interface RequestProps {
  request: DocumentRequest;
}

const Request: React.FC<RequestProps> = ({ request }) => {
  const router = useRouter();
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

      <RequestBody request={request} />

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
