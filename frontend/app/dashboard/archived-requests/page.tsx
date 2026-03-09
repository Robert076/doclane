import { notFound } from "next/navigation";
import { DocumentRequest } from "@/types";
import ArchivedRequestsSection from "@/components/Pages/ArchivedRequestsComponents/ArchivedRequestsSection";
import PageHeader from "@/components/PageHeader/PageHeader";
import { getDocumentRequests } from "@/lib/api/requests";

const ArchivedRequests = async () => {
        const requestsResponse = await getDocumentRequests("PROFESSIONAL");

        if (!requestsResponse?.data) {
                notFound();
        }

        const requests = requestsResponse.data as DocumentRequest[];

        return (
                <div>
                        <PageHeader
                                title="Dosare arhivate"
                                subtitle="Restaurează şi gestionează dosarele arhivate."
                        />
                        <ArchivedRequestsSection requests={requests} />
                </div>
        );
};

export default ArchivedRequests;
