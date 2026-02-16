import { User } from "@/types";
import React from "react";
import ClientHeader from "./_components/ClientHeader";
import { getDocumentRequests } from "@/lib/api/api";
import RequestsSection from "@/components/RequestComponents/RequestsSection/RequestsSection";

interface ClientViewProps {
        user: User;
}

export default async function ClientView({ user }: ClientViewProps) {
        const response = await getDocumentRequests(user.role);

        return (
                <div className="client-view">
                        <ClientHeader user={user} length={response.data.length} />
                        <RequestsSection user={user} requests={response.data} />
                </div>
        );
}
