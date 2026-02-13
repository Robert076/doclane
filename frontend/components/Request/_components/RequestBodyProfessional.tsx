import { DocumentRequest } from "@/types";
import React from "react";
import "./RequestBody.css";
import RequestInfoItem from "./RequestInfoItem";

interface RequestBodyProps {
        request: DocumentRequest;
        searchTerm?: string;
}

const RequestBody: React.FC<RequestBodyProps> = ({ request, searchTerm }) => {
        return (
                <div className="request-body">
                        <div className="request-info">
                                {RequestInfoItem(
                                        searchTerm,
                                        "Client email:",
                                        request.client_email,
                                )}
                                {RequestInfoItem(
                                        searchTerm,
                                        "Client name:",
                                        `${request.client_first_name} ${request.client_last_name}`,
                                )}
                        </div>
                </div>
        );
};

export default RequestBody;
