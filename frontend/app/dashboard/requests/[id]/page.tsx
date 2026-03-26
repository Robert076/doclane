import { notFound } from "next/navigation";
import DetailCardsActionSidebar from "@/components/Pages/RequestsComponents/RequestDetailsActions";
import RequestDetailsHeader from "@/components/Pages/RequestsComponents/RequestDetailsHeader";
import RequestTabs from "@/components/Pages/RequestsComponents/RequestTabs";
import { getRequestById, getFilesByRequestId, getCommentsByRequest } from "@/lib/api/requests";

interface PageProps {
        params: Promise<{ id: string }>;
}

export default async function RequestDetailsPage({ params }: PageProps) {
        const { id } = await params;
        const [request, filesResponse, commentsResponse] = await Promise.all([
                getRequestById(id),
                getFilesByRequestId(id),
                getCommentsByRequest(+id),
        ]);

        if (!request || !request.data) {
                notFound();
        }

        const data = request.data;
        const files = filesResponse?.data || [];
        const comments = commentsResponse?.data || [];

        return (
                <div className="details-container">
                        <RequestDetailsHeader data={data} />
                        <div className="details-grid">
                                <div className="main-content">
                                        <RequestTabs
                                                data={data}
                                                files={files}
                                                comments={comments}
                                                requestId={id}
                                        />
                                </div>
                                <DetailCardsActionSidebar id={id} />
                        </div>
                </div>
        );
}
