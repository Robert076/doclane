import { DocumentRequest } from "@/types";
import React from "react";
import HighlightText from "@/components/HighlightText/HighlightText";
import "./RequestBody.css";

interface RequestBodyProps {
  request: DocumentRequest;
  searchTerm?: string;
}

const formatDate = (dateString: string) => {
  const date = new Date(dateString);
  return date.toLocaleDateString("en-US", {
    year: "numeric",
    month: "long",
    day: "numeric",
  });
};

const RequestBody: React.FC<RequestBodyProps> = ({ request, searchTerm }) => {
  return (
    <div className="request-body">
      <div className="request-info">
        {request.created_at && (
          <p className="request-info-item">
            <span className="request-label">Created at:</span>
            <span className="request-value">
              <HighlightText text={`${formatDate(request.created_at)}`} search={searchTerm} />
            </span>
          </p>
        )}
      </div>
    </div>
  );
};

export default RequestBody;
