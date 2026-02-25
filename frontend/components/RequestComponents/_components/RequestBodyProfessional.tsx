import { DocumentRequest } from "@/types";
import React from "react";
import "./RequestBody.css";
import RequestInfoItem from "./RequestInfoItem";
import { UI_TEXT } from "@/locales/ro";

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
                                        UI_TEXT.request.card.clientEmail,
                                        request.client_email,
                                )}
                                {RequestInfoItem(
                                        searchTerm,
                                        UI_TEXT.request.card.clientName,
                                        `${request.client_first_name} ${request.client_last_name}`,
                                )}
                        </div>
                </div>
        );
};

export default RequestBody;
