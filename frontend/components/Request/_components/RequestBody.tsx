import { DocumentRequest } from "@/types";
import React from "react";
import "./RequestBody.css";

interface RequestBodyProps {
  request: DocumentRequest;
}

const RequestBody: React.FC<RequestBodyProps> = ({ request }) => {
  return (
    <div className="request-body">
      <p className="request-info">
        <strong>Client:</strong> {request.client_email}
      </p>
      {request.description && <p className="request-desc">{request.description}</p>}
    </div>
  );
};

export default RequestBody;
