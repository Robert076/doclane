import { DocumentRequest } from "@/types";
import React from "react";
import HighlightText from "@/components/HighlightText/HighlightText";
import "./RequestBody.css";

interface RequestBodyProps {
  request: DocumentRequest;
  searchTerm?: string;
}

const RequestBody: React.FC<RequestBodyProps> = ({ request, searchTerm }) => {
  return (
    <div className="request-body">
      <div className="request-info">
        <p className="request-info-item">
          <span className="request-label">Client email:</span>
          <span className="request-value">
            <HighlightText text={request.client_email} search={searchTerm} />
          </span>
        </p>
        <p className="request-info-item">
          <span className="request-label">Client name:</span>
          <span className="request-value">
            <HighlightText
              text={`${request.client_first_name} ${request.client_last_name}`}
              search={searchTerm}
            />
          </span>
        </p>
      </div>
    </div>
  );
};

export default RequestBody;
