import { DocumentRequest } from "@/types";
import React from "react";
import HighlightText from "@/components/HighlightText/HighlightText";
import "./RequestBody.css";
import { formatDate } from "@/lib/formatDate";

interface RequestBodyProps {
  request: DocumentRequest;
  searchTerm?: string;
}

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
