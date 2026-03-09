import React from "react";
import { DocumentRequest } from "@/types";
import RequestInfoItem from "./RequestInfoItem";
import { formatDate } from "@/lib/client/formatDate";
import "./RequestBody.css";

interface RequestBodyProps {
        request: DocumentRequest;
        searchTerm?: string;
}

export default function RequestBodyClient({ request, searchTerm }: RequestBodyProps) {
        return (
                <div className="request-body">
                        <div className="request-info">
                                <RequestInfoItem
                                        label="Creat la"
                                        value={formatDate(request.created_at)}
                                        searchTerm={searchTerm}
                                />
                        </div>
                </div>
        );
}
