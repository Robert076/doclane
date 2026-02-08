import { User } from "@/types";
import React from "react";
import ClientHeader from "./_components/ClientHeader";
import getDocumentRequests from "@/lib/api/getDocumentRequests";
import RequestsSection from "@/components/RequestsSection/RequestsSection";

interface ClientViewProps {
  user: User;
}

export default async function ClientView({ user }: ClientViewProps) {
  const requests = await getDocumentRequests(user.role);

  return (
    <div className="client-view">
      <ClientHeader user={user} length={requests.length} />
      <RequestsSection user={user} requests={requests} />
    </div>
  );
}
