import React from "react";
import { Request } from "@/types";
import RequestInfoItem from "./RequestInfoItem";
import "./RequestBody.css";

interface RequestBodyProps {
        request: Request;
        searchTerm?: string;
}

export default function RequestBodyProfessional({ request, searchTerm }: RequestBodyProps) {
        return (
                <div className="request-body">
                        <div className="request-info">
                                <RequestInfoItem
                                        label="Email client"
                                        value={request.client_email}
                                        searchTerm={searchTerm}
                                />
                                <RequestInfoItem
                                        label="Nume client"
                                        value={`${request.client_first_name} ${request.client_last_name}`}
                                        searchTerm={searchTerm}
                                />
                        </div>
                </div>
        );
}
