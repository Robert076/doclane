import { notFound } from "next/navigation";

import FileSection from "@/components/FileSectionComponents/FileSection/FileSection";
import DetailsCard from "@/components/Pages/RequestsComponents/DetailsCard";
import DetailCardsActionSidebar from "@/components/Pages/RequestsComponents/RequestDetailsActions";
import RequestDetailsHeader from "@/components/Pages/RequestsComponents/RequestDetailsHeader";
import { getDocumentRequestById, getFilesByRequestId } from "@/lib/api/requests";

interface PageProps {
        params: Promise<{ id: string }>;
}

export default async function RequestDetailsPage({ params }: PageProps) {
        const { id } = await params;

        const [request, filesResponse] = await Promise.all([
                getDocumentRequestById(id),
                getFilesByRequestId(id),
        ]);

        if (!request || !request.data) {
                notFound();
        }

        const data = request.data;
        const files = filesResponse?.data || [];

        return (
                <div className="details-container">
                        <RequestDetailsHeader data={data} />

                        <div className="details-grid">
                                <div className="main-content">
                                        <DetailsCard data={data} />

                                        <FileSection
                                                files={files}
                                                expectedDocuments={
                                                        data.expected_documents || []
                                                }
                                                requestId={id}
                                        />
                                </div>
                                <DetailCardsActionSidebar id={id} />
                        </div>
                </div>
        );
}
