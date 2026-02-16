import { DocumentRequest } from "@/types";
import React from "react";
import "./RequestBody.css";
import RequestInfoItem from "./RequestInfoItem";
import { formatDate } from "@/lib/client/formatDate";

interface RequestBodyProps {
        request: DocumentRequest;
        searchTerm?: string;
}

const RequestBody: React.FC<RequestBodyProps> = ({ request, searchTerm }) => {
        return (
                <div className="request-body">
                        <div className="request-info">
                                {request.created_at &&
                                        RequestInfoItem(
                                                searchTerm,
                                                "Created at:",
                                                formatDate(request.created_at),
                                        )}
                        </div>
                </div>
        );
};

export default RequestBody;
