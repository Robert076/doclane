import PageHeader from "@/components/PageHeader/PageHeader";
import RequestsSection from "@/components/Pages/RequestsComponents/RequestsSection";
import { getCurrentUser, getDocumentRequests } from "@/lib/api/api";
import { notFound, redirect } from "next/navigation";

export default async function RequestsPage() {
        const userResponse = await getCurrentUser();
        if (!userResponse.success || !userResponse.data) {
                redirect("/login");
        }

        const user = userResponse.data;

        const requestsResponse = await getDocumentRequests(userResponse.data.role);
        if (!requestsResponse.success || !requestsResponse.data) {
                notFound();
        }

        const requests = requestsResponse.data;

        return (
                <div>
                        <PageHeader
                                title="Dosarele tale"
                                subtitle="Administrează şi gestionează dosarele tale deschise."
                        />
                        <RequestsSection user={user} requests={requests} />
                </div>
        );
}
