import { notFound } from "next/navigation";

import getFilesByRequestId from "@/lib/getFilesByRequestId";
import FileSection from "@/components/FileSectionComponents/FileSection/FileSection";
import DetailsHeader from "./_components/DetailsHeader/DetailsHeader";
import DetailsCard from "./_components/DetailsCard/DetailsCard";
import DetailCardsActionSidebar from "./_components/DetailCardsActionSidebar/DetailCardsActionSidebar";
import "./style.css";
import { getDocumentRequestById } from "@/lib/api/api";

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
                        <DetailsHeader data={data} />
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
