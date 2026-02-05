"use client";
import { useRouter } from "next/navigation";
import { DocumentRequest, RequestStatus } from "@/types";
import StatusBadge from "./StatusBadge/StatusBadge";
import ButtonPrimary from "../Buttons/ButtonPrimary/ButtonPrimary";
import HighlightText from "../HighlightText/HighlightText";
import "./Request.css";
import RequestBody from "./_components/RequestBody";

interface RequestProps {
  request: DocumentRequest;
  searchTerm?: string;
}

const Request: React.FC<RequestProps> = ({ request, searchTerm }) => {
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
      <h3 className="request-title">
        <HighlightText text={request.title} search={searchTerm} />
      </h3>
      <RequestBody request={request} searchTerm={searchTerm} />
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
